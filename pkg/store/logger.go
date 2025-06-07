package store

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/text/gstr"
)

type Logger interface {
	Errorf(ctx context.Context, format string, v ...any)
	Debugf(ctx context.Context, format string, v ...any)
	Infof(ctx context.Context, format string, v ...any)
	Warningf(ctx context.Context, format string, v ...any)
	Fatalf(ctx context.Context, format string, v ...any)
}

type Log struct {
	Type     string `json:"type"`
	Duration int64  `json:"duration"`
	Key      string `json:"key"`
	Value    string `json:"value,omitempty"`
	Bucket   string `json:"bucket,omitempty"`
}

func (l *Log) String() string {
	var description string

	switch {
	case gstr.Contains(l.Type, "GET"):
		description = fmt.Sprintf("Fetching record from bucket '%s' with ID '%s'", l.Value, l.Key)
	case gstr.Contains(l.Type, "SET"):
		description = fmt.Sprintf("Updating record with ID '%s' in bucket '%s'", l.Key, l.Value)
	case gstr.Contains(l.Type, "DELETE"):
		description = fmt.Sprintf("Deleting record from bucket '%s' with ID '%s'", l.Value, l.Key)
	}
	return fmt.Sprintf("%-32s \u001B[38;5;162mNATS\u001B[0m   %8dμs \u001B[38;5;8m%s\u001B[0m\n",
		l.Type,
		l.Duration,
		description)
}
