package application

import (
	"context"

	natsclient "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
)

type backend struct {
	pointDataSet PointDataSetSrv
	jetStream    JetStreamSrv
	ncfact       natsclient.Factory
}

func (s *backend) PointDataSet() PointDataSetSrv {
	return s.pointDataSet
}
func (s *backend) JetStream() JetStreamSrv {
	return s.jetStream
}
func (s *backend) NatsConnFact() natsclient.Factory {
	return s.ncfact
}

// New 一个DDD的应用层
func New(ctx context.Context, pdsSrv PointDataSetSrv, jsSrv JetStreamSrv, ncfact natsclient.Factory) Service {
	return &backend{
		pointDataSet: pdsSrv,
		jetStream:    jsSrv,
		ncfact:       ncfact,
	}
}
