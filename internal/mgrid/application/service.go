package application

import (
	"context"

	natsio "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application/dto"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
)

type JetStreamSrv interface {
	DeleteStream(ctx context.Context, in *dto.DeleteStreamIn) error
	JetStreamInfo(ctx context.Context, in *dto.JetStreamInfoIn) (*dto.JetStreamInfoOut, error)
	JetStreamList(ctx context.Context, in *dto.JetStreamListIn) (*dto.JetStreamListOut, error)
}
type PointDataSetSrv interface {
	HandleMsg(ctx context.Context, msg *natsio.Msg) error
	HandleStream(ctx context.Context, msg *jetstream.Msg) error
}

type Service interface {
	PointDataSet() PointDataSetSrv
	JetStream() JetStreamSrv
	NatsConnFact() natsclient.Factory
}
