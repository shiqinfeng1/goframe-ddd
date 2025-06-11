package pubsub

import (
	"context"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/server"
	mqttclient "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/mqtt"
	natsclient "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
	"github.com/shiqinfeng1/goframe-ddd/pkg/recovery"
)

type pubsubConfig struct {
	Mqtt mqttclient.Config
	Nats natsclient.Config
}

type ControllerV1 struct {
	ctx        context.Context
	logger     server.Logger
	cfg        pubsubConfig
	group      sync.WaitGroup
	natsClient *natsclient.Client
	mqttClient *mqttclient.Client
	app        application.Service
}

func NewV1(logger server.Logger, app application.Service) *ControllerV1 {
	c := &ControllerV1{
		ctx:    gctx.New(),
		logger: logger,
		group:  sync.WaitGroup{},
		app:    app,
	}
	if err := g.Cfg().MustGet(c.ctx, "nats").Scan(&c.cfg.Nats); err != nil {
		logger.Fatalf(c.ctx, "get nats config fail:%v", err)
	}
	if err := g.Cfg().MustGet(c.ctx, "mqtt").Scan(&c.cfg.Mqtt); err != nil {
		logger.Fatalf(c.ctx, "get mqtt config fail:%v", err)
	}
	c.natsClient = natsclient.New(&c.cfg.Nats, logger)
	// 允许连接云端失败
	c.mqttClient = func() *mqttclient.Client {
		mqttc, err := mqttclient.New(c.ctx, &c.cfg.Mqtt, logger)
		if err != nil {
			logger.Warningf(c.ctx, "new mqtt client fail:%v", err)
		}
		return mqttc
	}()
	return c
}

func (c *ControllerV1) Stop() error {
	if err := c.natsClient.Close(c.ctx); err != nil {
		c.logger.Error(c.ctx, err)
	}
	if err := c.mqttClient.Close(c.ctx); err != nil {
		c.logger.Error(c.ctx, err)
	}
	return nil
}

// 注册topic的处理函数
// 默认一个topic注册一个处理函数， topic支持通配符
// 注意：一个topic起一个协程
func (c *ControllerV1) asyncSubscribe(
	factory natsclient.Factory) {

	subscriber, err := c.natsClient.NewAsyncSubscriber(c.logger, factory)
	if err != nil {
		c.logger.Fatal(c.ctx, err)
	}

	// 订阅主题
	err = subscriber.Subscribe(c.cfg.Nats.Subject1, func(ctx context.Context, msg *nats.Msg) error {
		// 处理数据
		out, err := c.app.PointDataSet().HandleMsg(ctx, msg)
		if err != nil {
			return err
		}
		// 转发数据
		// todo: 替换真正的主题
		if err := c.mqttClient.Publish(ctx, msg.Subject, out); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.logger.Fatal(c.ctx, err)
	}
	// 运行订阅者
	subscriber.Run()
}
func (c *ControllerV1) syncSubscribe(
	factory natsclient.Factory) {

	subscriber, err := c.natsClient.NewSyncSubscriber(c.logger, factory)
	if err != nil {
		c.logger.Fatal(c.ctx, err)
	}

	// 订阅主题
	err = subscriber.Subscribe(c.cfg.Nats.Subject2, func(ctx context.Context, msg *nats.Msg) error {
		// 处理数据
		out, err := c.app.PointDataSet().HandleMsg(ctx, msg)
		if err != nil {
			return err
		}
		// 转发数据
		// todo: 替换真正的主题
		if err := c.mqttClient.Publish(ctx, msg.Subject, out); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.logger.Fatal(c.ctx, err)
	}
	// 运行订阅者
	subscriber.Run()
}

func (c *ControllerV1) syncConsume(
	factory natsclient.Factory) {
	nc, err := factory.New(c.ctx, "GoMgridCreateOrUpdateStreamClient")
	if err != nil {
		c.logger.Fatal(c.ctx, err)
	}
	stream := c.cfg.Nats.StreamName
	err = c.natsClient.CreateOrUpdateStream(c.ctx, nc, c.cfg.Nats.StreamName, []string{c.cfg.Nats.JSSubject1, c.cfg.Nats.JSSubject2})
	if err != nil {
		nc.Close()
		c.logger.Fatal(c.ctx, err)
	}
	nc.Close()

	consumer, err := c.natsClient.NewSyncConsumer(c.logger, factory)
	if err != nil {
		c.logger.Fatal(c.ctx, err)
	}

	// 订阅主题
	err = consumer.Consume(stream, c.cfg.Nats.ConsumerName1, c.cfg.Nats.JSSubject1, func(ctx context.Context, msg jetstream.Msg) error {
		// 处理数据
		out, err := c.app.PointDataSet().HandleStream(ctx, msg)
		if err != nil {
			return err
		}
		// 转发数据
		// todo: 替换真正的主题
		if err := c.mqttClient.Publish(ctx, msg.Subject(), out); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.logger.Fatal(c.ctx, err)
	}
	err = consumer.Consume(stream, c.cfg.Nats.ConsumerName2, c.cfg.Nats.JSSubject2, func(ctx context.Context, msg jetstream.Msg) error {
		// 处理数据
		out, err := c.app.PointDataSet().HandleStream(ctx, msg)
		if err != nil {
			return err
		}
		// 转发数据
		// todo: 替换真正的主题
		if err := c.mqttClient.Publish(ctx, msg.Subject(), out); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.logger.Fatal(c.ctx, err)
	}
	// 运行订阅者
	consumer.Run()
}
func (c *ControllerV1) attachMqttHandler(connForward *nats.Conn, topic string, handler mqttclient.SubscribeFunc) error {

	publisher := natsclient.NewPublisher(connForward)
	cb := func(ctx context.Context, msg mqtt.Message) error {
		out, err := handler(ctx, msg)
		if err != nil {
			return err
		}
		// todo: 替换真正的主题
		if err := publisher.PublishStreamMsg(ctx, msg.Topic(), out); err != nil {
			return err
		}
		return nil
	}
	if err := c.mqttClient.Subscribe(c.ctx, topic, cb); err != nil {
		return err
	}
	return nil
}
func (c *ControllerV1) watchKey() error {
	connWatch, err := c.app.NatsConnFact().New(c.ctx, "GoMgridWatchClient")
	if err != nil {
		return err
	}
	defer connWatch.Close()
	if err := c.startWatch(c.ctx, connWatch); err != nil {
		c.logger.Errorf(c.ctx, "nats watch fail:%v", err)
		return err
	}
	c.logger.Infof(c.ctx, "[goroutine]exit nats watch ok")
	return nil
}
func (c *ControllerV1) registerMqttHandler() error {
	// 连接到nats服务端，用于转发云端消息给业务服务
	connForward, err := c.app.NatsConnFact().New(c.ctx, "GoMgridForwardClient")
	if err != nil {
		return err
	}
	defer connForward.Close()
	if err := c.attachMqttHandler(connForward, c.cfg.Mqtt.Topic1, c.app.PointDataSet().HandleMqttMsg); err != nil {
		c.logger.Warningf(c.ctx, "attach mqtt handler fail:%v", err)
	}
	if err := c.attachMqttHandler(connForward, c.cfg.Mqtt.Topic2, c.app.PointDataSet().HandleMqttMsg); err != nil {
		c.logger.Warningf(c.ctx, "attach mqtt handler fail:%v", err)
	}
	c.logger.Infof(c.ctx, "register mqtt handler ok")
	return nil
}

// 运行nats订阅客户端
func (c *ControllerV1) Run() error {
	routine := func(f func() error) error {
		c.group.Add(1)
		defer recovery.Recovery(c.ctx, func(ctx context.Context, exception error) {
			c.group.Done()
			c.logger.Errorf(c.ctx, "panic in controller:\n%v", exception)
		})
		return f()
	}
	go routine(func() error { c.asyncSubscribe(c.app.NatsConnFact()); return nil })
	go routine(func() error { c.syncSubscribe(c.app.NatsConnFact()); return nil })
	go routine(func() error { c.syncConsume(c.app.NatsConnFact()); return nil })
	go routine(func() error { return c.watchKey() })
	if err := c.registerMqttHandler(); err != nil {
		return err
	}
	// 阻塞等待协程退出：订阅连接断开后协程退出
	c.group.Wait()
	c.logger.Infof(c.ctx, "all subscribe & consume exited")
	return nil
}
