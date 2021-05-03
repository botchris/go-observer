package observer_test

import (
	"context"
	"io"
	"sync"
	"testing"

	"github.com/botchris/observer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperable_Filter(t *testing.T) {
	t.Run("GIVEN a stream of integers WHEN filtering evens/odds THEN even/odds are split", func(t *testing.T) {
		var wg sync.WaitGroup
		n := 50
		ctx := context.Background()
		prop := observer.NewProperty(nil)
		odds := make([]int, 0)
		evens := make([]int, 0)
		wg.Add(2)

		// Start even stream reader
		streamEven := observer.MakeOperable(ctx, prop.Observe()).
			Filter(func(v interface{}) bool {
				return v.(int)%2 == 0
			})

		go func() {
			defer wg.Done()

			for {
				<-streamEven.Changes()
				value := streamEven.Next()
				if value == io.EOF {
					return
				}

				read := streamEven.Value().(int)
				//println("even:", read)
				assert.True(t, read%2 == 0)
				evens = append(evens, read)
			}
		}()

		// Start odd stream reader
		streamOdd := observer.MakeOperable(ctx, prop.Observe()).
			Filter(func(v interface{}) bool {
				return v.(int)%2 != 0
			})

		go func() {
			defer wg.Done()

			for {
				<-streamOdd.Changes()
				value := streamOdd.Next()
				if value == io.EOF {
					return
				}

				read := streamOdd.Value().(int)
				//println("odd:", read)
				assert.True(t, read%2 != 0)
				odds = append(odds, read)
			}
		}()

		// start writing values
		for i := 0; i <= n; i++ {
			prop.Update(i)
		}
		prop.End()

		// wait for streams to finish
		wg.Wait()

		require.NotEmpty(t, odds)
		require.NotEmpty(t, evens)

		var collected []int
		collected = append(collected, odds...)
		collected = append(collected, evens...)

		require.Len(t, collected, n+1)
	})

	t.Run("GIVEN a filter for a stream of integers WHEN context ends THEN stream ends with io.EOF", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// Start even stream reader
		prop := observer.NewProperty(nil)
		stream := observer.MakeOperable(ctx, prop.Observe()).
			Filter(func(v interface{}) bool {
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

		stream := observer.MakeOperable(ctx, prop.Observe()).
			Filter(func(v interface{}) bool {
				return v.(int)%2 == 0
			}).
			Filter(func(v interface{}) bool {
				return v.(int)%8 == 0
			})

		var wg sync.WaitGroup

		wg.Add(1)
		seenValues := make([]int, 0)
		go func() {
			defer wg.Done()

			for {
				<-stream.Changes()
				value := stream.Next()
				if value == io.EOF {
					return
				}

				seenValues = append(seenValues, value.(int))
			}
		}()

		for i := 1; i <= 200; i++ {
			prop.Update(i)
		}
		prop.End()

		wg.Wait()

		require.NotEmpty(t, seenValues)
		for _, i := range seenValues {
			require.True(t, i%2 == 0 && i%8 == 0, "received values must be divisible by 2 and 8")
		}
	})
}
