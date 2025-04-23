package v1

import "github.com/gogf/gf/v2/frame/g"

type DeviceAuthReq struct {
	g.Meta   `path:"/auth/device" tags:"认证" method:"post" summary:"访问权限认证"`
	APIKey   string `header:"X-API-Key" dc:"api-key"`
	TenantID string `header:"X-Tenant-ID" dc:"租户id"`
	Role     string `header:"X-Role" dc:"角色"`
}

type DeviceAuthRes struct {
	g.Meta  `status:"200"`
	Running int `json:"runnings" dc:"正在运行的发送任务数量"`
}
