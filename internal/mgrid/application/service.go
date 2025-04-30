package application

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application/dto"
)

type JetStreamSrv interface {
	DeleteStream(ctx context.Context, in *dto.DeleteStreamIn) error
	JetStreamInfo(ctx context.Context, in *dto.JetStreamInfoIn) (*dto.JetStreamInfoOut, error)
	JetStreamList(ctx context.Context, in *dto.JetStreamListIn) (*dto.JetStreamListOut, error)
}
type PointDataSetSrv interface {
	HandleMsg(ctx context.Context, msg *nats.Msg) error
	HandleStream(ctx context.Context, msg *jetstream.Msg) error
}

type Service interface {
	PointDataSet() PointDataSetSrv
	JetStream() JetStreamSrv
}
