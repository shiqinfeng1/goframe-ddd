package utils

import (
	"fmt"
	"net"
	"os"

	"github.com/gogf/gf/v2/text/gstr"
)

func UidIsValid(uid string) bool {
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

// 为当前主机生成一个全局唯一uid
func GenUIDForHost() (string, error) {
	// 获取主机名
	hostname, err := os.Hostname()
	if err != nil {
		return "nohostname", err
	}

	// 获取第一个非空的 MAC 地址
	var macAddr string
	ifaces, err := net.Interfaces()
	if err != nil {
		return hostname + "-nointf", err
	}
	for _, i := range ifaces {
		if i.Flags&net.FlagUp != 0 && i.Flags&net.FlagLoopback == 0 {
			addrs, err := i.Addrs()
			if err == nil {
				for _, addr := range addrs {
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					}
					if ip != nil && !ip.IsLoopback() {
						macAddr = i.HardwareAddr.String()
						break
					}
				}
			}
		}
		if macAddr != "" {
			break
		}
	}
	if macAddr == "" {
		return hostname + "-nomac", fmt.Errorf("未找到有效的 MAC 地址")
	}

	// 组合信息
	return fmt.Sprintf("%s-%s", hostname, macAddr), nil
}
