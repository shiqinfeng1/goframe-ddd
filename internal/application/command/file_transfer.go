package command

import (
	"context"

	"github.com/rs/xid"
)

func (f *Handler) StartSendFile(ctx context.Context, in *StartSendFileInput) (*StartSendFileOutput, error) {
	taskId := xid.New().String()

	f.fileTransfer.AddTask(ctx, taskId, in.BaseName, in.Files)
	return nil, nil
}
