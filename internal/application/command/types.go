package command

type File struct {
	Files []string
}
type StartSendFileInput struct {
	BaseName string
	IsDir    bool
	Files    []string // 如果是目录，那么是该目录下的所有文件列表
}

type StartSendFileOutput struct{}
