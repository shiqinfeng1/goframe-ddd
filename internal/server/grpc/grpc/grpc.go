
package grpc

type Controller struct {
	file-mgr.UnimplementedHelloServiceServer
}

func Register(s *grpcx.GrpcServer) {
	file-mgr.RegisterHelloServiceServer(s.Server, &Controller{})
}

func (*Controller) SayHello(ctx context.Context, req *file-mgr.SayHelloRequest) (res *file-mgr.SayHelloResponse, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
