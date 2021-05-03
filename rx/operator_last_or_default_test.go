package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/botchris/observer"
	"github.com/botchris/observer/rx"
	"github.com/stretchr/testify/require"
)

func TestOperable_LastOrDefault(t *testing.T) {
	t.Run("GIVEN an empty emitter WHEN applying last or default THEN default value is received", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		prop := observer.NewProperty(nil)
		stream := rx.MakeOperable(ctx, prop.Observe()).LastOrDefault("default")

		prop.End()

		results := stream.ToSlice()

		require.Len(t, results, 1)
		require.EqualValues(t, "default", results[0].(string))
	})

	t.Run("GIVEN a non-empty emitter WHEN applying last or default THEN last value is received", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		prop := observer.NewProperty(nil)
		stream := rx.MakeOperable(ctx, prop.Observe()).LastOrDefault("default")

		prop.Update("no default")
		prop.End()

		results := stream.ToSlice()

		require.Len(t, results, 1)
		require.EqualValues(t, "no default", results[0].(string))
	})
}
