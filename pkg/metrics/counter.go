package metrics

import (
	"context"

	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmetric"
	"github.com/shiqinfeng1/goframe-ddd/pkg/cache"
	"go.opentelemetry.io/otel/exporters/prometheus"
)

var (
	NatsKVStats = "nats_kvstore_stats"

	NatsPublishTotalCount       = "nats_publish_total_count"
	NatsPublishSuccessCount     = "nats_publish_success_count"
	NatsJsPublishTotalCount     = "nats_jspublish_total_count"
	NatsJsPublishSuccessCount   = "nats_jspublish_success_count"
	NatsSubscribeTotalCount     = "nats_subscribe_total_count"
	NatsSubscribeSuccessCount   = "nats_subscribe_success_count"
	NatsJsSubscribeTotalCount   = "nats_jssubscribe_total_count"
	NatsJsSubscribeSuccessCount = "nats_jssubscribe_success_count"
	meter                       = gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{
		Instrument:        "mgrid_pubsub",
		InstrumentVersion: "v1.0",
	})
	counter = map[string]gmetric.Counter{
		NatsPublishTotalCount: meter.MustCounter(
			NatsPublishTotalCount,
			gmetric.MetricOption{
				Help: "Publish total counts to nats",
				Unit: "count",
			},
		),
		NatsPublishSuccessCount: meter.MustCounter(
			NatsPublishSuccessCount,
			gmetric.MetricOption{
				Help: "Publish success counts to nats",
				Unit: "count",
			},
		),
		NatsJsPublishTotalCount: meter.MustCounter(
			NatsJsPublishTotalCount,
			gmetric.MetricOption{
				Help: "Publish total counts to nats",
				Unit: "count",
			},
		),
		NatsJsPublishSuccessCount: meter.MustCounter(
			NatsJsPublishSuccessCount,
			gmetric.MetricOption{
				Help: "Publish success counts to nats",
				Unit: "count",
			},
		),
		NatsSubscribeTotalCount: meter.MustCounter(
			NatsSubscribeTotalCount,
			gmetric.MetricOption{
				Help: "Subscribe total counts to nats",
				Unit: "count",
			},
		),
		NatsSubscribeSuccessCount: meter.MustCounter(
			NatsSubscribeSuccessCount,
			gmetric.MetricOption{
				Help: "Subscribe success counts to nats",
				Unit: "count",
			},
		),
		NatsJsSubscribeTotalCount: meter.MustCounter(
			NatsJsSubscribeTotalCount,
			gmetric.MetricOption{
				Help: "Subscribe total counts to nats",
				Unit: "count",
			},
		),
		NatsJsSubscribeSuccessCount: meter.MustCounter(
			NatsJsSubscribeSuccessCount,
			gmetric.MetricOption{
				Help: "Subscribe success counts to nats",
				Unit: "count",
			},
		),
	}
	natsKVStats = meter.MustHistogram(
		NatsKVStats,
		gmetric.MetricOption{
			Help:    "Response time of NATS KV operations in milliseconds.",
			Unit:    "ms",
			Buckets: []float64{.05, .075, .1, .125, .15, .2, .3, .5, .75, 1, 2, 3, 4, 5, 7.5, 10},
		},
	)
	provider   gmetric.Provider
	labelCache = cache.KV()
)

func init() {
	// Prometheus exporter to export metrics as Prometheus format.
	ctx := gctx.New()
	exporter, err := prometheus.New(
		prometheus.WithoutCounterSuffixes(),
		prometheus.WithoutUnits(),
	)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	// OpenTelemetry provider.
	provider = otelmetric.MustProvider(otelmetric.WithReader(exporter))
	provider.SetAsGlobal()
}

func Shutdown(ctx context.Context) {
	provider.Shutdown(ctx)
}

func IncCnt(ctx context.Context, name, label, value string) {
	v, err := labelCache.GetOrSet(ctx, value+label, gmetric.Option{
		Attributes: []gmetric.Attribute{
			gmetric.NewAttribute(label, value),
		}}, 0)
	if err != nil {
		return
	}
	var opt gmetric.Option
	if err := v.Scan(&opt); err != nil {
		return
	}

	if c, ok := counter[name]; ok {
		c.Inc(ctx, opt)
	}
}

func RecordHistogram(ctx context.Context, t float64, labels ...string) {
	attrs := []gmetric.Attribute{}
	if labels != nil {
		for i := 0; i < len(labels)-1; i += 2 {
			attrs = append(attrs, gmetric.NewAttribute(labels[i], labels[i+1]))
		}
	}
	natsKVStats.Record(t, gmetric.Option{
		Attributes: attrs,
	})
}
