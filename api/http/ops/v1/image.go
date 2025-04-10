package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/internal/application"
)

type ImageListReq struct {
	g.Meta `path:"/ops/imageList" tags:"运维" method:"post" summary:"镜像列表"`
}
type ImageListRes struct {
	g.Meta `status:"200"`
	Images []application.ImageSummary `json:"images" dc:"所有镜像列表"`
}

type ContainerImageReq struct {
	g.Meta `path:"/ops/containerImage" tags:"运维" method:"post" summary:"当前容器的镜像"`
}
type ContainerImageRes struct {
	g.Meta `status:"200"`
	Images []application.ImageSummary `json:"images" dc:"当前运行容器对应的镜像"`
}
