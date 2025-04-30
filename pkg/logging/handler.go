package logging

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/text/gstr"
)

// JsonOutputsForLogger is for JSON marshaling in sequence.
type JsonOutputsForLogger struct {
	Time    string `json:"ts"`
	Level   string `json:"lvl"`
	Scope   string `json:"scope"`
	TraceId string `json:"traceId"`
	Msg     []any  `json:"msg"`
}

var LoggingJsonHandler glog.Handler = func(ctx context.Context, in *glog.HandlerInput) {
	jsonForLogger := JsonOutputsForLogger{
		Time:    in.TimeFormat,
		Level:   gstr.Trim(in.LevelFormat, "[]"),
		TraceId: gstr.Trim(in.TraceId, "{}"),
		Scope:   in.Prefix,
		Msg:     in.Values,
	}
	encoder := json.NewEncoder(in.Buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(jsonForLogger)
	in.Next(ctx)
}
