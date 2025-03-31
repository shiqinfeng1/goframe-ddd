package pubsub

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shiqinfeng1/goframe-ddd/internal/application"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
	"golang.org/x/sync/errgroup"
)

// 消息处理函数

type ControllerV1 struct {
	subscriptions map[string]pubsub.SubscribeFunc
	group         errgroup.Group
	natsClient    pubsub.Client
	app           *application.Application
}

func NewV1() *ControllerV1 {
	ctx := gctx.New()
	c := &ControllerV1{
		subscriptions: make(map[string]pubsub.SubscribeFunc),
		group:         errgroup.Group{},
		natsClient: nats.New(&nats.Config{
			Server: g.Cfg().MustGet(ctx, "nats.serverAddr").String(),
			// Stream: nats.StreamConfig{
			// 	Name:     g.Cfg().MustGet(ctx, "nats.streamName").String(),
			// 	Subjects: g.Cfg().MustGet(ctx, "nats.subjects").Strings(),
			// 	MaxBytes: g.Cfg().MustGet(ctx, "nats.maxBytes").Int64() * 1024 * 1024 * 1024, //
			// },
			MaxWait:      5 * time.Second,
			ConsumerName: g.Cfg().MustGet(ctx, "nats.consumerName").String(),
		}),
		app: application.App(ctx),
	}
	// 注册topic的处理函数
	subjs := g.Cfg().MustGet(ctx, "nats.subjects").Strings()
	// 默认一个topic注册一个处理函数， topic支持通配符
	// 注意：一个topic起一个协程， 如果对于统一topic但不同的具体subject，如果需要并行处理，也可以为具体的subject起一个处理协程
	c.RegisterSubscription(ctx, subjs[0], c.app.HandleTopic1)
	c.RegisterSubscription(ctx, subjs[1], c.app.HandleTopic2)
	return c
}

func (s *ControllerV1) Stop(ctx context.Context) {
	s.natsClient.Close(ctx)
}
func (s *ControllerV1) Topics() (topics []string) {
	for topic := range s.Subscriptions() {
		topics = append(topics, topic)
	}
	return
}

// 运行nats订阅客户端
func (s *ControllerV1) Run(ctx context.Context) error {
	if err := s.natsClient.Connect(ctx); err != nil {
		return gerror.Wrapf(err, "run subscription manager fail")
	}

	// 使用协程并发订阅消息主题
	for topic, handler := range s.Subscriptions() {
		s.group.Go(func() error {
			g.Log().Debugf(ctx, "start subscribe topic %v ...", topic)
			err := s.StartSubscriber(ctx, topic, handler)
			g.Log().Debugf(ctx, "exit subscribe topic %v %v", topic, err)
			return err
		})
	}
	// 阻塞等待协程退出：订阅连接断开后协程退出
	err := s.group.Wait()
	g.Log().Debugf(ctx, "all subscribe exited:%v", err)
	return err
}

// startSubscriber continuously subscribes to a topic and handles messages using the provided handler.
func (s *ControllerV1) StartSubscriber(ctx context.Context, topic string, handler pubsub.SubscribeFunc) error {
	err := s.natsClient.Subscribe(ctx, topic, handler)
	if err != nil {
		return gerror.Wrapf(err, "subscribe topic %v fail", topic)
	}
	return nil
}

// 注册topic的处理函数
func (s *ControllerV1) RegisterSubscription(ctx context.Context, topic string, handler pubsub.SubscribeFunc) {
	s.subscriptions[topic] = handler
}

// 返回所有注册函数
func (s *ControllerV1) Subscriptions() map[string]pubsub.SubscribeFunc {
	return s.subscriptions
}
