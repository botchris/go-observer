package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/botchris/observer"
	"github.com/botchris/observer/rx"
	"github.com/stretchr/testify/require"
)

func TestOperable_All(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	t.Run("GIVEN an emitter of constant values WHEN all operator is applied for such value THEN bool true is emitted at the end", func(t *testing.T) {
		prop := observer.NewProperty[string]("")
		stream := rx.MakeOperable[string](ctx, prop.Observe()).
			All(func(_ context.Context, v string) bool {
				return v == "yolo"
			})

		for i := 1; i <= 10; i++ {
			prop.Update("yolo")
		}
		prop.End()

		results := stream.ToSlice()
		require.Len(t, results, 1)
		require.True(t, results[0].(bool))
	})

	t.Run("GIVEN an emitter of constant values WHEN all operator is applied for another value THEN bool false is emitted at the end", func(t *testing.T) {
		prop := observer.NewProperty[string]("")
		stream := rx.MakeOperable[string](ctx, prop.Observe()).
			All(func(_ context.Context, v string) bool {
				return v == "yolo"
			})

		for i := 1; i <= 10; i++ {
			prop.Update("no yolo")
		}
		prop.End()

		results := stream.ToSlice()
		require.Len(t, results, 1)
		require.False(t, results[0].(bool))
	})
}
