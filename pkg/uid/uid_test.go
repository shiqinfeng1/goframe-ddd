package uid_test

import (
	"errors"
	"net"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/shiqinfeng1/goframe-ddd/pkg/uid"
	"github.com/shiqinfeng1/goframe-ddd/pkg/uid/mocks"

	"github.com/stretchr/testify/assert"
)

// TestGetFirstValidMACAddress 测试 GetFirstValidMACAddress 函数
func TestGetFirstValidMACAddress(t *testing.T) {
	tests := []struct {
		name string
		// 模拟网络接口列表
		interfaces []net.Interface
		// 获取网络接口列表时的错误
		interfacesErr error
		// 每个接口对应的地址列表
		addrs []net.Addr
		// 获取接口地址时的错误
		addrsErr error
		// 期望的MAC地址
		expectedMAC string
		// 期望的错误
		expectedErr error
	}{
		{
			name: "成功获取有效的MAC地址",
			interfaces: []net.Interface{
				{
					Name:         "eth0",
					HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
					Flags:        net.FlagUp,
				},
			},
			interfacesErr: nil,
			addrs:         []net.Addr{&net.IPNet{IP: net.ParseIP("192.168.1.100")}},
			addrsErr:      nil,
			expectedMAC:   "00:11:22:33:44:55",
			expectedErr:   nil,
		},
		{
			name:          "获取网络接口列表失败",
			interfaces:    nil,
			interfacesErr: errors.New("获取网络接口列表失败"),
			addrs:         nil,
			addrsErr:      nil,
			expectedMAC:   "",
			expectedErr:   errors.New("获取网络接口列表失败"),
		},
		{
			name: "未找到有效的MAC地址（环回接口）",
			interfaces: []net.Interface{
				{
					Name:         "lo",
					HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
					Flags:        net.FlagLoopback,
				},
			},
			interfacesErr: nil,
			addrs:         []net.Addr{&net.IPNet{IP: net.ParseIP("127.0.0.1")}},
			addrsErr:      nil,
			expectedMAC:   "",
			expectedErr:   errors.New("no valid mac address found"),
		},
		{
			name: "未找到有效的MAC地址（接口未启用）",
			interfaces: []net.Interface{
				{
					Name:         "eth0",
					HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
					Flags:        0,
				},
			},
			interfacesErr: nil,
			addrs:         []net.Addr{&net.IPNet{IP: net.ParseIP("192.168.1.100")}},
			addrsErr:      nil,
			expectedMAC:   "",
			expectedErr:   errors.New("no valid mac address found"),
		},
		{
			name: "未找到有效的MAC地址（无硬件地址）",
			interfaces: []net.Interface{
				{
					Name:  "eth0",
					Flags: net.FlagUp,
				},
			},
			interfacesErr: nil,
			addrs:         []net.Addr{&net.IPNet{IP: net.ParseIP("192.168.1.100")}},
			addrsErr:      nil,
			expectedMAC:   "",
			expectedErr:   errors.New("no valid mac address found"),
		},
		{
			name: "未找到有效的MAC地址（仅本地IP）",
			interfaces: []net.Interface{
				{
					Name:         "eth0",
					HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
					Flags:        net.FlagUp,
				},
			},
			interfacesErr: nil,
			addrs:         []net.Addr{&net.IPNet{IP: net.ParseIP("127.0.0.1")}},
			addrsErr:      nil,
			expectedMAC:   "",
			expectedErr:   errors.New("no valid mac address found"),
		},
		{
			name: "未找到有效的MAC地址（全零MAC地址）",
			interfaces: []net.Interface{
				{
					Name:         "eth0",
					HardwareAddr: net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
					Flags:        net.FlagUp,
				},
			},
			interfacesErr: nil,
			addrs:         []net.Addr{&net.IPNet{IP: net.ParseIP("192.168.1.100")}},
			addrsErr:      nil,
			expectedMAC:   "",
			expectedErr:   errors.New("no valid mac address found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNP := mocks.NewNetworkProvider(t)
			mockNP.On("Interfaces").Return(tt.interfaces, tt.interfacesErr)

			if tt.interfacesErr == nil {
				for _, iface := range tt.interfaces {
					mockNP.On("InterfaceAddrs", iface).Return(tt.addrs, tt.addrsErr)
				}
			}

			mac, err := uid.GetFirstValidMACAddress(mockNP)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedMAC, mac)

			mockNP.AssertExpectations(t)
		})
	}
}

// TestUidIsValid 测试 ClientIDIsValid 函数
func TestUidIsValid(t *testing.T) {
	tests := []struct {
		name     string
		uid      string
		expected bool
	}{
		{
			name:     "valid uid with MAC",
			uid:      "prefix-00:1A:2B:3C:4D:5E",
			expected: true,
		},
		{
			name:     "invalid uid - no delimiter",
			uid:      "noDelimiter",
			expected: false,
		},
		{
			name:     "invalid uid - empty string",
			uid:      "",
			expected: false,
		},
		{
			name:     "invalid uid - only prefix",
			uid:      "prefix-",
			expected: false,
		},
		{
			name:     "invalid uid - invalid MAC format",
			uid:      "prefix-invalidMAC",
			expected: false,
		},
		{
			name:     "invalid uid - incomplete MAC",
			uid:      "prefix-00:1A:2B",
			expected: false,
		},
		{
			name:     "valid uid with lowercase MAC",
			uid:      "prefix-00:1a:2b:3c:4d:5e",
			expected: true,
		},
		{
			name:     "invalid uid with hyphenated MAC",
			uid:      "prefix-00-1A-2B-3C-4D-5E",
			expected: false,
		},
		{
			name:     "invalid uid - multiple delimiters but invalid MAC",
			uid:      "prefix1-prefix2-invalidMAC",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uid.ClientIDIsValid(tt.uid)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestNewClientID 测试 NewClientID 函数
func TestNewClientID(t *testing.T) {
	tests := []struct {
		name        string
		hostname    string
		hostnameErr error
		macAddr     string
		macAddrErr  error
		expectedID  string
		expectedErr error
	}{
		{
			name:        "正常生成客户端ID",
			hostname:    "testhost",
			hostnameErr: nil,
			macAddr:     "00:11:22:33:44:55",
			macAddrErr:  nil,
			expectedID:  "testhost-00:11:22:33:44:55",
			expectedErr: nil,
		},
		{
			name:        "获取主机名失败",
			hostname:    "",
			hostnameErr: errors.New("获取主机名失败"),
			macAddr:     "",
			macAddrErr:  nil,
			expectedID:  "nohostname",
			expectedErr: errors.New("获取主机名失败"),
		},
		{
			name:        "获取MAC地址失败",
			hostname:    "testhost",
			hostnameErr: nil,
			macAddr:     "",
			macAddrErr:  errors.New("获取MAC地址失败"),
			expectedID:  "testhost-nomac",
			expectedErr: errors.New("获取MAC地址失败"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建 Mock 对象
			mockHP := mocks.NewHostnameProvider(t)
			mockNP := mocks.NewNetworkProvider(t)

			// 设置 Mock 对象的行为
			mockHP.On("Hostname").Return(tt.hostname, tt.hostnameErr)

			// 使用 mockey 替换 GetFirstValidMACAddress 函数
			patch := gomonkey.ApplyFunc(uid.GetFirstValidMACAddress, func(np uid.NetworkProvider) (string, error) {
				return tt.macAddr, tt.macAddrErr
			})
			defer patch.Reset()

			// 调用被测试函数
			clientID, err := uid.NewClientID(mockHP, mockNP)

			// 断言结果
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedID, clientID)

			// 验证 Mock 对象的调用情况
			mockHP.AssertExpectations(t)
			mockNP.AssertExpectations(t)
		})
	}
}
