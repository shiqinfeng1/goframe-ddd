package application

import (
	"context"
)

type backend struct {
	pointDataSet PointDataSetSrv
	jetStream    JetStreamSrv
}

func (s *backend) PointDataSet() PointDataSetSrv {
	return s.pointDataSet
}
func (s *backend) JetStream() JetStreamSrv {
	return s.jetStream
}

// New 一个DDD的应用层
func New(ctx context.Context, pdsSrv PointDataSetSrv, jsSrv JetStreamSrv) Service {
	return &backend{
		pointDataSet: pdsSrv,
		jetStream:    jsSrv,
	}
}
