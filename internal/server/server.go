package server

import (
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	grpc_filemgr "github.com/shiqinfeng1/goframe-ddd/internal/server/grpc/filemgr"
	http_filemgr "github.com/shiqinfeng1/goframe-ddd/internal/server/http/filemgr"
	"github.com/shiqinfeng1/goframe-ddd/internal/server/pubsub"
)

func NewHttpServer() *ghttp.Server {
	// 启动http服务
	s := g.Server()
	if g.Cfg().MustGet(gctx.New(), "pprof").Bool() {
		s.EnablePProf()
	}
	// 设置cors
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.Response.CORSDefault()
		r.Middleware.Next()
	})
	// 健康检查的接口
	s.BindHandler("/health", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{
			"status": "OK",
		})
	})
	// 业务api接口注册
	s.Group("/mgrid", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(
			http_filemgr.NewV1(),
		)
	})
	// 设置openapi的api接口返回格式
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

func NewSubscriptions() *pubsub.SubscriptionManager {
	var subMgr = pubsub.NewSubscriptionManager()
	subMgr.RegisterSubscription("", nil)

	return subMgr
}
