package pointdata

import (
	"context"

	v1 "github.com/shiqinfeng1/goframe-ddd/api/http/pointdata/v1"
	"github.com/shiqinfeng1/goframe-ddd/internal/application"
)

const (
	DefaultNumMsgs     = 100000
	DefaultNumPubs     = 1
	DefaultNumSubs     = 1
	DefaultMessageSize = 128
	DefaultSubject     = "benchmark-test"
)

func (c *ControllerV1) PubSubBenchmark(ctx context.Context, req *v1.PubSubBenchmarkReq) (res *v1.PubSubBenchmarkRes, err error) {
	in := &application.PubSubBenchmarkInput{
		NumPubs:  req.NumPubs,
		NumSubs:  req.NumSubs,
		NumMsgs:  req.NumMsgs,
		MsgSize:  req.MsgSize,
		Subjects: req.Subjects,
		Typ:      req.Typ,
	}
	if req.MsgSize == 0 {
		in.MsgSize = DefaultMessageSize
	}
	if req.NumPubs == 0 {
		in.NumPubs = DefaultNumPubs
	}
	if req.NumSubs == 0 {
		in.NumSubs = DefaultNumSubs
	}
	if req.NumMsgs == 0 {
		in.NumMsgs = DefaultNumMsgs
	}
	if len(req.Subjects) == 0 {
		in.Subjects = []string{DefaultSubject}
	}
	err = c.app.PubSubBenchmark(ctx, in)
	return
}
func (c *ControllerV1) GetStreamInfo(ctx context.Context, req *v1.GetStreamInfoReq) (res *v1.GetStreamInfoRes, err error) {
	in := &application.PubSubStreamInfoInput{}
	streams, err := c.app.PubSubStreamInfo(ctx, in)
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
