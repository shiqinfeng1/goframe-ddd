package dto

type UpgradeAppInput struct {
	AppName string
}
type UpgradeImageInput struct {
	Version string
}

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
