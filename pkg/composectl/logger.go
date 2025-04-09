package composectl

import (
	"context"

	"github.com/docker/compose/v2/pkg/api"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/logrusorgru/aurora"
)

var colorPool = []aurora.Color{
	aurora.BlueFg,
	aurora.GreenFg,
	aurora.CyanFg,
	aurora.MagentaFg,
	aurora.YellowFg,
	aurora.RedFg,
}

type logger struct {
	ctx             context.Context
	containerColors map[string]aurora.Color
}

var _ api.LogConsumer = (*logger)(nil)

func newLogConsumer(ctx context.Context) (*logger, error) {
	return &logger{
		ctx:             ctx,
		containerColors: map[string]aurora.Color{},
	}, nil
}

func (l *logger) colorize(cid string) string {
	color, ok := l.containerColors[cid]
	if ok {
		return aurora.Colorize(cid, color).String()
	}

	color = colorPool[len(l.containerColors)%len(colorPool)]
	l.containerColors[cid] = color

	return aurora.Colorize(cid, color).String()
}

func (l *logger) Log(container, msg string) {
	g.Log().Infof(l.ctx, "docker service [%s] %s", l.colorize(container), msg)
}
func (l *logger) Err(container, msg string) {
	g.Log().Infof(l.ctx, "docker service [%s] err: %s", l.colorize(container), msg)
}

func (l *logger) Status(container, msg string) {
	g.Log().Infof(l.ctx, "docker service [%s] status: %s", l.colorize(container), msg)
}

func (l *logger) Register(container string) {
	g.Log().Infof(l.ctx, "docker service [%s] registered ok", l.colorize(container))
}
