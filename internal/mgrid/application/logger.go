package application

import "context"

type Logger interface {
	Errorf(ctx context.Context, format string, v ...any)
	Debugf(ctx context.Context, format string, v ...any)
	Infof(ctx context.Context, format string, v ...any)
	Warningf(ctx context.Context, format string, v ...any)
}
