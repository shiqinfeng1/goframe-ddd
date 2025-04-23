package server

import (
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/shiqinfeng1/goframe-ddd/internal/demo/server/grpc/hello"
)

func NewGrpcServer() *grpcx.GrpcServer {
	s := grpcx.Server.New()
	hello.Register(s)
	return s
}
