package service

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application/dto"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/domain/repository"
	"github.com/shiqinfeng1/goframe-ddd/pkg/errors"
	pkgnats "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
	"github.com/shiqinfeng1/goframe-ddd/pkg/utils"
	"golang.org/x/time/rate"
)

type JetStreamMgr struct {
	repo repository.Repository
}

// NewFileSendQueue 创建一个新的文件发送队列
func NeJetStreamService(_ context.Context, repo repository.Repository) application.JetStreamSrv {
	return &JetStreamMgr{
		repo: repo,
	}
}

func (app *JetStreamMgr) SendStreamForTest(ctx context.Context) error {
	jssubjects := g.Cfg().MustGet(ctx, "nats.jsSubjects").Strings()
	exjssubs := utils.ExpandSubjectRange(strings.TrimSuffix(jssubjects[0], ">") + "IED.1~50.point.1~2")
	for _, j := range exjssubs {
		go runStreamPublisherToRemote(j)
	}
	return nil
}
func (app *JetStreamMgr) DeleteStream(ctx context.Context, in *dto.DeleteStreamIn) error {

	client := pkgnats.New(
		g.Cfg().MustGet(ctx, "nats.serverAddr").String(),
		"Delete Stream Client",
	)

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
func (app *JetStreamMgr) JetStreamInfo(ctx context.Context, in *dto.JetStreamInfoIn) (*dto.JetStreamInfoOut, error) {

	client := pkgnats.New(
		g.Cfg().MustGet(ctx, "nats.serverAddr").String(),
		"Stream Query Client",
	)

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
	return &dto.JetStreamInfoOut{
		StreamInfo:    si,
		ConsumerInfos: cis,
	}, nil
}

// 基准测试无业务逻辑处理，不需要domian层参与，因此直接在app层实现测试逻辑
func (app *JetStreamMgr) PubSubBenchmark(ctx context.Context, in *dto.PubSubBenchmarkIn) error {

	// 再运行发布者，一个发布者一个连接
	for j := range len(in.Subjects) {
		go runPublisher(in.Subjects[j])
	}
	for j := range len(in.JsSubjects) {
		go runStreamPublisher(in.JsSubjects[j])
	}
	return nil
}

var defaultDelay = 24 * time.Hour
var limiter = rate.NewLimiter(rate.Limit(250000), 50)
var jslimiter = rate.NewLimiter(rate.Limit(100), 100)

var pubOnce sync.Once
var jspubOnce sync.Once
var pubClient *pkgnats.Client
var jspubClient *pkgnats.Client

var defaultMessageSize = 1 * 1024 * 1024
var msg = []byte(grand.Letters(defaultMessageSize))
var jsmsg = []byte(grand.Letters(defaultMessageSize))

func runPublisher(subj string) {
	ctx := gctx.New()
	var cancel chan struct{}
	// 创建一个发布客户端
	pubOnce.Do(func() {
		pubClient = pkgnats.New(
			g.Cfg().MustGet(ctx, "nats.serverAddr").String(),
			"Benchmark Publisher",
		)
		if err := pubClient.Connect(ctx,
			pkgnats.WithJsMgr(pkgnats.NewJsSub()),
			pkgnats.WithSubMgr(pkgnats.NewSub()),
		); err != nil {
			return
		}
		cancel = make(chan struct{})
		time.AfterFunc(defaultDelay, func() {
			close(cancel)
		})
	})
	defer func() {
		conn, err := pubClient.Conn()
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		conn.Flush()
		pubClient.Close(ctx)
	}()

	for {
		select {
		case <-cancel:
			return
		default:
			limiter.WaitN(ctx, 50)
			// g.Log().Debugf(ctx, "pub msg: subject='%v' %s...", subj, msg[:8])
			for range 50 {
				pubClient.Publish(ctx, subj, msg)
			}
		}
	}
}

func runStreamPublisherToRemote(subj string) {
	ctx := gctx.New()
	var cancel chan struct{}

	// 创建一个发布客户端
	jspubOnce.Do(func() {
		jspubClient = pkgnats.New(
			"nats://10.17.14.35:4222",
			"Benchmark JsPublisher",
		)
		if err := jspubClient.Connect(ctx,
			pkgnats.WithJsMgr(pkgnats.NewJsSub()),
			pkgnats.WithSubMgr(pkgnats.NewSub()),
		); err != nil {
			g.Log().Error(ctx, err)
			return
		}

		cancel = make(chan struct{})
		time.AfterFunc(defaultDelay, func() {
			close(cancel)
		})
	})

	defer func() {
		conn, err := jspubClient.Conn()
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		conn.Flush()
		jspubClient.Close(ctx)
	}()

	for {
		select {
		case <-cancel:
			return
		default:
			jslimiter.Wait(ctx)
			// g.Log().Debugf(ctx, "js pub msg: subject='%v' (%v)%s...", subj, len(msg), msg[:8])
			jspubClient.JsPublish(ctx, subj, jsmsg)
		}
	}
}
func runStreamPublisher(subj string) {
	ctx := gctx.New()
	var cancel chan struct{}

	// 创建一个发布客户端
	jspubOnce.Do(func() {
		jspubClient = pkgnats.New(
			g.Cfg().MustGet(ctx, "nats.serverAddr").String(),
			"Benchmark JsPublisher",
		)
		if err := jspubClient.Connect(ctx,
			pkgnats.WithJsMgr(pkgnats.NewJsSub()),
			pkgnats.WithSubMgr(pkgnats.NewSub()),
		); err != nil {
			return
		}

		cancel = make(chan struct{})
		time.AfterFunc(defaultDelay, func() {
			close(cancel)
		})
	})

	defer func() {
		conn, err := jspubClient.Conn()
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
		conn.Flush()
		jspubClient.Close(ctx)
	}()

	for {
		select {
		case <-cancel:
			return
		default:
			jslimiter.Wait(ctx)
			// g.Log().Debugf(ctx, "js pub msg: subject='%v' %s...", subj, msg[:8])
			jspubClient.JsPublish(ctx, subj, jsmsg)
		}
	}
}
