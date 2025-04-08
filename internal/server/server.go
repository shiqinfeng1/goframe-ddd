package server

import (
	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	grpc_filemgr "github.com/shiqinfeng1/goframe-ddd/internal/server/grpc/filemgr"
	http_filemgr "github.com/shiqinfeng1/goframe-ddd/internal/server/http/filemgr"
	"github.com/shiqinfeng1/goframe-ddd/internal/server/http/ops"
	"github.com/shiqinfeng1/goframe-ddd/internal/server/http/pointdata"
	"github.com/shiqinfeng1/goframe-ddd/internal/server/pubsub"
)

func NewHttpServer() *ghttp.Server {
	ctx := gctx.New()
	// 启动http服务
	s := g.Server()
	if g.Cfg().MustGet(ctx, "pprof").Bool() {
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
	// 服务监控指标输出接口注册
	s.BindHandler("/metrics", otelmetric.PrometheusHandler)

	// 业务api接口注册
	s.Group("/mgrid", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		handle := []any{
			pointdata.NewV1(),
			ops.NewV1(),
		}
		if g.Cfg().MustGet(ctx, "filemgr.enable").Bool() {
			handle = append(handle, http_filemgr.NewV1())
		}
		group.Bind(
			handle...,
		)
	})

	// 使能管理页面
	// s.EnableAdmin()
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

func NewSubscriptions() *pubsub.ControllerV1 {
	subMgr := pubsub.NewV1()
	return subMgr
}
