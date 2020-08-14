package observer_test

import (
	"context"
	"io"
	"math"
	"sync"
	"testing"

	"github.com/christophercastro/observer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperable_Map(t *testing.T) {
	t.Run("GIVEN a stream of integers WHEN mapping to sqrt THEN sqrt values are received", func(t *testing.T) {
		var wg sync.WaitGroup
		type mappedValue struct {
			number float64
			sqrt   float64
		}

		ctx := context.Background()
		prop := observer.NewProperty(nil)
		result := make([]*mappedValue, 0)
		wg.Add(1)

		// Start even stream reader
		streamSqrt := observer.MakeOperable(ctx, prop.Observe()).
			Map(func(v interface{}) interface{} {
				return &mappedValue{
					number: v.(float64),
					sqrt:   math.Sqrt(v.(float64)),
				}
			})

		go func() {
			defer wg.Done()

			for {
				<-streamSqrt.Changes()
				value := streamSqrt.Next()
				if value == io.EOF {
					return
				}

				read := value.(*mappedValue)
				result = append(result, read)

				//println("number:", fmt.Sprintf("%.2f", read.number), "sqrt:", fmt.Sprintf("%.2f", read.sqrt))
			}
		}()

		// start writing values
		for i := 0; i <= 50; i++ {
			prop.Update(float64(i))
		}
		prop.End()

		// wait for streams to finish
		wg.Wait()

		require.NotEmpty(t, result)
		require.Len(t, result, 51)

		for _, r := range result {
			assert.Equal(t, math.Sqrt(r.number), r.sqrt)
		}
	})
}
