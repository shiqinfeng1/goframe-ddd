package filemgr

import (
	"context"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	v1 "github.com/shiqinfeng1/goframe-ddd/api/grpc/filemgr/v1"
)

type Controller struct {
	v1.UnimplementedHelloServiceServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterHelloServiceServer(s.Server, &Controller{})
}

func (*Controller) SayHello(ctx context.Context, req *v1.SayHelloRequest) (res *v1.SayHelloResponse, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
