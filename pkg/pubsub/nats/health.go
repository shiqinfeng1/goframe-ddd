package nats

import (
	"context"

	"github.com/shiqinfeng1/goframe-ddd/pkg/health"
)

const (
	natsBackend            = "Client"
	jetStreamStatusOK      = "OK"
	jetStreamStatusError   = "Error"
	jetStreamConnected     = "CONNECTED"
	jetStreamDisconnecting = "DISCONNECTED"
)

// Health checks the health of the NATS connection.
func (c *Client) Health(ctx context.Context) *health.Health {
	if c.connManager == nil {
		return &health.Health{
			Status: health.StatusDown,
		}
	}

	health := c.connManager.Health()
	health.Details["backend"] = natsBackend

	js, err := c.connManager.GetJetStream()
	if err != nil {
		health.Details["jetstream_enabled"] = false
		health.Details["jetstream_status"] = jetStreamStatusError + ": " + err.Error()
		return health
	}
	sm := NewStreamManager(js)
	// Call AccountInfo() to get jStream status
	jetStreamStatus, err := sm.GetJetStreamStatus(ctx)
	if err != nil {
		jetStreamStatus = jetStreamStatusError + ": " + err.Error()
	}

	health.Details["jetstream_enabled"] = true
	health.Details["jetstream_status"] = jetStreamStatus
	return health
}
