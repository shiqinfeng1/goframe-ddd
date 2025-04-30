package application

import "context"

type Logger interface {
	Errorf(ctx context.Context, format string, v ...interface{})
	Debugf(ctx context.Context, format string, v ...interface{})
	Infof(ctx context.Context, format string, v ...interface{})
	Warningf(ctx context.Context, format string, v ...interface{})
}
