package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type ImageSummary struct {
	Name string `json:"name" dc:"镜像名称"`
	Tag  string `json:"tag" dc:"镜像标签"`
}
type ImagesOutput struct {
	Images []ImageSummary `json:"images" dc:"当前运行容器对应的镜像"`
}
type ComposeImagesOutput struct {
	Images []ImageSummary `json:"images" dc:"当前运行容器对应的镜像"`
}

type ImageListReq struct {
	g.Meta `path:"/ops/imageList" tags:"运维" method:"post" summary:"镜像列表"`
}
type ImageListRes struct {
	g.Meta `status:"200"`
	Images []ImageSummary `json:"images" dc:"所有镜像列表"`
}

type ContainerImageReq struct {
	g.Meta `path:"/ops/containerImage" tags:"运维" method:"post" summary:"当前容器的镜像"`
}
type ContainerImageRes struct {
	g.Meta `status:"200"`
	Images []ImageSummary `json:"images" dc:"当前运行容器对应的镜像"`
}
