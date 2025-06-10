package pubsub

import (
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server"
	mqttclient "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/mqtt"
	natsclient "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
	"golang.org/x/sync/errgroup"
)

type pubsubConfig struct {
	Mqtt mqttclient.Config
	Nats natsclient.Config
}

type ControllerV1 struct {
	logger            server.Logger
	cfg               pubsubConfig
	natsSubscriptions []string
	natsConsumes      []string
	mqttSubscriptions []string
	group             errgroup.Group
	natsClient        *natsclient.Client
	mqttClient        *mqttclient.Client
	app               application.Service
}

func NewV1(logger server.Logger, app application.Service) *ControllerV1 {
	c := &ControllerV1{
		logger:            logger,
		natsSubscriptions: []string{},
		natsConsumes:      []string{},
		group:             errgroup.Group{},
		app:               app,
	}
	ctx := gctx.New()
	if err := g.Cfg().MustGet(ctx, "nats").Scan(&c.cfg.Nats); err != nil {
		logger.Fatalf(ctx, "get nats config fail:%v", err)
	}
	if err := g.Cfg().MustGet(ctx, "mqtt").Scan(&c.cfg.Mqtt); err != nil {
		logger.Fatalf(ctx, "get mqtt config fail:%v", err)
	}
	c.natsClient = natsclient.New(&c.cfg.Nats, logger)
	// 允许连接云端失败
	c.mqttClient = func() *mqttclient.Client {
		c, err := mqttclient.New(ctx, &c.cfg.Mqtt, logger)
		if err != nil {
			logger.Warningf(ctx, "new mqtt client fail:%v", err)
		}
		return c
	}()
	return c
}

func (c *ControllerV1) Stop(ctx context.Context) error {
	if err := c.natsClient.Close(ctx); err != nil {
		c.logger.Error(ctx, err)
	}
	if err := c.mqttClient.Close(ctx); err != nil {
		c.logger.Error(ctx, err)
	}
	return nil
}
func (c *ControllerV1) NatsTopics() []string {
	return c.natsSubscriptions
}
func (c *ControllerV1) MqttTopics() []string {
	return c.mqttSubscriptions
}
func (c *ControllerV1) StreamTopics() (tpics []string) {
	return c.natsConsumes
}

// 注册topic的处理函数
// 默认一个topic注册一个处理函数， topic支持通配符
// 注意：一个topic起一个协程
func (c *ControllerV1) attachSubscribeHandler(
	ctx context.Context,
	nc *natsclient.Conn,
	subject string,
	handler natsclient.SubscribeFunc) error {
	c.natsSubscriptions = append(c.natsSubscriptions, subject)

	cb := func(ctx context.Context, msg *nats.Msg) error {
		// 处理数据
		out, err := handler(ctx, msg)
		if err != nil {
			return err
		}
		// 转发数据
		// todo: 替换真正的主题
		if err := c.mqttClient.Publish(ctx, subject, out); err != nil {
			return err
		}
		return nil
	}
	err := c.natsClient.SubMsg(ctx, nc, subject, natsclient.SUBASYNC, cb)
	if err != nil {
		return err
	}
	c.logger.Infof(ctx, "subscribe topic ok. topic:%v", subject)
	return nil
}
func (c *ControllerV1) attachConsumeHandler(ctx context.Context, nc *natsclient.Conn, subject string, handler natsclient.ConsumeFunc) {
	c.natsConsumes = append(c.natsConsumes, subject)

	// 创建或更新一个流
	if err := c.natsClient.CreateOrUpdateStream(ctx, nc, c.cfg.Nats.StreamName, c.natsConsumes); err != nil {
		c.logger.Fatalf(ctx, "create stream fail:%v", err)
	}
	// 设置流的消费者
	c.group.Go(func() error {
		c.logger.Debugf(ctx, "start consume stream. topic:%v", subject)
		cb := func(ctx context.Context, msg *jetstream.Msg) error {
			out, err := handler(ctx, msg)
			if err != nil {
				return err
			}
			// todo: 替换真正的主题
			if err := c.mqttClient.Publish(ctx, subject, out); err != nil {
				return err
			}
			return nil
		}
		err := c.natsClient.ConsumeStream(ctx, nc, c.cfg.Nats.StreamName, c.cfg.Nats.ConsumerName, subject, natsclient.JSNEXT, cb)
		if err != nil {
			return gerror.Wrapf(err, "consume stream topic fail:%v", subject)
		}
		c.logger.Infof(ctx, "[goroutine]exit consume stream for topic ok. topic:%v", subject)
		return nil
	})
}
func (c *ControllerV1) attachMqttHandler(ctx context.Context, connForward *natsclient.Conn, topic string, handler mqttclient.SubscribeFunc) error {
	cb := func(ctx context.Context, msg *mqtt.Message) error {
		out, err := handler(ctx, msg)
		if err != nil {
			return err
		}
		// todo: 替换真正的主题
		if err := connForward.PubStream(ctx, (*msg).Topic(), out); err != nil {
			return err
		}
		return nil
	}
	if err := c.mqttClient.Subscribe(ctx, topic, cb); err != nil {
		return err
	}
	return nil
}

// 运行nats订阅客户端
func (c *ControllerV1) Run(ctx context.Context) error {

	// 连接到nats服务端，用于接收消息
	connSub, err := c.app.NatsConnFact().New(ctx, "GoMgridSubscribeClient")
	if err != nil {
		return err
	}
	defer connSub.Close(ctx)
	if err := c.attachSubscribeHandler(ctx, connSub, c.cfg.Nats.Subject1, c.app.PointDataSet().HandleMsg); err != nil {
		return err
	}
	if err := c.attachSubscribeHandler(ctx, connSub, c.cfg.Nats.Subject2, c.app.PointDataSet().HandleMsg); err != nil {
		return err
	}

	// 连接到nats服务端，用于消费流
	connConsume, err := c.app.NatsConnFact().New(ctx, "GoMgridConsumeClient")
	if err != nil {
		return err
	}
	defer connConsume.Close(ctx)
	c.attachConsumeHandler(ctx, connConsume, c.cfg.Nats.JSSubject1, c.app.PointDataSet().HandleStream)
	c.attachConsumeHandler(ctx, connConsume, c.cfg.Nats.JSSubject2, c.app.PointDataSet().HandleStream)

	// 连接到nats服务端，用于监听值变化
	c.group.Go(func() error {
		connWatch, err := c.app.NatsConnFact().New(ctx, "GoMgridWatchClient")
		if err != nil {
			return err
		}
		defer connWatch.Close(ctx)
		if err := c.startWatch(ctx, connWatch); err != nil {
			c.logger.Errorf(ctx, "nats watch fail:%v", err)
			return err
		}
		c.logger.Infof(ctx, "[goroutine]exit nats watch ok")
		return nil
	})

	// 连接到nats服务端，用于转发云端消息给业务服务
	connForward, err := c.app.NatsConnFact().New(ctx, "GoMgridForwardClient")
	if err != nil {
		return err
	}
	defer connForward.Close(ctx)
	if err := c.attachMqttHandler(ctx, connForward, c.cfg.Mqtt.Topic1, c.app.PointDataSet().HandleMqttMsg); err != nil {
		c.logger.Warningf(ctx, "attach mqtt handler fail:%v", err)
	}
	if err := c.attachMqttHandler(ctx, connForward, c.cfg.Mqtt.Topic2, c.app.PointDataSet().HandleMqttMsg); err != nil {
		c.logger.Warningf(ctx, "attach mqtt handler fail:%v", err)
	}

	// 阻塞等待协程退出：订阅连接断开后协程退出
	err = c.group.Wait()
	c.logger.Infof(ctx, "all subscribe & consume exited %v", err)
	return err
}
