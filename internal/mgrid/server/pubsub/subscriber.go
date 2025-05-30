package pubsub

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
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
	c := &ControllerV1{
		logger:        logger,
		subscriptions: make(map[string]nats.SubscribeFunc),
		consumes:      make(map[string]nats.ConsumeFunc),
		group:         errgroup.Group{},
		natsClient:    nats.New(logger),
		app:           app,
	}

	return c
}

func (c *ControllerV1) Stop(ctx context.Context) error {
	return c.natsClient.Close(ctx)
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

// 注册topic的处理函数
// 默认一个topic注册一个处理函数， topic支持通配符
// 注意：一个topic起一个协程， 如果对于统一topic但不同的具体subject，如果需要并行处理，也可以为具体的subject起一个处理协程
func (c *ControllerV1) attachSubscribeHandler(ctx context.Context, nc *nats.Conn) error {
	subs := g.Cfg().MustGet(ctx, "nats.subjects").Strings()
	exsubs := utils.ExpandSubjectRange(subs[0])
	for _, exsub := range exsubs {
		_, exist := c.subscriptions[exsub]
		if exist {
			c.logger.Warningf(ctx, "subscriber of topic '%v' is already registered handler", exsub)
			return nil
		}
		c.subscriptions[exsub] = c.app.PointDataSet().HandleMsg
		c.logger.Infof(ctx, "subscriber of topic '%v' register handler ok", exsub)
	}

	// 注册主题处理函数, 一个主题一个协程
	for topic, handler := range c.Subscriptions() {
		c.group.Go(func() error {
			err := c.natsClient.SubMsg(ctx, nc, topic, nats.SubTypeSubAsync, handler)
			if err != nil {
				return gerror.Wrapf(err, "subscribe topic '%v' fail", topic)
			}
			c.logger.Debugf(ctx, "exit subscribe for topic '%v' ok", topic)
			return err
		})
	}
	return nil
}
func (c *ControllerV1) attachConsumeHandler(ctx context.Context, nc *nats.Conn) error {

	subjs := g.Cfg().MustGet(ctx, "nats.jsSubjects").Strings()
	exsubjs := utils.ExpandSubjectRange(subjs[0])
	for _, exsub := range exsubjs {
		_, exist := c.consumes[exsub]
		if exist {
			c.logger.Warningf(ctx, "stream consumer of topic '%v' is already registered handler", exsub)
			return nil
		}
		c.consumes[exsub] = c.app.PointDataSet().HandleStream
		c.logger.Infof(ctx, "stream consumer of topic '%v' register handler ok", exsub)
	}

	consumerName := g.Cfg().MustGet(ctx, "nats.consumerName").String()
	streamName := g.Cfg().MustGet(ctx, "nats.streamName").String()

	// 创建一个流
	if err := c.natsClient.CreateOrUpdateStream(ctx, nc, streamName, c.StreamTopics()); err != nil {
		return err
	}
	// 设置流的消费者
	for topic, handler := range c.Consumes() {
		c.group.Go(func() error {
			c.logger.Debug(ctx, "start consume stream", g.Map{"topic": topic})
			err := c.natsClient.ConsumeStream(ctx, nc, streamName, consumerName, topic, nats.SubTypeJSConsumeNext, handler)
			if err != nil {
				return gerror.Wrapf(err, "consume stream topic fail:%v", topic)
			}
			c.logger.Debug(ctx, "exit consume stream for topic ok", g.Map{"topic": topic})
			return nil
		})
	}
	return nil
}

// 运行nats订阅客户端
func (c *ControllerV1) Run(ctx context.Context) error {

	// 连接到nats服务端，用于接收消息
	connSub, err := c.app.NatsConnFact().New(ctx, "go-mgrid subscribe client")
	if err != nil {
		return err
	}
	defer connSub.Close(ctx)
	if err := c.attachSubscribeHandler(ctx, connSub); err != nil {
		return nil
	}

	// 连接到nats服务端，用于消费流
	connConsume, err := c.app.NatsConnFact().New(ctx, "go-mgrid consume client")
	if err != nil {
		return err
	}
	defer connConsume.Close(ctx)
	if err := c.attachConsumeHandler(ctx, connConsume); err != nil {
		return nil
	}

	// 连接到nats服务端，用于消费流
	connWatch, err := c.app.NatsConnFact().New(ctx, "go-mgrid consume client")
	if err != nil {
		return err
	}
	defer connWatch.Close(ctx)
	if err := c.startWatch(ctx, connWatch); err != nil {
		return nil
	}

	// 阻塞等待协程退出：订阅连接断开后协程退出
	err = c.group.Wait()
	c.logger.Infof(ctx, "all subscribe & consume exited %v", err)
	return err
}

// 返回所有注册函数
func (c *ControllerV1) Subscriptions() map[string]nats.SubscribeFunc {
	return c.subscriptions
}
func (c *ControllerV1) Consumes() map[string]nats.ConsumeFunc {
	return c.consumes
}
