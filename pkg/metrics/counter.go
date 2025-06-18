package metrics

import (
	"context"

	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gmetric"
	"github.com/shiqinfeng1/goframe-ddd/pkg/cache"
	"go.opentelemetry.io/otel/exporters/prometheus"
)

var (
	provider    gmetric.Provider // 全局的监控指标管理工厂
	meterPubsub gmetric.Meter    // 管理全局的 Instrument，不同的 Meter 可以看做是不同的程序组件,例如订阅发布的meter，http的meter
	counter     map[string]gmetric.Counter
	natsKVStats gmetric.Histogram

	labelCache  = cache.KV()
	NatsKVStats = "nats.kvstore.stats"

	NatsPublishTotalCount     = "nats.publish.total.count"
	NatsPublishSuccessCount   = "nats.publish.success.count"
	NatsJSPublishTotalCount   = "nats.jspublish.total.count"
	NatsJSPublishSuccessCount = "nats.jspublish.success.count"
	NatsSubscribeTotalCount   = "nats.subscribe.total.count"
	NatsSubscribeSuccessCount = "nats.subscribe.success.count"
	NatsJSConsumeTotalCount   = "nats.jsconsume.total.count"
	NatsJSConsumeSuccessCount = "nats.jsconsume.success.count"

	NatsKVSetTotalCount     = "nats.kvset.total.count"
	NatsKVSetSuccessCount   = "nats.kvset.success.count"
	NatsKVGetTotalCount     = "nats.kvget.total.count"
	NatsKVGetSuccessCount   = "nats.kvget.success.count"
	NatsObjSetTotalCount    = "nats.objset.total.count"
	NatsObjSetSuccessCount  = "nats.objset.success.count"
	NatsObjGetTotalCount    = "nats.objget.total.count"
	NatsObjGetSuccessCount  = "nats.objget.success.count"
	NatsFileSetTotalCount   = "nats.fileset.total.count"
	NatsFileSetSuccessCount = "nats.fileset.success.count"
	NatsFileGetTotalCount   = "nats.fileget.total.count"
	NatsFileGetSuccessCount = "nats.fileget.success.count"
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
	if genv.Get("ENV").String() == "prod" {
		provider = otelmetric.MustProvider(otelmetric.WithReader(exporter))
	} else {
		provider = otelmetric.MustProvider(
			otelmetric.WithReader(exporter),
			otelmetric.WithBuiltInMetrics())
	}
	provider.SetAsGlobal()

	// 实例化一个pubsub的meter,一个Instrument
	meterPubsub = gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{
		Instrument:        "mgrid_pubsub",
		InstrumentVersion: "v1.0",
	})

	//pubsub的meter的各个不同类型的指标
	counter = map[string]gmetric.Counter{
		NatsPublishTotalCount: meterPubsub.MustCounter(
			NatsPublishTotalCount,
			gmetric.MetricOption{
				Help: "Publish total counts to nats",
				Unit: "count",
			},
		),
		NatsPublishSuccessCount: meterPubsub.MustCounter(
			NatsPublishSuccessCount,
			gmetric.MetricOption{
				Help: "Publish success counts to nats",
				Unit: "count",
			},
		),
		NatsJSPublishTotalCount: meterPubsub.MustCounter(
			NatsJSPublishTotalCount,
			gmetric.MetricOption{
				Help: "JS Publish total counts to nats",
				Unit: "count",
			},
		),
		NatsJSPublishSuccessCount: meterPubsub.MustCounter(
			NatsJSPublishSuccessCount,
			gmetric.MetricOption{
				Help: "JS Publish success counts to nats",
				Unit: "count",
			},
		),
		NatsSubscribeTotalCount: meterPubsub.MustCounter(
			NatsSubscribeTotalCount,
			gmetric.MetricOption{
				Help: "Subscribe total counts from nats",
				Unit: "count",
			},
		),
		NatsSubscribeSuccessCount: meterPubsub.MustCounter(
			NatsSubscribeSuccessCount,
			gmetric.MetricOption{
				Help: "Subscribe success counts from nats",
				Unit: "count",
			},
		),
		NatsJSConsumeTotalCount: meterPubsub.MustCounter(
			NatsJSConsumeTotalCount,
			gmetric.MetricOption{
				Help: "JS Subscribe total counts from nats",
				Unit: "count",
			},
		),
		NatsJSConsumeSuccessCount: meterPubsub.MustCounter(
			NatsJSConsumeSuccessCount,
			gmetric.MetricOption{
				Help: "JS Consume success counts from nats",
				Unit: "count",
			},
		),
		NatsKVSetTotalCount: meterPubsub.MustCounter(
			NatsKVSetTotalCount,
			gmetric.MetricOption{
				Help: "Set kv total counts to nats",
				Unit: "count",
			},
		),
		NatsKVSetSuccessCount: meterPubsub.MustCounter(
			NatsKVSetSuccessCount,
			gmetric.MetricOption{
				Help: "Set kv success counts to nats",
				Unit: "count",
			},
		),
		NatsKVGetTotalCount: meterPubsub.MustCounter(
			NatsKVGetTotalCount,
			gmetric.MetricOption{
				Help: "Get kv total counts from nats",
				Unit: "count",
			},
		),
		NatsKVGetSuccessCount: meterPubsub.MustCounter(
			NatsKVGetSuccessCount,
			gmetric.MetricOption{
				Help: "Set kv success counts to nats",
				Unit: "count",
			},
		),
		NatsObjSetTotalCount: meterPubsub.MustCounter(
			NatsObjSetTotalCount,
			gmetric.MetricOption{
				Help: "Set obj total counts to nats",
				Unit: "count",
			},
		),
		NatsObjSetSuccessCount: meterPubsub.MustCounter(
			NatsObjSetSuccessCount,
			gmetric.MetricOption{
				Help: "Set obj success counts to nats",
				Unit: "count",
			},
		),
		NatsObjGetTotalCount: meterPubsub.MustCounter(
			NatsObjGetTotalCount,
			gmetric.MetricOption{
				Help: "Get obj total counts from nats",
				Unit: "count",
			},
		),
		NatsObjGetSuccessCount: meterPubsub.MustCounter(
			NatsObjGetSuccessCount,
			gmetric.MetricOption{
				Help: "Get obj success counts from nats",
				Unit: "count",
			},
		),
		NatsFileSetTotalCount: meterPubsub.MustCounter(
			NatsFileSetTotalCount,
			gmetric.MetricOption{
				Help: "Set file total counts to nats",
				Unit: "count",
			},
		),
		NatsFileSetSuccessCount: meterPubsub.MustCounter(
			NatsFileSetSuccessCount,
			gmetric.MetricOption{
				Help: "Set file success counts to nats",
				Unit: "count",
			},
		),
		NatsFileGetTotalCount: meterPubsub.MustCounter(
			NatsFileGetTotalCount,
			gmetric.MetricOption{
				Help: "Get file total counts from nats",
				Unit: "count",
			},
		),
		NatsFileGetSuccessCount: meterPubsub.MustCounter(
			NatsFileGetSuccessCount,
			gmetric.MetricOption{
				Help: "Get file success counts from nats",
				Unit: "count",
			},
		),
	}
	natsKVStats = meterPubsub.MustHistogram(
		NatsKVStats,
		gmetric.MetricOption{
			Help:    "Response time of NATS KV operations in milliseconds.",
			Unit:    "ms",
			Buckets: []float64{.05, .075, .1, .125, .15, .2, .3, .5, .75, 1, 2, 3, 4, 5, 7.5, 10},
		},
	)
}

func Shutdown(ctx context.Context) {
	provider.Shutdown(ctx)
}

func Inc(ctx context.Context, name, label, value string) {
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
