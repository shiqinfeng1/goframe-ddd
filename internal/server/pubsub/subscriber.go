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
type SubscribeFunc func(ctx context.Context, msg *pubsub.Message) error

type ControllerV1 struct {
	subscriptions map[string]SubscribeFunc
	group         errgroup.Group
	subClient     pubsub.Client
	app           *application.Application
}

func NewV1() *ControllerV1 {
	ctx := gctx.New()
	c := &ControllerV1{
		subscriptions: make(map[string]SubscribeFunc),
		group:         errgroup.Group{},
		subClient: nats.New(&nats.Config{
			Server: g.Cfg().MustGet(ctx, "nats.serverAddr").String(),
			Stream: nats.StreamConfig{
				Name:     g.Cfg().MustGet(ctx, "nats.streamName").String(),
				Subjects: g.Cfg().MustGet(ctx, "nats.subjects").Strings(),
			},
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
	s.subClient.Close(ctx)
}
func (s *ControllerV1) Topics() (topics []string) {
	for topic := range s.Subscriptions() {
		topics = append(topics, topic)
	}
	return
}

// 运行nats订阅客户端
func (s *ControllerV1) Run(ctx context.Context) error {
	if err := s.subClient.Connect(ctx); err != nil {
		return gerror.Wrapf(err, "run subscription manager fail")
	}

	if err := s.subClient.CreateTopic(ctx, s.Topics()); err != nil {
		g.Log().Fatal(ctx, err)
	}
	// Start subscribers concurrently using go-routines
	for topic, handler := range s.Subscriptions() {
		s.group.Go(func() error {
			return s.StartSubscriber(ctx, topic, handler)
		})
	}

	return s.group.Wait()
}

// 注册topic的处理函数
func (s *ControllerV1) RegisterSubscription(ctx context.Context, topic string, handler SubscribeFunc) {
	s.subscriptions[topic] = handler
}

// 返回所有注册函数
func (s *ControllerV1) Subscriptions() map[string]SubscribeFunc {
	return s.subscriptions
}

// startSubscriber continuously subscribes to a topic and handles messages using the provided handler.
func (s *ControllerV1) StartSubscriber(ctx context.Context, topic string, handler SubscribeFunc) error {
	for {
		select {
		case <-ctx.Done():
			g.Log().Infof(ctx, "shutting down subscriber for topic %s", topic)
			return nil
		default:
			err := s.handleSubscription(ctx, topic, handler)
			if err != nil {
				g.Log().Errorf(ctx, "error in subscription: %v", err)
			}
		}
	}
}

func (s *ControllerV1) handleSubscription(ctx context.Context, topic string, handler SubscribeFunc) error {
	msg, err := s.subClient.Subscribe(ctx, topic)
	if err != nil {
		return gerror.Wrapf(err, "error while reading from topic %v", topic)
	}

	if msg == nil {
		return nil
	}

	err = func() error {
		ctx := gctx.New()
		defer func() {
			panicRecovery(ctx, recover())
		}()

		return handler(ctx, msg)
	}()
	if err != nil {
		g.Log().Errorf(ctx, "error in handler for topic '%s': %v", topic, err)
		return nil
	}

	if msg.Committer != nil {
		// commit the message if the subscription function does not return error
		msg.Commit(ctx)
	}

	return nil
}

func panicRecovery(ctx context.Context, re any) {
	if re == nil {
		return
	}

	var e string
	switch t := re.(type) {
	case string:
		e = t
	case error:
		e = t.Error()
	default:
		e = "Unknown panic type"
	}

	g.Log().Error(ctx, e)
}
