package pointdata

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/shiqinfeng1/goframe-ddd/api/http/pointdata/v1"
	"github.com/shiqinfeng1/goframe-ddd/internal/application"
	"github.com/shiqinfeng1/goframe-ddd/pkg/utils"
)

const (
	DefaultMessageSize = 128
)

func (c *ControllerV1) PubSubBenchmark(ctx context.Context, req *v1.PubSubBenchmarkReq) (res *v1.PubSubBenchmarkRes, err error) {

	in := &application.PubSubBenchmarkInput{
		MsgSize:      req.MsgSize,
		StreamName:   g.Cfg().MustGet(ctx, "nats.streamName").String(),
		ConsumerName: g.Cfg().MustGet(ctx, "nats.consumerName").String(),
	}
	// 将范围主题展开, 例如：pubsub.station.1.IED.1~50.* 转换成：pubsub.station.1.IED.1.* pubsub.station.1.IED.2.*  ...
	subjects := g.Cfg().MustGet(ctx, "nats.subjects").Strings()
	exsubs := utils.ExpandSubjectRange(strings.TrimSuffix(subjects[0], ">") + "point1~100")
	in.Subjects = append(in.Subjects, exsubs...)
	jssubjects := g.Cfg().MustGet(ctx, "nats.jsSubjects").Strings()
	exjssubs := utils.ExpandSubjectRange(strings.TrimSuffix(jssubjects[0], ">") + "IED.1~50.point.1~2")
	in.JsSubjects = append(in.JsSubjects, exjssubs...)

	if req.MsgSize == 0 {
		in.MsgSize = DefaultMessageSize
	}

	err = c.app.PubSubBenchmark(ctx, in)
	return
}
func (c *ControllerV1) GetStreamInfo(ctx context.Context, req *v1.GetStreamInfoReq) (res *v1.GetStreamInfoRes, err error) {
	in := &application.JetStreamInfoInput{
		Name: req.StreamName,
	}
	streams, err := c.app.JetStreamInfo(ctx, in)
	if err != nil {
		return &v1.GetStreamInfoRes{}, err
	}
	si := &application.StreamInfo{
		Subjects:  streams.StreamInfo.Config.Subjects,
		Retention: streams.StreamInfo.Config.Retention.String(),
		State:     streams.StreamInfo.State,
	}
	var cis []*application.ConsumerInfo
	for _, ci := range streams.ConsumerInfos {
		cis = append(cis, &application.ConsumerInfo{
			Name:           ci.Name,
			Durable:        ci.Config.Durable,
			Description:    ci.Config.Description,
			DeliverPolicy:  ci.Config.DeliverPolicy.String(),
			AckPolicy:      ci.Config.AckPolicy.String(),
			FilterSubject:  ci.Config.FilterSubject,
			FilterSubjects: ci.Config.FilterSubjects,
			NumAckPending:  ci.NumAckPending,
			NumRedelivered: ci.NumRedelivered,
			NumWaiting:     ci.NumWaiting,
			NumPending:     ci.NumPending,
		})
	}
	res = &v1.GetStreamInfoRes{
		StreamInfo:    si,
		ConsumerInfos: cis,
	}
	return res, nil
}
func (c *ControllerV1) DeleteStream(ctx context.Context, req *v1.DeleteStreamReq) (res *v1.DeleteStreamRes, err error) {
	err = c.app.DeleteStream(ctx, &application.DeleteStreamInput{Name: req.StreamName})
	return
}
