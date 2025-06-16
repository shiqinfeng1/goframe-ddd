package store

import (
	"context"
	"time"

	"github.com/shiqinfeng1/goframe-ddd/pkg/clock"
	"github.com/shiqinfeng1/goframe-ddd/pkg/metrics"
)

type OpType struct {
	name     string
	cntTotal string
	cntSucc  string
}

var (
	GET_VAULE = OpType{
		name:     "GET_VAULE",
		cntTotal: metrics.NatsKVGetTotalCount,
		cntSucc:  metrics.NatsKVGetSuccessCount}
	SET_VAULE = OpType{
		name:     "SET_VAULE",
		cntTotal: metrics.NatsKVSetTotalCount,
		cntSucc:  metrics.NatsKVSetSuccessCount}
	SET_OBJ = OpType{
		name:     "SET_OBJ",
		cntTotal: metrics.NatsObjSetTotalCount,
		cntSucc:  metrics.NatsObjSetSuccessCount}
	GET_OBJ = OpType{
		name:     "GET_OBJ",
		cntTotal: metrics.NatsObjGetTotalCount,
		cntSucc:  metrics.NatsObjGetSuccessCount}
	GET_FILE = OpType{
		name:     "GET_FILE",
		cntTotal: metrics.NatsObjGetTotalCount,
		cntSucc:  metrics.NatsObjGetSuccessCount}
	SET_FILE = OpType{
		name:     "SET_FILE",
		cntTotal: metrics.NatsObjGetTotalCount,
		cntSucc:  metrics.NatsObjGetSuccessCount}
)

func (o OpType) MetricTotalName() string {
	return o.cntTotal
}
func (o OpType) MetricSuccName() string {
	return o.cntSucc
}
func (o OpType) String() string {
	return o.name
}

func SendOperationStats(ctx context.Context, ot OpType, bucket string, key string, f func() error) error {
	start := clock.Now()
	// 统计总数
	metrics.Inc(ctx, ot.MetricTotalName(), key, bucket)
	if err := f(); err != nil {
		return err
	}
	duration := time.Since(start)
	// 统计成功次数
	metrics.Inc(ctx, ot.MetricSuccName(), key, bucket)
	// 记录耗时统计
	metrics.RecordHistogram(ctx, float64(duration.Milliseconds()),
		"bucket", bucket,
		"operation", ot.String())
	return nil
}
