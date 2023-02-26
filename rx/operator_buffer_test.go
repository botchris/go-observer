package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/botchris/go-observer"
	"github.com/botchris/go-observer/rx"
	"github.com/stretchr/testify/require"
)

func TestOperable_Buffer(t *testing.T) {
	t.Run("GIVEN an observable that emits 15 values WHEN buffering by 5 THEN 3 packages of 5 elements are emitted instead", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		prop := observer.NewProperty(nil)
		stream := rx.MakeOperable(ctx, prop.Observe()).BufferWithCount(5)

		for i := 1; i <= 15; i++ {
			prop.Update(i)
		}
		prop.End()

		results := stream.ToSlice()
		require.Len(t, results, 3)

		for _, pack := range results {
			require.Len(t, pack, 5)
		}
	})
}
