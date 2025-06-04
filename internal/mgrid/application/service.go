package application

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	natsio "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application/dto"
	natsclient "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
)

type JetStreamSrv interface {
	DeleteStream(ctx context.Context, in *dto.DeleteStreamIn) error
	JetStreamInfo(ctx context.Context, in *dto.JetStreamInfoIn) (*dto.JetStreamInfoOut, error)
	JetStreamList(ctx context.Context, in *dto.JetStreamListIn) (*dto.JetStreamListOut, error)
}
type PointDataSetSrv interface {
	HandleMsg(ctx context.Context, msg *natsio.Msg) ([]byte, error)
	HandleStream(ctx context.Context, msg *jetstream.Msg) ([]byte, error)
	HandleMqttMsg(ctx context.Context, msg *mqtt.Message) ([]byte, error)
}

type Service interface {
	PointDataSet() PointDataSetSrv
	JetStream() JetStreamSrv
	NatsConnFact() natsclient.Factory
}
