package http

import (
	"context"

	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application/service"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server/http/auth"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server/http/ops"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server/http/pointdata"
	"github.com/shiqinfeng1/goframe-ddd/pkg/dockerctl"
	"github.com/shiqinfeng1/goframe-ddd/pkg/locale"
)

func NewServer(ctx context.Context, logger server.Logger, app application.Service, dockerOps dockerctl.DockerOps) *ghttp.Server {
	// 初始化i18n
	g.I18n().SetPath("config/i18n")
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
	s.Group("/mgrid/v1", func(g *ghttp.RouterGroup) {
		g.Middleware(locale.Locale)
		g.Middleware(service.Auth)
		g.Middleware(ghttp.MiddlewareHandlerResponse)
		g.Bind(
			pointdata.NewV1(logger, app),
			ops.NewV1(logger, app, dockerOps),
			auth.NewV1(logger, app),
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
