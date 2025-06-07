package application

import (
	"context"

	natsclient "github.com/shiqinfeng1/goframe-ddd/pkg/pubsub/nats"
)

type backend struct {
	pointDataSet PointdataService
	jetStream    JetstreamService
	auth         AuthService
	ncfact       natsclient.Factory
}

func (s *backend) PointDataSet() PointdataService {
	return s.pointDataSet
}
func (s *backend) JetStream() JetstreamService {
	return s.jetStream
}
func (s *backend) Auth() AuthService {
	return s.auth
}
func (s *backend) NatsConnFact() natsclient.Factory {
	return s.ncfact
}

// New 一个DDD的应用层
func New(ctx context.Context, pdsSrv PointdataService, auth AuthService, jsSrv JetstreamService, ncfact natsclient.Factory) Service {
	return &backend{
		pointDataSet: pdsSrv,
		jetStream:    jsSrv,
		auth:         auth,
		ncfact:       ncfact,
	}
}
