package ops

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
	v1 "github.com/shiqinfeng1/goframe-ddd/api/http/ops/v1"
)

func (c *ControllerV1) UpgradeApp(ctx context.Context, req *v1.UpgradeAppReq) (res *v1.UpgradeAppRes, err error) {
	res = &v1.UpgradeAppRes{}
	go func() {
		gproc.ShellExec(gctx.New(), `supervisorctl restart `+req.AppName)
	}()
	g.Log().Debugf(ctx, "upgrade app '%v' ...", req.AppName)
	return
}
func (c *ControllerV1) UpgradeImage(ctx context.Context, req *v1.UpgradeImageReq) (res *v1.UpgradeImageRes, err error) {
	res = &v1.UpgradeImageRes{}
	go func() {
		nctx := gctx.New()
		if err := c.dockerOps.ComposeUp(nctx, req.Version); err != nil {
			g.Log().Debugf(nctx, "exec docker compose up fail':%v", err)
		}
	}()
	g.Log().Debugf(ctx, "upgrade image from '' to '%v' ...", req.Version)
	return
}
