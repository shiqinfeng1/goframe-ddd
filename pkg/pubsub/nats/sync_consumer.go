package natsclient

import (
	"context"
	"errors"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub"
	"github.com/shiqinfeng1/goframe-ddd/pkg/recovery"
)

// 订阅器管理

type SyncConsumer struct {
	logger   pubsub.Logger
	nc       *nats.Conn
	js       *JetStreamWrapper
	msgIters map[string]map[string]jetstream.MessagesContext // stream:consumer:
	ctx      context.Context
	cancel   context.CancelFunc
	errChan  chan error
}

func NewSyncConsumer(logger pubsub.Logger, f Factory) (*SyncConsumer, error) {
	// 创建上下文用于优雅关闭
	ctx, cancel := context.WithCancel(gctx.New())

	// 连接NATS
	nc, err := f.New(ctx, "GoMgridSyncConsumeClient")
	if err != nil {
		cancel()
		return nil, err
	}
	js, _ := jetstream.New(nc)
	jswapper := NewJetStreamWrapper(logger, js)

	return &SyncConsumer{
		logger:   logger,
		nc:       nc,
		js:       jswapper,
		msgIters: make(map[string]map[string]jetstream.MessagesContext),
		ctx:      ctx,
		cancel:   cancel,
		errChan:  make(chan error, 100),
	}, nil
}

func (s *SyncConsumer) Consume(stream, consumer, subject string, handler func(ctx context.Context, msg jetstream.Msg) error) error {

	cons, err := s.js.CreateOrUpdateConsumer(s.ctx, stream, consumer, subject)
	if err != nil {
		return err
	}
	if _, ok := s.msgIters[stream]; !ok {
		s.msgIters[stream] = make(map[string]jetstream.MessagesContext)
	}
	// 创建持久化消费者
	iter, err := cons.Messages(jetstream.PullMaxMessages(1))
	if err != nil {
		return gerror.Wrap(err, "get consumer msg iter fail")
	}
	s.msgIters[stream][consumer] = iter

	go func() {
		s.processMessages(iter, handler)
	}()
	s.logger.Infof(s.ctx, "sync consumer for subject=%v ok", subject)
	return nil
}
func (s *SyncConsumer) processMessages(iter jetstream.MessagesContext, handler func(ctx context.Context, msg jetstream.Msg) error) {
	for {
		msg, err := iter.Next()
		if err != nil {
			if errors.Is(err, jetstream.ErrMsgIteratorClosed) {
				s.errChan <- err
				return
			}
			time.Sleep(consumeMessageDelay)
			s.errChan <- err
			continue
		}
		err = func() error {
			defer recovery.Recovery(s.ctx, func(ctx context.Context, exception error) {
				s.logger.Errorf(ctx, "sync consumer: panic in handler: \n%v", exception)
			})
			if err := handler(s.ctx, msg); err != nil {
				time.Sleep(consumeMessageDelay)
				return err
			}
			return nil
		}()
		if err != nil {
			s.errChan <- err
			continue
		}
		// 处理完成
		if err := msg.Ack(); err != nil {
			s.errChan <- err
		}
	}
}
func (s *SyncConsumer) Run() {
	go func() {
		for err := range s.errChan {
			s.logger.Errorf(s.ctx, "sync consumer error: %v", err)
		}
	}()
	s.logger.Infof(s.ctx, "sync comsuner start running...")
	<-s.ctx.Done()
}

func (s *SyncConsumer) Shutdown() {
	s.logger.Infof(s.ctx, "sync consumer: starting shutdown...")

	// 取消所有订阅
	for stream, v := range s.msgIters {
		for consumer, iter := range v {
			iter.Drain()
			s.js.DeleteConsumer(s.ctx, stream, consumer)
		}
	}
	s.cancel()
	close(s.errChan)
	s.nc.Close()
	s.logger.Infof(s.ctx, "sync consumer: shutdown complete")
}
