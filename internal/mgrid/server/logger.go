package server

import "context"

type Logger interface {
	Errorf(ctx context.Context, format string, v ...interface{})
	Debugf(ctx context.Context, format string, v ...interface{})
	Infof(ctx context.Context, format string, v ...interface{})
	Warningf(ctx context.Context, format string, v ...interface{})

	Error(ctx context.Context, v ...interface{})
	Debug(ctx context.Context, v ...interface{})
	Info(ctx context.Context, v ...interface{})
	Warning(ctx context.Context, v ...interface{})
}
