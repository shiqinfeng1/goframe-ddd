package pubsub

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shiqinfeng1/goframe-ddd/internal/application"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
	pkgnats "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
	"github.com/shiqinfeng1/goframe-ddd/pkg/utils"
	"golang.org/x/sync/errgroup"
)

// 消息处理函数

type ControllerV1 struct {
	subscriptions   map[string]pubsub.SubscribeFunc
	jsSubscriptions map[string]pubsub.JsSubscribeFunc
	group           errgroup.Group
	natsClient      *pkgnats.Client
	app             *application.Application
}

func NewV1() *ControllerV1 {
	ctx := gctx.New()
	c := &ControllerV1{
		subscriptions:   make(map[string]pubsub.SubscribeFunc),
		jsSubscriptions: make(map[string]pubsub.JsSubscribeFunc),
		group:           errgroup.Group{},
		natsClient: pkgnats.New(
			g.Cfg().MustGet(ctx, "nats.serverAddr").String(),
			"Client For Subscriber",
		),
		app: application.App(ctx),
	}
	// 注册topic的处理函数
	// 默认一个topic注册一个处理函数， topic支持通配符
	// 注意：一个topic起一个协程， 如果对于统一topic但不同的具体subject，如果需要并行处理，也可以为具体的subject起一个处理协程
	subs := g.Cfg().MustGet(ctx, "nats.subjects").Strings()
	exsubs := utils.ExpandSubjectRange(subs[0])
	for _, exsub := range exsubs {
		c.RegisterSubscription(ctx, exsub, c.app.HandleTopic1)
	}

	subjs := g.Cfg().MustGet(ctx, "nats.jsSubjects").Strings()
	exsubjs := utils.ExpandSubjectRange(subjs[0])
	for _, exsub := range exsubjs {
		c.RegisterJsSubscription(ctx, exsub, c.app.HandleTopic2)
	}
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
func (s *ControllerV1) JsTopics() (topics []string) {
	for topic := range s.JsSubscriptions() {
		topics = append(topics, topic)
	}
	return
}

// 运行nats订阅客户端
func (s *ControllerV1) Run(ctx context.Context) error {
	if err := s.natsClient.Connect(ctx,
		pkgnats.WithJsManager(pkgnats.NewJsSubMgr()),
		pkgnats.WithSubManager(pkgnats.NewSubMgr()),
	); err != nil {
		return gerror.Wrapf(err, "run subscription manager fail")
	}

	// 使用协程并发订阅消息主题
	for topic, handler := range s.Subscriptions() {
		s.group.Go(func() error {
			err := s.natsClient.Subscribe(ctx, topic, handler)
			if err != nil {
				return gerror.Wrapf(err, "subscribe topic '%v' fail", topic)
			}
			g.Log().Debugf(ctx, "exit subscribe topic '%v' %v", topic, err)
			return err
		})
	}

	js, _ := s.natsClient.JetStream()
	jsMgr := pkgnats.NewStreamManager(js)
	consumer := g.Cfg().MustGet(ctx, "nats.consumerName").String()
	stream := g.Cfg().MustGet(ctx, "nats.streamName").String()
	if _, err := jsMgr.GetStream(ctx, stream); err == nil {
		g.Log().Warningf(ctx, "stream '%v' is already created", stream)
		if err := jsMgr.DeleteStream(ctx, stream); err != nil {
			return err
		}
	}
	if err := jsMgr.CreateStream(ctx, stream, s.JsTopics()); err != nil {
		return err
	}
	for topic, handler := range s.JsSubscriptions() {
		identity := []string{stream, consumer, topic}
		s.group.Go(func() error {
			g.Log().Debugf(ctx, "start jsSubscribe topic '%v' ...", topic)
			err := s.natsClient.JsSubscribe(ctx, jsMgr, identity, pkgnats.SubTypeJSConsumeNext, handler)
			if err != nil {
				return gerror.Wrapf(err, "jsSubscribe topic '%v' fail", topic)
			}
			g.Log().Debugf(ctx, "exit jsSubscribe topic '%v': %v", topic, err)
			return err
		})
	}
	// 阻塞等待协程退出：订阅连接断开后协程退出
	err := s.group.Wait()
	g.Log().Debugf(ctx, "all subscribe exited: '%v'", err)
	return err
}

// 注册topic的处理函数
func (s *ControllerV1) RegisterSubscription(ctx context.Context, topic string, handler pubsub.SubscribeFunc) {
	_, exist := s.subscriptions[topic]
	if exist {
		g.Log().Warningf(ctx, "subscriber of topic '%v' is already registered handler", topic)
		return
	}
	s.subscriptions[topic] = handler
	g.Log().Infof(ctx, "subscriber of topic '%v' register handler ok", topic)
}
func (s *ControllerV1) RegisterJsSubscription(ctx context.Context, topic string, handler pubsub.JsSubscribeFunc) {
	_, exist := s.jsSubscriptions[topic]
	if exist {
		g.Log().Warningf(ctx, "js consumer of topic '%v' is already registered handler", topic)
		return
	}
	s.jsSubscriptions[topic] = handler
	g.Log().Infof(ctx, "js consumer of topic '%v' register handler ok", topic)
}

// 返回所有注册函数
func (s *ControllerV1) Subscriptions() map[string]pubsub.SubscribeFunc {
	return s.subscriptions
}
func (s *ControllerV1) JsSubscriptions() map[string]pubsub.JsSubscribeFunc {
	return s.jsSubscriptions
}
