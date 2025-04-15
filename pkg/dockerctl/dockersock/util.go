package dockersock

import (
	"context"
	"fmt"
	"math"
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/gogf/gf/v2/errors/gerror"
)

type PortConfig struct {
	Port         string
	OriginalPort string
	Protocol     string
	Services     []string
}

type PortConfigs []PortConfig

func (c PortConfigs) String() string {
	s := ""
	for _, cfg := range c {
		s += fmt.Sprintf(
			" %s/%s\t%s\n",
			cfg.Port,
			cfg.Protocol,
			strings.Join(cfg.Services, ", "),
		)
	}
	return s
}
func protoPort(port string, proto string) string {
	return fmt.Sprintf("%s/%s", port, proto)
}

// resolvePortConflicts 将冲突的端口映射到可用的端口上
func resolvePortConflicts(conflicts PortConfigs) (PortConfigs, error) {
	// indexes reserved ports for easier lookup
	protoPortCfgs := map[string]bool{}
	for _, cfg := range conflicts {
		protoPortCfgs[protoPort(cfg.Port, cfg.Protocol)] = true
	}

	resolved := []PortConfig{}

	for _, cfg := range conflicts {
		for i, service := range cfg.Services {
			// skip the first service, as it's either the host or we want to keep this service's ports and change
			// the ports of the other services
			// 第一个是host， 或者保持端口号不变， 接下来的service改变端口号
			if i == 0 {
				resolved = append(resolved, PortConfig{
					Port:         cfg.Port,
					OriginalPort: cfg.Port,
					Protocol:     cfg.Protocol,
					Services:     []string{service},
				})
			}

			// check the next port
			intPort, _ := strconv.Atoi(cfg.Port)

			port := intPort + 1
			for {
				pp := protoPort(fmt.Sprint(port), cfg.Protocol)

				if _, ok := protoPortCfgs[pp]; !ok && portAvailable(fmt.Sprint(port), cfg.Protocol) {
					// we found an available port that is not already reserved
					protoPortCfgs[pp] = true
					break
				}

				port += 1
				if port == math.MaxUint16 {
					return nil, fmt.Errorf("no ports available")
				}
			}

			resolved = append(resolved, PortConfig{
				Port:         fmt.Sprint(port),
				OriginalPort: cfg.Port,
				Protocol:     cfg.Protocol,
				Services:     []string{service},
			})
		}
	}
	return resolved, nil
}

// 更新project中的端口信息
func applyPortMapping(p *types.Project, mapping PortConfigs) {
	// index that associates service with its port configs
	servicePorts := map[string]PortConfigs{}
	for _, cfg := range mapping {
		service := cfg.Services[0]
		servicePorts[service] = append(servicePorts[service], cfg)
	}

	for i, service := range p.Services {
		servicePortCfgs := servicePorts[service.Name]

		for j, port := range service.Ports {
			for _, cfg := range servicePortCfgs {
				if port.Published == cfg.OriginalPort && port.Protocol == cfg.Protocol {
					port.Published = cfg.Port
					service.Ports[j] = port
					p.Services[i] = service
				}
			}
		}
	}
}

// portConflicts 返回有冲突的端口列表
func portConflicts(cfgs PortConfigs) PortConfigs {
	conflicts := []PortConfig{}

	for _, cfg := range cfgs {
		if len(cfg.Services) > 1 {
			// port conflict detected
			conflicts = append(conflicts, cfg)
		}
	}

	return conflicts
}

// hrojectPortConfigs 检查端口是否被多个服务使用
func hrojectPortConfigs(cfgs PortConfigs) bool {
	for _, cfg := range cfgs {
		if len(cfg.Services) > 1 {
			// port conflict detected
			return true
		}
	}

	return false
}

// 根据type获取端口配置
func portConfigs(proj *types.Project, typ string) PortConfigs {
	portServices := map[string][]string{}

	for _, s := range proj.Services {
		for _, spCfg := range s.Ports {
			port := spCfg.Published // 对外开放的主机端口
			if spCfg.Protocol == typ {
				if !portAvailable(port, typ) && len(portServices[port]) == 0 { // 端口不可用。记录为host
					portServices[port] = append(portServices[port], "host")
				}

				portServices[port] = append(portServices[port], s.Name) // 端口可用。记录为服务名称,可能会记录多个sevice(端口冲突)
			}
		}
	}

	portCfgs := []PortConfig{}
	for port, services := range portServices {
		portCfgs = append(portCfgs, PortConfig{
			Port:         port,
			OriginalPort: port,
			Protocol:     typ,
			Services:     services,
		})
	}

	return portCfgs
}

// projectPortConfigs 获取需要映射的端口列表
func projectPortConfigs(p *types.Project) PortConfigs {
	tcpCfg := portConfigs(p, "tcp")
	udpCfg := portConfigs(p, "udp")
	cfgs := append(tcpCfg, udpCfg...)

	// sort ports for consistent ordering
	sort.Slice(cfgs, func(i, j int) bool {
		return protoPort(cfgs[i].Port, cfgs[i].Protocol) > protoPort(cfgs[j].Port, cfgs[j].Protocol)
	})

	return cfgs
}

var defaultDockerComposeYmlPath = "docker-compose.yml"

// ProjectFromConfig 加载 docker-compose 配置文件，并实例化一个docker compose的project
func ProjectFromConfig(composePath string) (p *types.Project, err error) {
	if composePath == "" {
		composePath = defaultDockerComposeYmlPath
	}
	opts, err := cli.NewProjectOptions(
		[]string{composePath},
		cli.WithDotEnv,
	)
	if err != nil {
		return nil, gerror.Wrap(err, "new project fail")
	}
	// 创建一个project实例
	p, err = cli.ProjectFromOptions(context.Background(), opts)
	if err != nil {
		return nil, gerror.Wrap(err, "error loading docker-compose file")
	}
	// 设置service的标签
	for i, s := range p.Services {
		s.CustomLabels = map[string]string{
			api.ProjectLabel:     p.Name,
			api.ServiceLabel:     s.Name,
			api.VersionLabel:     api.ComposeVersion,
			api.WorkingDirLabel:  p.WorkingDir,
			api.ConfigFilesLabel: strings.Join(p.ComposeFiles, ","),
			api.OneoffLabel:      "False",
		}
		p.Services[i] = s
	}

	return p, nil
}

// portAvailable 判断端口是否可用
func portAvailable(port string, proto string) bool {
	switch proto {
	case "tcp":
		ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
		if err != nil {
			return false
		}
		_ = ln.Close()
		return true
	case "udp":
		ln, err := net.ListenPacket("udp", fmt.Sprintf(":%s", port))
		if err != nil {
			return false
		}
		_ = ln.Close()
		return true
	default:
		return false
	}
}
