package rx_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/botchris/go-observer"
	"github.com/botchris/go-observer/rx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperable_Filter(t *testing.T) {
	t.Run("GIVEN a stream of integers from 1 to n WHEN filtering evens/odds THEN even/odds are split", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		n := 50
		prop := observer.NewProperty(nil)

		streamEven := rx.MakeOperable(ctx, prop.Observe()).
			Filter(func(_ context.Context, v interface{}) bool {
				return v.(int)%2 == 0
			})

		streamOdd := rx.MakeOperable(ctx, prop.Observe()).
			Filter(func(_ context.Context, v interface{}) bool {
				return v.(int)%2 != 0
			})

		for i := 1; i <= n; i++ {
			prop.Update(i)
		}
		prop.End()

		odds := streamOdd.ToSlice()
		evens := streamEven.ToSlice()

		require.NotEmpty(t, odds)
		require.NotEmpty(t, evens)

		for _, number := range odds {
			assert.NotZero(t, number.(int)%2)
		}

		for _, number := range evens {
			assert.Zero(t, number.(int)%2)
		}

		var collected []interface{}
		collected = append(collected, odds...)
		collected = append(collected, evens...)

		require.Len(t, collected, n)
	})

	t.Run("GIVEN a filter for a stream of integers WHEN context ends THEN stream ends with io.EOF", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// Start even stream reader
		prop := observer.NewProperty(nil)
		stream := rx.MakeOperable(ctx, prop.Observe()).
			Filter(func(_ context.Context, v interface{}) bool {
				return v.(int)%2 == 0
			})

		<-stream.Changes()
		stream.Next()

		read := stream.Value()
		require.NotNil(t, read)

		err, ok := read.(error)
		require.True(t, ok)
		require.Error(t, err)
		require.EqualError(t, err, io.EOF.Error())
	})

	t.Run("GIVEN a filtered stream WHEN filtering again THEN both filters are applied", func(t *testing.T) {
		ctx := context.Background()
		prop := observer.NewProperty(nil)

		stream := rx.MakeOperable(ctx, prop.Observe()).
			Filter(func(_ context.Context, v interface{}) bool {
				return v.(int)%2 == 0
			}).
			Filter(func(_ context.Context, v interface{}) bool {
				return v.(int)%8 == 0
			})

		for i := 1; i <= 200; i++ {
			prop.Update(i)
		}
		prop.End()

		results := stream.ToSlice()
		require.NotEmpty(t, results)

		for _, i := range results {
			require.True(t, i.(int)%2 == 0 && i.(int)%8 == 0, "received values must be divisible by 2 and 8")
		}
	})
}
