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
