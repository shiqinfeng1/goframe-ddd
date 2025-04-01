package application

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/errors"
	pkgnats "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
	"golang.org/x/time/rate"
)

func (h *Application) DeleteStream(ctx context.Context, in *DeleteStreamInput) error {

	client := pkgnats.New(
		g.Cfg().MustGet(ctx, "nats.serverAddr").String(),
		nats.Name("Delete Stream Client"),
	)
	defer client.Close(ctx)

	if err := client.Connect(ctx); err != nil {
		return errors.ErrNatsConnectFail(err)
	}
	defer client.Close(ctx)

	js, err := client.JetStream()
	if err != nil {
		return errors.ErrNatsConnectFail(err)
	}
	if err := js.DeleteStream(ctx, in.Name); err != nil {
		return errors.ErrNatsDeleteStreamFail(err)
	}
	return nil
}
func (h *Application) JetStreamInfo(ctx context.Context, in *JetStreamInfoInput) (*JetStreamInfoOutput, error) {

	client := pkgnats.New(
		g.Cfg().MustGet(ctx, "nats.serverAddr").String(),
		nats.Name("Delete Stream Query"))
	defer client.Close(ctx)

	if err := client.Connect(ctx); err != nil {
		return nil, errors.ErrNatsConnectFail(err)
	}
	defer client.Close(ctx)

	js, err := client.JetStream()
	if err != nil {
		return nil, errors.ErrNatsConnectFail(err)
	}

	// 获取 Stream 信息
	stream, err := js.Stream(ctx, in.Name)
	if err != nil {
		if gerror.Is(err, jetstream.ErrStreamNotFound) {
			return nil, errors.ErrNatsNotFooundStream(in.Name)
		}
		return nil, errors.ErrNatsStreamFail(err)
	}
	si, err := stream.Info(ctx)
	if err != nil {
		return nil, errors.ErrNatsStreamFail(err)
	}
	var cis []*jetstream.ConsumerInfo
	for consumer := range stream.ListConsumers(ctx).Info() {
		cis = append(cis, consumer)
	}
	return &JetStreamInfoOutput{
		StreamInfo:    si,
		ConsumerInfos: cis,
	}, nil
}

// 基准测试无业务逻辑处理，不需要domian层参与，因此直接在app层实现测试逻辑
func (h *Application) PubSubBenchmark(ctx context.Context, in *PubSubBenchmarkInput) error {

	// 再运行发布者，一个发布者一个连接
	for j := range len(in.Subjects) {
		go runPublisher(in.Subjects[j], in.MsgSize)
	}
	for j := range len(in.JsSubjects) {
		go runStreamPublisher(in.JsSubjects[j], in.MsgSize)
	}
	return nil
}

var pubOnce sync.Once
var pubClient *pkgnats.Client

func runPublisher(subj string, msgSize int) {
	ctx := gctx.New()
	// 创建一个发布客户端
	pubOnce.Do(func() {
		pubClient = pkgnats.New(
			g.Cfg().MustGet(ctx, "nats.serverAddr").String(),
			nats.Name("NATS Benchmark"),
		)
		if err := pubClient.Connect(ctx,
			pkgnats.WithJsManager(pkgnats.NewJsSubMgr()),
			pkgnats.WithSubManager(pkgnats.NewSubMgr()),
		); err != nil {
			return
		}
	})
	defer pubClient.Close(ctx)

	var msg []byte
	if msgSize > 0 {
		msg = []byte(grand.Letters(msgSize))
	}
	defer func() {
		conn, err := pubClient.Conn()
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		conn.Flush()
	}()
	cancel := make(chan struct{})
	oneDay := 24 * time.Hour
	time.AfterFunc(oneDay, func() {
		close(cancel)
	})
	// 限速每秒发送50个
	limiter := rate.NewLimiter(rate.Limit(1), 1)
	for {
		select {
		case <-cancel:
			return
		default:
			limiter.Wait(ctx)
			g.Log().Debugf(ctx, "pub msg:%s", msg)
			pubClient.Publish(ctx, subj, msg)
		}
	}
}

var jspubOnce sync.Once
var jspubClient *pkgnats.Client

func runStreamPublisher(subj string, msgSize int) {
	ctx := gctx.New()

	// 创建一个发布客户端
	jspubOnce.Do(func() {
		jspubClient = pkgnats.New(
			g.Cfg().MustGet(ctx, "nats.serverAddr").String(),
			nats.Name("NATS Benchmark"),
		)
		if err := jspubClient.Connect(ctx,
			pkgnats.WithJsManager(pkgnats.NewJsSubMgr()),
			pkgnats.WithSubManager(pkgnats.NewSubMgr()),
		); err != nil {
			return
		}
	})
	defer jspubClient.Close(ctx)

	var msg []byte
	if msgSize > 0 {
		msg = []byte(grand.Letters(msgSize))
	}
	defer func() {
		conn, err := jspubClient.Conn()
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		conn.Flush()
	}()

	cancel := make(chan struct{})
	oneDay := 24 * time.Hour
	time.AfterFunc(oneDay, func() {
		close(cancel)
	})
	// 限速每秒发送1个
	limiter := rate.NewLimiter(rate.Limit(1), 1)
	for {
		select {
		case <-cancel:
			return
		default:
			limiter.Wait(ctx)
			g.Log().Debugf(ctx, "js pub msg:%s", msg)
			jspubClient.JsPublish(ctx, subj, msg)
		}
	}
}
