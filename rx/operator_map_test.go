package rx_test

import (
	"context"
	"math"
	"testing"

	"github.com/botchris/go-observer"
	"github.com/botchris/go-observer/rx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperable_Map(t *testing.T) {
	t.Run("GIVEN a stream of integers WHEN mapping to sqrt THEN sqrt values are received", func(t *testing.T) {
		type mappedValue struct {
			number float64
			sqrt   float64
		}

		ctx := context.Background()
		prop := observer.NewProperty(nil)

		stream := rx.MakeOperable(ctx, prop.Observe()).
			Map(func(_ context.Context, v interface{}) interface{} {
				return mappedValue{
					number: v.(float64),
					sqrt:   math.Sqrt(v.(float64)),
				}
			})

		// start writing values
		for i := 0; i <= 50; i++ {
			prop.Update(float64(i))
		}
		prop.End()

		results := stream.ToSlice()
		require.NotEmpty(t, results)
		require.Len(t, results, 51)

		for _, r := range results {
			require.IsType(t, mappedValue{}, r)
			assert.Equal(t, math.Sqrt(r.(mappedValue).number), r.(mappedValue).sqrt)
		}
	})
}
