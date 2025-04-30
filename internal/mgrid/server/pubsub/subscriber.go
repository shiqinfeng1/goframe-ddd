package pubsub

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	natsio "github.com/nats-io/nats.go"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
	"github.com/shiqinfeng1/goframe-ddd/pkg/utils"
	"golang.org/x/sync/errgroup"
)

// 消息处理函数

type ControllerV1 struct {
	logger        server.Logger
	subscriptions map[string]nats.SubscribeFunc
	consumes      map[string]nats.ConsumeFunc
	group         errgroup.Group
	natsClient    *nats.Client
	app           application.Service
}

func NewV1(logger server.Logger, app application.Service) *ControllerV1 {
	ctx := gctx.New()
	c := &ControllerV1{
		logger:        logger,
		subscriptions: make(map[string]nats.SubscribeFunc),
		consumes:      make(map[string]nats.ConsumeFunc),
		group:         errgroup.Group{},
		natsClient: nats.New(
			ctx,
			logger,
			g.Cfg().MustGet(ctx, "nats.serverAddr").String(),
			nil,
		),
		app: app,
	}

	return c
}

func (c *ControllerV1) Stop(ctx context.Context) {
	c.natsClient.Close(ctx)
}
func (c *ControllerV1) Topics() (topics []string) {
	for topic := range c.Subscriptions() {
		topics = append(topics, topic)
	}
	return
}
func (c *ControllerV1) StreamTopics() (topics []string) {
	for topic := range c.Consumes() {
		topics = append(topics, topic)
	}
	return
}

// 运行nats订阅客户端
func (c *ControllerV1) Run(ctx context.Context) error {
	// 注册topic的处理函数
	// 默认一个topic注册一个处理函数， topic支持通配符
	// 注意：一个topic起一个协程， 如果对于统一topic但不同的具体subject，如果需要并行处理，也可以为具体的subject起一个处理协程
	subs := g.Cfg().MustGet(ctx, "nats.subjects").Strings()
	exsubs := utils.ExpandSubjectRange(subs[0])
	for _, exsub := range exsubs {
		c.registerSubscribeFunc(ctx, exsub, c.app.PointDataSet().HandleMsg)
	}

	subjs := g.Cfg().MustGet(ctx, "nats.jsSubjects").Strings()
	exsubjs := utils.ExpandSubjectRange(subjs[0])
	for _, exsub := range exsubjs {
		c.registerConsumeFunc(ctx, exsub, c.app.PointDataSet().HandleStream)
	}
	// 连接到nats服务端，用于接收消息
	connSub, err := c.natsClient.NewConn(ctx, natsio.Name("mgrid sub client"))
	if err != nil {
		return err
	}
	defer connSub.Close(ctx)

	// 使用协程并发订阅消息主题
	for topic, handler := range c.Subscriptions() {
		c.group.Go(func() error {
			err := c.natsClient.SubMsg(ctx, connSub, topic, nats.SubTypeSubAsync, handler)
			if err != nil {
				return gerror.Wrapf(err, "subscribe topic '%v' fail", topic)
			}
			c.logger.Debugf(ctx, "exit subscribe for topic '%v' ok", topic)
			return err
		})
	}

	// 连接到nats服务端，用于消费流
	connConsume, err := c.natsClient.NewConn(ctx, natsio.Name("mgrid consume client"))
	if err != nil {
		return err
	}
	defer connConsume.Close(ctx)

	consumer := g.Cfg().MustGet(ctx, "nats.consumerName").String()
	streamName := g.Cfg().MustGet(ctx, "nats.streamName").String()
	if err := c.natsClient.CreateOrUpdateStream(ctx, connConsume, streamName, c.StreamTopics()); err != nil {
		return err
	}
	for topic, handler := range c.Consumes() {
		c.group.Go(func() error {
			c.logger.Debug(ctx, "start consume stream", g.Map{"topic": topic})
			err := c.natsClient.ConsumeStream(ctx, connConsume, streamName, consumer, topic, nats.SubTypeJSConsumeNext, handler)
			if err != nil {
				return gerror.Wrapf(err, "consume stream topic fail:%v", topic)
			}
			c.logger.Debug(ctx, "exit consume stream for topic ok", g.Map{"topic": topic})
			return nil
		})
	}
	// 阻塞等待协程退出：订阅连接断开后协程退出
	err = c.group.Wait()
	c.logger.Infof(ctx, "all subscribe & consume exited %v", err)
	return err
}

// 注册topic的处理函数
func (c *ControllerV1) registerSubscribeFunc(ctx context.Context, topic string, handler nats.SubscribeFunc) {
	_, exist := c.subscriptions[topic]
	if exist {
		c.logger.Warningf(ctx, "subscriber of topic '%v' is already registered handler", topic)
		return
	}
	c.subscriptions[topic] = handler
	c.logger.Infof(ctx, "subscriber of topic '%v' register handler ok", topic)
}

func (c *ControllerV1) registerConsumeFunc(ctx context.Context, topic string, handler nats.ConsumeFunc) {
	_, exist := c.consumes[topic]
	if exist {
		c.logger.Warningf(ctx, "stream consumer of topic '%v' is already registered handler", topic)
		return
	}
	c.consumes[topic] = handler
	c.logger.Infof(ctx, "stream consumer of topic '%v' register handler ok", topic)
}

// 返回所有注册函数
func (c *ControllerV1) Subscriptions() map[string]nats.SubscribeFunc {
	return c.subscriptions
}
func (c *ControllerV1) Consumes() map[string]nats.ConsumeFunc {
	return c.consumes
}
