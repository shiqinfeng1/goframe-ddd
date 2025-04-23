package server

import (
	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shiqinfeng1/goframe-ddd/internal/filexfer/server/http/filemgr"
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
	s.Group("/filexfer", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(
			filemgr.NewV1(),
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
