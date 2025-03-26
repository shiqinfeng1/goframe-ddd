package pointdata

import (
	"context"

	v1 "github.com/shiqinfeng1/goframe-ddd/api/http/pointdata/v1"
	"github.com/shiqinfeng1/goframe-ddd/internal/application"
)

const (
	DefaultNumMsgs     = 100000
	DefaultNumPubs     = 1
	DefaultNumSubs     = 1
	DefaultMessageSize = 128
)

func (c *ControllerV1) PubSubBenchmark(ctx context.Context, req *v1.PubSubBenchmarkReq) (res *v1.PubSubBenchmarkRes, err error) {
	in := &application.PubSubBenchmarkInput{
		NumPubs: req.NumPubs,
		NumSubs: req.NumSubs,
		NumMsgs: req.NumMsgs,
		MsgSize: req.MsgSize,
	}
	if req.MsgSize == 0 {
		in.MsgSize = DefaultMessageSize
	}
	if req.NumPubs == 0 {
		in.NumPubs = DefaultNumPubs
	}
	if req.NumSubs == 0 {
		in.NumSubs = DefaultNumSubs
	}
	if req.NumMsgs == 0 {
		in.NumMsgs = DefaultNumMsgs
	}
	err = c.app.PubSubBenchmark(ctx, in)
	return
}
