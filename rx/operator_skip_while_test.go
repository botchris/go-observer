package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/botchris/observer"
	"github.com/botchris/observer/rx"
	"github.com/stretchr/testify/require"
)

func TestOperable_SkipWhile(t *testing.T) {
	t.Run("GIVEN an integer emitter ranging 1 to 100 WHEN skipping numbers from 1 to 50 THEN only integer greater than 50 are received", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		prop := observer.NewProperty(nil)
		stream := rx.MakeOperable(ctx, prop.Observe()).
			SkipWhile(func(_ context.Context, v interface{}) bool {
				return v.(int) <= 50
			})

		for i := 1; i <= 100; i++ {
			prop.Update(i)
		}
		prop.End()

		items := stream.ToSlice()
		require.Greater(t, len(items), 0)

		for _, item := range items {
			require.Greater(t, item.(int), 50)
		}
	})
}
