package ops

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	v1 "github.com/shiqinfeng1/goframe-ddd/api/mgrid/http/ops/v1"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application/dto"
	"github.com/shiqinfeng1/goframe-ddd/pkg/errors"
)

func (c *ControllerV1) GetStreamInfo(ctx context.Context, req *v1.GetStreamInfoReq) (res *v1.GetStreamInfoRes, err error) {
	lang := ghttp.RequestFromCtx(ctx).GetCtxVar("lang").String()
	in := &dto.JetStreamInfoIn{
		Name: req.StreamName,
	}
	streams, err := c.app.JetStream().JetStreamInfo(ctx, in)
	if err != nil {
		c.logger.Error(ctx, err)
		return nil, errors.ErrNatsGetStreamInfoFail(lang)
	}
	si := &dto.StreamInfo{
		Name:      streams.StreamInfo.Config.Name,
		Subjects:  streams.StreamInfo.Config.Subjects,
		Retention: streams.StreamInfo.Config.Retention.String(),
		State:     streams.StreamInfo.State,
	}
	var cis []*dto.ConsumerInfo
	for _, ci := range streams.ConsumerInfos {
		cis = append(cis, &dto.ConsumerInfo{
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
	lang := ghttp.RequestFromCtx(ctx).GetCtxVar("lang").String()
	res = &v1.DeleteStreamRes{}
	err = c.app.JetStream().DeleteStream(ctx, &dto.DeleteStreamIn{Name: req.StreamName})
	if err != nil {
		c.logger.Error(ctx, err)
		return nil, errors.ErrNatsDeleteStreamFail(lang)
	}
	return
}
func (c *ControllerV1) GetStreamList(ctx context.Context, req *v1.GetStreamListReq) (res *v1.GetStreamListRes, err error) {
	lang := ghttp.RequestFromCtx(ctx).GetCtxVar("lang").String()
	streams, err := c.app.JetStream().JetStreamList(ctx, &dto.JetStreamListIn{})
	if err != nil {
		c.logger.Error(ctx, err)
		return nil, errors.ErrNatsGetStreamListFail(lang)
	}
	var cis []*dto.StreamInfo
	for _, si := range streams.Streams {
		cis = append(cis, &dto.StreamInfo{
			Name:      si.Config.Name,
			Subjects:  si.Config.Subjects,
			Retention: si.Config.Retention.String(),
			State:     si.State,
		})
	}
	res = &v1.GetStreamListRes{
		Streams: cis,
	}
	return res, nil
}
