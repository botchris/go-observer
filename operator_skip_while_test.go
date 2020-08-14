package observer_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/christophercastro/observer"
	"github.com/stretchr/testify/require"
)

func TestOperable_SkipWhile(t *testing.T) {
	t.Run("GIVEN an integer emitter ranging 1 to 100 WHEN skipping numbers from 1 to 50 THEN only integer greater than 50 are received", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		prop := observer.NewProperty(nil)
		stream := observer.MakeOperable(ctx, prop.Observe()).
			SkipWhile(func(v interface{}) bool {
				return v.(int) <= 50
			})

		for i := 1; i <= 100; i++ {
			prop.Update(i)
		}
		prop.End()

		for {
			<-stream.Changes()
			v := stream.Next()
			if v == io.EOF {
				break
			}

			require.Greater(t, v.(int), 50)
		}
	})
}
