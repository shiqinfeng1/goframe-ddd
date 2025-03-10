package command

import (
	"context"

	"github.com/rs/xid"
)

func (h *Handler) StartSendFile(ctx context.Context, in *StartSendFileInput) (*StartSendFileOutput, error) {
	taskId := xid.New().String()
	h.fileTransfer.AddTask(ctx, taskId, in.BaseName, in.NodeId, in.Files)
	return nil, nil
}

func (h *Handler) PauseSendFile(ctx context.Context, in *PauseSendFileInput) (*PauseSendFileOutput, error) {
	h.fileTransfer.PauseTask(ctx, in.TaskId)
	return nil, nil
}

func (h *Handler) CancelSendFile(ctx context.Context, in *CancelSendFileInput) (*CancelSendFileOutput, error) {
	h.fileTransfer.CancelTask(ctx, in.TaskId)
	return nil, nil
}

func (h *Handler) ResumeSendFile(ctx context.Context, in *ResumeSendFileInput) (*ResumeSendFileOutput, error) {
	h.fileTransfer.ResumeTask(ctx, in.TaskId)
	return nil, nil
}

func (h *Handler) SendingTaskList(ctx context.Context, in *ResumeSendFileInput) (*ResumeSendFileOutput, error) {
	h.fileTransfer.ResumeTask(ctx, in.TaskId)
	return nil, nil
}
