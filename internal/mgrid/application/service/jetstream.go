package service

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application/dto"
	natsclient "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
)

type jetstreamService struct {
	logger application.Logger
	nc     *nats.Conn
}

func NewJetstreamService(ctx context.Context, logger application.Logger, ncfact natsclient.Factory) application.JetstreamService {
	jsm := &jetstreamService{
		logger: logger,
	}
	var err error
	jsm.nc, err = ncfact.New(ctx, "GoMgridStreamMgrClient")
	if err != nil {
		logger.Errorf(ctx, "%v", err)
		return nil
	}
	return jsm
}

func (jsm *jetstreamService) DeleteStream(ctx context.Context, in *dto.DeleteStreamIn) error {

	jstream, err := jsm.nc.JetStream()
	if err != nil {
		return err
	}
	if err := jstream.DeleteStream(in.Name); err != nil {
		return gerror.Wrapf(err, "delete stream fail: name=%v", in.Name)
	}
	return nil
}

func (jsm *jetstreamService) JetStreamInfo(ctx context.Context, in *dto.JetStreamInfoIn) (*dto.JetStreamInfoOut, error) {

	jstream, err := jetstream.New(jsm.nc)
	if err != nil {
		return nil, err
	}

	// 获取 Stream 信息
	stream, err := jstream.Stream(ctx, in.Name)
	if err != nil {
		return nil, gerror.Wrapf(err, "get stream info fail: name=%v", in.Name)
	}
	si, err := stream.Info(ctx)
	if err != nil {
		return nil, gerror.Wrapf(err, "get stream info fail: name=%v", in.Name)
	}
	var cis []*jetstream.ConsumerInfo
	for consumer := range stream.ListConsumers(ctx).Info() {
		cis = append(cis, consumer)
	}
	return &dto.JetStreamInfoOut{
		StreamInfo:    si,
		ConsumerInfos: cis,
	}, nil
}

func (jsm *jetstreamService) JetStreamList(ctx context.Context, in *dto.JetStreamListIn) (*dto.JetStreamListOut, error) {

	jstream, err := jetstream.New(jsm.nc)
	if err != nil {
		return nil, err
	}

	// 获取 Stream 信息
	lister := jstream.ListStreams(ctx)
	var cis []*jetstream.StreamInfo
	for stream := range lister.Info() {
		cis = append(cis, stream)
	}
	return &dto.JetStreamListOut{
		Streams: cis,
	}, nil
}
