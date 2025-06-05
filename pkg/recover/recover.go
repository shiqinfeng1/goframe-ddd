package recover

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

type RecoverFunc func(ctx context.Context, exception error)

func Recovery(ctx context.Context, recoverFunc RecoverFunc) {
	if exception := recover(); exception != nil {
		if recoverFunc != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				recoverFunc(ctx, v)
			} else {
				recoverFunc(ctx, gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception))
			}
		}
	}
}
