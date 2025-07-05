package xcontext

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDetach(t *testing.T) {

	t.Run("given success, when Detach(ctx), then parent-child Deadline is detached", func(t *testing.T) {
		ctx := context.Background()
		d := time.Now().Add(5 * time.Second)
		ctx, cancel := context.WithDeadline(ctx, d)
		defer cancel()
		ctxD, ctxOk := ctx.Deadline()

		detachedCtx := Detach(ctx)
		detachedD, detachedOk := detachedCtx.Deadline()

		assert.NotEqual(t, ctxD, detachedD)
		assert.NotEqual(t, ctxOk, detachedOk)
	})

	t.Run("given success, when Detach(ctx), then parent-child Done is detached", func(t *testing.T) {
		ctx := context.Background()
		d := time.Now().Add(1 * time.Second)
		ctx, cancel := context.WithDeadline(ctx, d)
		defer cancel()

		detachedCtx := Detach(ctx)

		assert.Nil(t, detachedCtx.Done())
		assert.NoError(t, detachedCtx.Err())

		assert.NotNil(t, <-ctx.Done())
		assert.Error(t, ctx.Err())
	})

	t.Run("given success, when Detach(ctx), then parent-child context holds equal value", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "key", "val")
		val := ctx.Value("key")

		detachedCtx := Detach(ctx)
		detachedVal := detachedCtx.Value("key")

		assert.Equal(t, val, detachedVal)
	})
}
