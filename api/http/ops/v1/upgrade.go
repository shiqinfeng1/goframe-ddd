package v1

import "github.com/gogf/gf/v2/frame/g"

type UpgradeAppReq struct {
	g.Meta  `path:"/ops/appUpgrade" tags:"运维" method:"post" summary:"应用版本升级"`
	AppName string `p:"app_name" v:"required|in:mgrid#请输入应用名称|目前只支持平滑重启的应用:mgrid" dc:"应用名称"`
}
type UpgradeAppRes struct {
	g.Meta `status:"200"`
}

type UpgradeImageReq struct {
	g.Meta  `path:"/ops/imageUpgrade" tags:"运维" method:"post" summary:"镜像版本升级"`
	Version string `p:"version" v:"required#请输入镜像版本" dc:"镜像版本号"`
}
type UpgradeImageRes struct {
	g.Meta `status:"200"`
}
