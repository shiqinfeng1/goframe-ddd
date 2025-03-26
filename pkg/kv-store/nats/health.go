package nats

import (
	context "context"
	"fmt"
	time "time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/pkg/health"
)

func (c *Client) Health(ctx context.Context) *health.Health {
	start := time.Now()

	h := &health.Health{
		Details: make(map[string]any),
	}

	h.Details["url"] = c.configs.Server
	h.Details["bucket"] = c.configs.Bucket

	_, err := c.js.AccountInfo(ctx)
	if err != nil {
		h.Status = "DOWN"

		g.Log().Debug(ctx, &Log{
			Type:     "HEALTH CHECK",
			Key:      "health",
			Value:    fmt.Sprintf("Connection failed for bucket '%s' at '%s'", c.configs.Bucket, c.configs.Server),
			Duration: time.Since(start).Microseconds(),
		})

		return h
	}

	h.Status = "UP"

	g.Log().Debug(ctx, &Log{
		Type:     "HEALTH CHECK",
		Key:      "health",
		Value:    fmt.Sprintf("Checking connection status for bucket '%s' at '%s'", c.configs.Bucket, c.configs.Server),
		Duration: time.Since(start).Microseconds(),
	})

	return h
}
