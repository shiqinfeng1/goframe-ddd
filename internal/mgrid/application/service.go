package application

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application/dto"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/domain/entity"
	natsclient "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
)

type JetstreamService interface {
	DeleteStream(ctx context.Context, in *dto.DeleteStreamIn) error
	JetStreamInfo(ctx context.Context, in *dto.JetStreamInfoIn) (*dto.JetStreamInfoOut, error)
	JetStreamList(ctx context.Context, in *dto.JetStreamListIn) (*dto.JetStreamListOut, error)
}
type AuthService interface {
	RequestSendVerifyCode(ctx context.Context, email, mobilePhone string) error
	ResetPassword(ctx context.Context, verifyCode, newPassword string) error
	CreateUser(ctx context.Context, in *dto.CreateUserIn) error
	UserIsExisted(ctx context.Context, username, mobilePhone, email string) (exist bool, err error)
	Login(ctx context.Context, user *entity.User) (*dto.Token, error)
	VerifyCredentials(ctx context.Context, lang, username, plainPassword string) (*entity.User, error)
	RefreshToken(ctx context.Context) (*dto.Token, error)
	Logout(ctx context.Context) error
}
type PointdataService interface {
	HandleMsg(ctx context.Context, msg *nats.Msg) ([]byte, error)
	HandleStream(ctx context.Context, msg jetstream.Msg) ([]byte, error)
	HandleMqttMsg(ctx context.Context, msg mqtt.Message) ([]byte, error)
}

type Service interface {
	PointDataSet() PointdataService
	JetStream() JetstreamService
	Auth() AuthService
	NatsConnFact() natsclient.Factory
}
