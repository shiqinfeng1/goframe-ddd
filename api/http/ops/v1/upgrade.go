package v1

import "github.com/gogf/gf/v2/frame/g"

type UpgradeReq struct {
	g.Meta  `path:"/ops/appUpgrade" tags:"运维" method:"post" summary:"应用版本升级"`
	AppName string `p:"app_name" v:"required|in:mgrid" dc:"应用名称"`
}
type UpgradeRes struct {
	g.Meta `status:"200"`
}
