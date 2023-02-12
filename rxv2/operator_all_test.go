package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/botchris/observer"
	"github.com/botchris/observer/rxv2"
	"github.com/stretchr/testify/require"
)

func TestOperable_All(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("GIVEN an emitter of constant values WHEN all operator is applied for such value THEN bool true is emitted at the end", func(t *testing.T) {
		prop := observer.NewProperty[string]("")
		stream := rx.All[string](ctx, prop.Observe(), func(v string) bool {
			return v == "yolo"
		})

		for i := 1; i <= 10; i++ {
			prop.Update("yolo")
		}

		prop.End()

		require.True(t, stream.AllMeet())
	})

	t.Run("GIVEN an emitter of constant values WHEN all operator is applied for another value THEN bool false is emitted at the end", func(t *testing.T) {
		prop := observer.NewProperty[string]("")
		stream := rx.All[string](ctx, prop.Observe(), func(v string) bool {
			return v == "yolo"
		})

		for i := 1; i <= 10; i++ {
			prop.Update("no yolo")
		}

		prop.End()

		require.False(t, stream.AllMeet())
	})
}
