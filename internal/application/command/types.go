package command

type File struct {
	Files []string
}
type StartSendFileInput struct {
	NodeId   string
	BaseName string
	Files    []string // 如果是目录，那么是该目录下的所有文件列表
}
type StartSendFileOutput struct{}

type ClientIdsOutput struct {
	Ids []string
}

type PauseSendFileInput struct {
	TaskId string
}
type PauseSendFileOutput struct{}

type CancelSendFileInput struct {
	TaskId string
}
type CancelSendFileOutput struct{}

type ResumeSendFileInput struct {
	TaskId string
}
type ResumeSendFileOutput struct{}

type RemoveTaskInput struct {
	TaskIds []string
}
type RemoveTaskOutput struct{}
