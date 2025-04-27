package application

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application/dto"
)

type JetStreamSrv interface {
	SendStreamForTest(ctx context.Context) error
	DeleteStream(ctx context.Context, in *dto.DeleteStreamIn) error
	JetStreamInfo(ctx context.Context, in *dto.JetStreamInfoIn) (*dto.JetStreamInfoOut, error)
	PubSubBenchmark(ctx context.Context, in *dto.PubSubBenchmarkIn) error
}
type PointDataSetSrv interface {
	HandleTopic1(ctx context.Context, msg *nats.Msg) error
	HandleTopic2(ctx context.Context, msg *jetstream.Msg) error
}

type Service interface {
	PointDataSet() PointDataSetSrv
	JetStream() JetStreamSrv
}
