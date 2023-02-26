package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/botchris/go-observer"
	"github.com/botchris/go-observer/rx"
	"github.com/stretchr/testify/require"
)

func TestOperable_Contains(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	t.Run("GIVEN an emitter of number from 1 to 25 WHEN operator contains for number 22 is applied THEN bool true is emitted at the end", func(t *testing.T) {
		prop := observer.NewProperty(nil)
		stream := rx.MakeOperable(ctx, prop.Observe()).
			Contains(func(_ context.Context, v interface{}) bool {
				return v == 22
			})

		for i := 1; i <= 25; i++ {
			prop.Update(i)
		}
		prop.End()

		results := stream.ToSlice()

		require.Len(t, results, 1)
		require.True(t, results[0].(bool))
	})

	t.Run("GIVEN an emitter of number from 1 to 25 WHEN operator contains for number 27 is applied THEN bool false is emitted at the end", func(t *testing.T) {
		prop := observer.NewProperty(nil)
		stream := rx.MakeOperable(ctx, prop.Observe()).
			Contains(func(_ context.Context, v interface{}) bool {
				return v == 27
			})

		for i := 1; i <= 25; i++ {
			prop.Update(i)
		}
		prop.End()

		results := stream.ToSlice()

		require.Len(t, results, 1)
		require.False(t, results[0].(bool))
	})
}
