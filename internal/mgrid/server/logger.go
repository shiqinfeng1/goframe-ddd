package server

import "context"

type Logger interface {
	Errorf(ctx context.Context, format string, v ...any)
	Debugf(ctx context.Context, format string, v ...any)
	Infof(ctx context.Context, format string, v ...any)
	Warningf(ctx context.Context, format string, v ...any)
	Fatalf(ctx context.Context, format string, v ...any)

	Error(ctx context.Context, v ...any)
	Debug(ctx context.Context, v ...any)
	Info(ctx context.Context, v ...any)
	Warning(ctx context.Context, v ...any)
	Fatal(ctx context.Context, v ...any)
}
