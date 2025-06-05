package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClock is a mock implementation of the Clock interface
type MockClock struct {
	mock.Mock
}

func (m *MockClock) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func TestNow(t *testing.T) {
	t.Run("default backend returns current time", func(t *testing.T) {
		// Setup
		originalBackend := defaultBackend
		defer func() { defaultBackend = originalBackend }()

		mockClock := new(MockClock)
		expectedTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		mockClock.On("Now").Return(expectedTime)
		defaultBackend = mockClock

		// Execute
		result := Now()

		// Verify
		assert.Equal(t, expectedTime, result)
		mockClock.AssertExpectations(t)
	})

	t.Run("nil default backend panics", func(t *testing.T) {
		// Setup
		originalBackend := defaultBackend
		defer func() {
			defaultBackend = originalBackend
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		defaultBackend = nil

		// Execute
		Now()
	})
}
