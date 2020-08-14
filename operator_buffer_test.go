package observer_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/christophercastro/observer"
	"github.com/stretchr/testify/require"
)

func TestOperable_Buffer(t *testing.T) {
	t.Run("GIVEN an observable that emits 15 values WHEN buffering by 5 THEN 3 packages of 5 elements are emitted instead", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		prop := observer.NewProperty(nil)
		stream := observer.MakeOperable(ctx, prop.Observe()).BufferWithCount(5)

		for i := 1; i <= 15; i++ {
			prop.Update(i)
		}
		prop.End()

		// wait for stream to complete
		<-stream.Done()

		packsCount := 0
		for {
			<-stream.Changes()
			value := stream.Next()

			if value == io.EOF {
				break
			}

			require.Len(t, value, 5)
			packsCount++
		}

		require.EqualValues(t, 3, packsCount)
	})
}
