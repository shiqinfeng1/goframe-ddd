package application

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
)

func (h *Application) UpgradeApp(ctx context.Context, in *UpgradeAppInput) error {
	go func() {
		r, err := gproc.ShellExec(gctx.New(), `supervisorctl restart `+in.AppName)
		g.Log().Info(ctx, "result:", r, err)
	}()

	return nil
}
