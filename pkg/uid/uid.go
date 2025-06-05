package uid

import (
	"fmt"
	"net"
	"os"

	"github.com/gogf/gf/v2/text/gstr"
)

// HostnameProvider 定义获取主机名的接口
type HostnameProvider interface {
	Hostname() (string, error)
}

// NetworkProvider 定义获取网络接口相关信息的接口
type NetworkProvider interface {
	Interfaces() ([]net.Interface, error)
	InterfaceAddrs(i net.Interface) ([]net.Addr, error)
}

type backend struct{}

var DefaultHostnameProvider HostnameProvider
var DefaultNetworkProvider NetworkProvider

func init() {
	DefaultHostnameProvider = &backend{}
	DefaultNetworkProvider = &backend{}
}

func (p *backend) Hostname() (string, error) {
	return os.Hostname()
}
func (p *backend) Interfaces() ([]net.Interface, error) {
	return net.Interfaces()
}
func (p *backend) InterfaceAddrs(i net.Interface) ([]net.Addr, error) {
	return i.Addrs()
}

func ClientIDIsValid(uid string) bool {
	suid := gstr.Split(uid, "-")
	// 检查uid格式
	if len(suid) < 2 {
		return false
	}
	// 检查mac
	if _, err := net.ParseMAC(suid[len(suid)-1]); err != nil {
		return false
	}
	return true
}

// GetFirstValidMACAddress 获取第一个有效的 MAC 地址
func GetFirstValidMACAddress(np NetworkProvider) (string, error) {
	interfaces, err := np.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		// 检查是否有IPv4或IPv6地址
		addrs, err := np.InterfaceAddrs(iface)
		if err != nil || len(addrs) == 0 {
			continue
		}

		// 跳过无效的MAC地址 跳过环回接口 跳过未启用的接口
		if len(iface.HardwareAddr) == 0 || iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		hasNonLocalIP := false
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipNet.IP
			// 跳过本地地址和链路本地地址
			if !ip.IsLoopback() && !ip.IsLinkLocalUnicast() {
				hasNonLocalIP = true
				break
			}
		}

		if !hasNonLocalIP {
			continue
		}

		hwAddr := iface.HardwareAddr.String()

		if iface.HardwareAddr.String() == "00:00:00:00:00:00" {
			continue
		}

		return hwAddr, nil
	}

	return "", fmt.Errorf("no valid mac address found")
}

// NewClientID 为当前主机生成一个全局唯一uid，接收 HostnameProvider 和 NetworkProvider 接口实例作为参数
func NewClientID(hp HostnameProvider, np NetworkProvider) (string, error) {
	// 获取主机名
	hostname, err := hp.Hostname()
	if err != nil {
		return "nohostname", err
	}

	// 获取第一个非空的 MAC 地址
	macAddr, err := GetFirstValidMACAddress(np)
	if err != nil {
		return hostname + "-nomac", err
	}

	// 组合信息
	return fmt.Sprintf("%s-%s", hostname, macAddr), nil
}

// NewClientIDWithDefault 使用默认 Provider 生成客户端 ID
func NewClientIDWithDefault() (string, error) {
	return NewClientID(DefaultHostnameProvider, DefaultNetworkProvider)
}
