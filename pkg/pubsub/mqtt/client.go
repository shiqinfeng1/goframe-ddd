package mqttclient

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/shiqinfeng1/goframe-ddd/pkg/panic"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	subMessageDelay = 100 * time.Millisecond
	qos             = 1
)

type Config struct {
	BrokerUrl string
	ClientID  string
}

type Client struct {
	config Config
	logger pubsub.Logger
	mqttco *mqtt.ClientOptions
	mqttc  mqtt.Client
	topics []string
}
type SubscribeFunc func(ctx context.Context, msg mqtt.Message) error

func New(ctx context.Context, config Config, logger pubsub.Logger) (*Client, error) {
	u, err := url.Parse(config.BrokerUrl)
	if err != nil {
		return nil, err
	}
	// 创建文件存储（断链时缓存消息）
	store := mqtt.NewFileStore("./mqtt_store")
	if err := gfile.Mkdir("./mqtt_store"); err != nil {
		return nil, gerror.Wrap(err, "init mqtt store fail")
	}
	opts := mqtt.NewClientOptions()
	opts.SetStore(store) // 启用消息缓存
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", u.Hostname(), u.Port()))
	opts.SetClientID(config.ClientID)
	opts.SetUsername(u.User.Username())
	if pwd, has := u.User.Password(); has {
		opts.SetPassword(pwd)
	}
	// 配置 QoS 1 相关参数
	opts.SetAutoReconnect(true)                   // 启用自动重连
	opts.SetMaxReconnectInterval(5 * time.Second) // 最大重连间隔
	opts.SetConnectRetry(true)                    // 连接失败时重试
	opts.SetCleanSession(false)                   // 保持会话状态，用于 QoS 1 消息确认

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		optsr := client.OptionsReader()
		logger.Infof(ctx, "connect to mqtt broker ok: addr=%v clientId=%v", fmt.Sprintf("%v", optsr.Servers()), optsr.ClientID())
	})

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		logger.Warningf(ctx, "mqtt broker connect lost:  %v", err)
	})
	opts.SetReconnectingHandler(func(client mqtt.Client, opts *mqtt.ClientOptions) {
		optsr := client.OptionsReader()
		logger.Warningf(ctx, "reconnect to mqtt broker: addr=%v clientId=%v", fmt.Sprintf("%v", optsr.Servers()), optsr.ClientID())
	})

	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &Client{
		config: config,
		logger: logger,
		mqttco: opts,
		mqttc:  c,
	}, nil
}

func (c *Client) Publish(topic string, message []byte) error {
	if !c.mqttc.IsConnected() {
		return gerror.New("mqtt not connected")
	}
	token := c.mqttc.Publish(topic, qos, false, message)

	if token.WaitTimeout(1*time.Second) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *Client) Subscribe(ctx context.Context, topic string, handler SubscribeFunc) error {

	cb := func(mc mqtt.Client, msg mqtt.Message) {
		defer func() {
			panic.Recovery(ctx, func(ctx context.Context, exception error) {
				c.logger.Errorf(ctx, "panic in mqtt handler: %v", exception)
			})
		}()
		if err := handler(ctx, msg); err != nil {
			c.logger.Errorf(ctx, "mqtt handler: %v", err)
			time.Sleep(subMessageDelay)
			return
		}
	}

	if token := c.mqttc.Subscribe(topic, qos, cb); token.Wait() && token.Error() != nil {
		return gerror.Wrap(token.Error(), "mqtt sub fail")
	}
	c.topics = append(c.topics, topic)
	c.logger.Infof(ctx, "mqtt sub success. topic=%v")
	return nil
}

func (c *Client) Close(ctx context.Context) error {
	// 取消订阅
	for _, topic := range c.topics {
		if token := c.mqttc.Unsubscribe(topic); token.Wait() && token.Error() != nil {
			return gerror.Wrap(token.Error(), "mqtt unsub fail")
		}
	}

	c.logger.Infof(ctx, "mqtt pub/sub closed")

	c.mqttc.Disconnect(2000)
	return nil
}
