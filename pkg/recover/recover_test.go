package recover

import (
	"context"
	"errors"
	"testing"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/stretchr/testify/assert"
)

func TestRecovery(t *testing.T) {
	t.Run("no panic occurs", func(t *testing.T) {
		called := false
		recoverFunc := func(ctx context.Context, exception error) {
			called = true
		}

		Recovery(context.Background(), recoverFunc)
		assert.False(t, called)
	})

	t.Run("panic with error that has stack", func(t *testing.T) {
		called := false
		expectedErr := gerror.New("test error")
		recoverFunc := func(ctx context.Context, exception error) {
			called = true
			assert.Equal(t, expectedErr, exception)
		}

		func() {
			defer Recovery(context.Background(), recoverFunc)
			panic(expectedErr)
		}()
		assert.True(t, called)
	})

	t.Run("panic with error without stack", func(t *testing.T) {
		called := false
		expectedErr := errors.New("test error")
		recoverFunc := func(ctx context.Context, exception error) {
			called = true
			assert.Equal(t, expectedErr.Error(), exception.Error())
			assert.True(t, gerror.HasStack(exception))
		}

		func() {
			defer Recovery(context.Background(), recoverFunc)
			panic(expectedErr)
		}()
		assert.True(t, called)
	})

	t.Run("panic with non-error value", func(t *testing.T) {
		called := false
		expectedValue := "test panic"
		recoverFunc := func(ctx context.Context, exception error) {
			called = true
			assert.Contains(t, exception.Error(), expectedValue)
			assert.Equal(t, gcode.CodeInternalPanic, gerror.Code(exception))
		}

		func() {
			defer Recovery(context.Background(), recoverFunc)
			panic(expectedValue)
		}()
		assert.True(t, called)
	})

	t.Run("nil recoverFunc", func(t *testing.T) {
		assert.NotPanics(t, func() {
			func() {
				defer Recovery(context.Background(), nil)
				panic("test panic")
			}()
		})
	})
}
