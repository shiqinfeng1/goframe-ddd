package metrics

import (
	"context"

	"github.com/gogf/gf/contrib/metric/otelmetric/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmetric"
	"go.opentelemetry.io/otel/exporters/prometheus"
)

var (
	NatsPublishTotalCount     = "nats_publish_total_count"
	NatsPublishSuccessCount   = "nats_publish_success_count"
	NatsSubscribeTotalCount   = "nats_subscribe_total_count"
	NatsSubscribeSuccessCount = "nats_subscribe_success_count"
	meter                     = gmetric.GetGlobalProvider().Meter(gmetric.MeterOption{
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
	}
	provider gmetric.Provider
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

func IncrementCounter(ctx context.Context, name, label, value string) {
	if c, ok := counter[name]; ok {
		c.Inc(ctx, gmetric.Option{
			Attributes: []gmetric.Attribute{
				gmetric.NewAttribute(label, value),
			},
		})
	}
}
