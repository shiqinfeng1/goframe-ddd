package server

import (
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	grpc_filemgr "github.com/shiqinfeng1/goframe-ddd/internal/server/grpc/filemgr"
	http_filemgr "github.com/shiqinfeng1/goframe-ddd/internal/server/http/filemgr"
)

func NewHttpServer() *ghttp.Server {
	// 启动http服务
	s := g.Server()
	if g.Cfg().MustGet(gctx.New(), "pprof").Bool() {
		s.EnablePProf()
	}
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.Response.CORSDefault()
		r.Middleware.Next()
	})

	s.Group("/mgrid", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(
			http_filemgr.NewV1(),
		)
	})
	oai := s.GetOpenApi()
	oai.Config.CommonResponse = ghttp.DefaultHandlerResponse{}
	oai.Config.CommonResponseDataField = `Data`
	return s
}

func NewGrpcServer() *grpcx.GrpcServer {
	s := grpcx.Server.New()
	grpc_filemgr.Register(s)
	return s
}
