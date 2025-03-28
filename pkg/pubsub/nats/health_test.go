package nats

import (
	"testing"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/shiqinfeng1/goframe-ddd/pkg/health"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	NATSServer = "nats://localhost:4222"
)

func TestNATSClient_Health(t *testing.T) {
	testCases := defineHealthTestCases()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runHealthTestCase(t, tc)
		})
	}
}

func defineHealthTestCases() []healthTestCase {
	return []healthTestCase{
		{
			name: "HealthyConnection",
			setupMocks: func(mockConnManager *MockConnectionManagerIntf, mockJS *MockJetStream) {
				mockConnManager.EXPECT().Health().Return(&health.Health{
					Status: health.StatusUp,
					Details: map[string]any{
						"host":              NATSServer,
						"connection_status": jetStreamConnected,
					},
				})
				mockConnManager.EXPECT().getJetStream().Return(mockJS, nil)
				mockJS.EXPECT().AccountInfo(gomock.Any()).Return(&jetstream.AccountInfo{}, nil)
			},
			expectedStatus: health.StatusUp,
			expectedDetails: map[string]any{
				"host":              NATSServer,
				"backend":           natsBackend,
				"connection_status": jetStreamConnected,
				"jetstream_enabled": true,
				"jetstream_status":  jetStreamStatusOK,
			},
		},
		{
			name: "DisconnectedStatus",
			setupMocks: func(mockConnManager *MockConnectionManagerIntf, _ *MockJetStream) {
				mockConnManager.EXPECT().Health().Return(&health.Health{
					Status: health.StatusDown,
					Details: map[string]any{
						"host":              NATSServer,
						"connection_status": jetStreamDisconnecting,
					},
				})
				mockConnManager.EXPECT().getJetStream().Return(nil, errJetStreamNotConfigured)
			},
			expectedStatus: health.StatusDown,
			expectedDetails: map[string]any{
				"host":              NATSServer,
				"backend":           natsBackend,
				"connection_status": jetStreamDisconnecting,
				"jetstream_enabled": false,
				"jetstream_status":  jetStreamStatusError + ": jStream is not configured",
			},
		},
		{
			name: "JetStreamError",
			setupMocks: func(mockConnManager *MockConnectionManagerIntf, mockJS *MockJetStream) {
				mockConnManager.EXPECT().Health().Return(&health.Health{
					Status: health.StatusUp,
					Details: map[string]any{
						"host":              NATSServer,
						"connection_status": jetStreamConnected,
					},
				})
				mockConnManager.EXPECT().getJetStream().Return(mockJS, nil)
				mockJS.EXPECT().AccountInfo(gomock.Any()).Return(nil, errJetStream)
			},
			expectedStatus: health.StatusUp,
			expectedDetails: map[string]any{
				"host":              NATSServer,
				"backend":           natsBackend,
				"connection_status": jetStreamConnected,
				"jetstream_enabled": true,
				"jetstream_status":  jetStreamStatusError + ": " + errJetStream.Error(),
			},
		},
	}
}

func runHealthTestCase(t *testing.T, tc healthTestCase) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnManager := NewMockConnectionManagerIntf(ctrl)
	mockJS := NewMockJetStream(ctrl)

	tc.setupMocks(mockConnManager, mockJS)

	client := &Client{
		connManager: mockConnManager,
		Config:      &Config{Server: NATSServer},
	}

	h := client.Health(t.Context())

	assert.Equal(t, tc.expectedStatus, h.Status)
	assert.Equal(t, tc.expectedDetails, h.Details)
}

type healthTestCase struct {
	name            string
	setupMocks      func(*MockConnectionManagerIntf, *MockJetStream)
	expectedStatus  string
	expectedDetails map[string]any
}
