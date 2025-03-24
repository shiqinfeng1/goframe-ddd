package pubsub

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
	"golang.org/x/sync/errgroup"
)

// 消息处理函数
type SubscribeFunc func(ctx context.Context, msg pubsub.Message) error

type SubscriptionManager struct {
	subscriptions map[string]SubscribeFunc
	group         errgroup.Group
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscriptions: make(map[string]SubscribeFunc),
		group:         errgroup.Group{},
	}
}

func (s *SubscriptionManager) Stop() {
	s.GetSubscriber().Close()
}

// 运行nats订阅客户端
func (s *SubscriptionManager) Run(ctx context.Context) error {
	// Start subscribers concurrently using go-routines
	for topic, handler := range s.Subscriptions() {
		s.group.Go(func() error {
			return s.StartSubscriber(ctx, topic, handler)
		})
	}

	return s.group.Wait()
}

// 注册topic的处理函数
func (s *SubscriptionManager) RegisterSubscription(topic string, handler SubscribeFunc) {
	s.subscriptions[topic] = handler
}

// 返回所有注册函数
func (s *SubscriptionManager) Subscriptions() map[string]SubscribeFunc {
	return s.subscriptions
}

// startSubscriber continuously subscribes to a topic and handles messages using the provided handler.
func (s *SubscriptionManager) StartSubscriber(ctx context.Context, topic string, handler SubscribeFunc) error {
	for {
		select {
		case <-ctx.Done():
			g.Log().Infof(ctx, "shutting down subscriber for topic %s", topic)
			return nil
		default:
			err := s.handleSubscription(ctx, topic, handler)
			if err != nil {
				g.Log().Errorf(ctx, "error in subscription for topic %s: %v", topic, err)
			}
		}
	}
}

func (s *SubscriptionManager) handleSubscription(ctx context.Context, topic string, handler SubscribeFunc) error {
	msg, err := s.GetSubscriber().Subscribe(ctx, topic)
	if err != nil {
		g.Log().Errorf(ctx, "error while reading from topic %v, err: %v", topic, err.Error())

		return err
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
		g.Log().Errorf(ctx, "error in handler for topic %s: %v", topic, err)

		return nil
	}

	if msg.Committer != nil {
		// commit the message if the subscription function does not return error
		msg.Commit()
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

	g.Log().Errorf(ctx, e)
}
