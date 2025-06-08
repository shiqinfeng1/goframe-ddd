package logging

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/text/gstr"
	gelf "github.com/robertkowalski/graylog-golang"
)

// glog.HandlerOutputJson
// JsonOutputsForLogger is for JSON marshaling in sequence.
type JsonOutputsForLogger struct {
	Time    string `json:"ts"`
	Level   string `json:"lvl"`
	Scope   string `json:"scope"`
	TraceId string `json:"traceId"`
	Msg     []any  `json:"msg"`
}

var grayLogClient = gelf.New(gelf.Config{
	// GraylogHostname: "graylog-host.com",
	// GraylogPort:     80,
	Connection:      "wan",
	MaxChunkSizeWan: 42,
	MaxChunkSizeLan: 1337,
})

// LoggingGrayLogHandler is an example handler for logging content to remote GrayLog service.
var LoggingGrayLogHandler glog.Handler = func(ctx context.Context, in *glog.HandlerInput) {
	in.Next(ctx)

	// 将日志输出转换为json格式
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
	grayLogClient.Log(in.Buffer.String())
}
