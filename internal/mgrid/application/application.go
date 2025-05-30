package application

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
)

type backend struct {
	pointDataSet PointDataSetSrv
	jetStream    JetStreamSrv
	ncfact       nats.ConnFactory
}

func (s *backend) PointDataSet() PointDataSetSrv {
	return s.pointDataSet
}
func (s *backend) JetStream() JetStreamSrv {
	return s.jetStream
}
func (s *backend) NatsConnFact() nats.ConnFactory {
	return s.ncfact
}

// New 一个DDD的应用层
func New(ctx context.Context, pdsSrv PointDataSetSrv, jsSrv JetStreamSrv, ncfact nats.ConnFactory) Service {
	return &backend{
		pointDataSet: pdsSrv,
		jetStream:    jsSrv,
		ncfact:       ncfact,
	}
}
