package pointdata

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/shiqinfeng1/goframe-ddd/api/mgrid/http/pointdata/v1"
)

func (c *ControllerV1) PubSubBenchmark(ctx context.Context, req *v1.PubSubBenchmarkReq) (res *v1.PubSubBenchmarkRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
func (c *ControllerV1) StreamSend(ctx context.Context, req *v1.StreamSendReq) (res *v1.StreamSendRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
func (c *ControllerV1) GetStreamInfo(ctx context.Context, req *v1.GetStreamInfoReq) (res *v1.GetStreamInfoRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
func (c *ControllerV1) DeleteStream(ctx context.Context, req *v1.DeleteStreamReq) (res *v1.DeleteStreamRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
