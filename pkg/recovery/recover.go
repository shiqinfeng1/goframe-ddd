package recovery

import (
	"context"
	"os"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

type RecoverFunc func(ctx context.Context, exception error)

func Recovery(ctx context.Context, recoverFunc RecoverFunc) {
	if exception := recover(); exception != nil {
		if recoverFunc != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				recoverFunc(ctx, v)
			}
			recoverFunc(ctx, gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception))
		}
		os.Exit(1)
	}
}
