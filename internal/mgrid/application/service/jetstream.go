package service

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	natsio "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application/dto"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/domain/repository"
	"github.com/shiqinfeng1/goframe-ddd/pkg/errors"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
)

type JetStreamMgr struct {
	logger application.Logger
	repo   repository.Repository
}

func NeJetStreamService(_ context.Context, logger application.Logger, repo repository.Repository) application.JetStreamSrv {
	return &JetStreamMgr{
		logger: logger,
		repo:   repo,
	}
}

func (js *JetStreamMgr) DeleteStream(ctx context.Context, in *dto.DeleteStreamIn) error {
	conn, err := nats.NewFactory(
		js.logger,
		g.Cfg().MustGet(ctx, "nats.serverAddr").String(),
		nil,
	).New(ctx, natsio.Name("client for delete stream:"+in.Name))
	if err != nil {
		return errors.ErrNatsConnectFail(err)
	}
	defer conn.Close(ctx)

	jstream, err := conn.JetStream()
	if err != nil {
		return errors.ErrNatsConnectFail(err)
	}
	if err := jstream.DeleteStream(ctx, in.Name); err != nil {
		return errors.ErrNatsDeleteStreamFail(err)
	}
	return nil
}

func (js *JetStreamMgr) JetStreamInfo(ctx context.Context, in *dto.JetStreamInfoIn) (*dto.JetStreamInfoOut, error) {

	conn, err := nats.NewFactory(
		js.logger,
		g.Cfg().MustGet(ctx, "nats.serverAddr").String(),
		nil,
	).New(ctx, natsio.Name("client for query stream:"+in.Name))
	if err != nil {
		return nil, errors.ErrNatsConnectFail(err)
	}
	defer conn.Close(ctx)

	jstream, err := conn.JetStream()
	if err != nil {
		return nil, errors.ErrNatsConnectFail(err)
	}

	// 获取 Stream 信息
	stream, err := jstream.Stream(ctx, in.Name)
	if err != nil {
		if gerror.Is(err, jetstream.ErrStreamNotFound) {
			return nil, errors.ErrNatsNotFooundStream(in.Name)
		}
		return nil, errors.ErrNatsStreamFail(err)
	}
	si, err := stream.Info(ctx)
	if err != nil {
		return nil, errors.ErrNatsStreamFail(err)
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

func (js *JetStreamMgr) JetStreamList(ctx context.Context, in *dto.JetStreamListIn) (*dto.JetStreamListOut, error) {

	conn, err := nats.NewFactory(
		js.logger,
		g.Cfg().MustGet(ctx, "nats.serverAddr").String(),
		nil,
	).New(ctx, natsio.Name("client for query streamlist"))
	if err != nil {
		return nil, errors.ErrNatsConnectFail(err)
	}
	defer conn.Close(ctx)

	jstream, err := conn.JetStream()
	if err != nil {
		return nil, errors.ErrNatsConnectFail(err)
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
