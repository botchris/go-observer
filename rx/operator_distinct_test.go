package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/botchris/go-observer"
	"github.com/botchris/go-observer/rx"
	"github.com/stretchr/testify/require"
)

func TestOperable_Distinct(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	prop := observer.NewProperty(nil)
	stream := rx.MakeOperable(ctx, prop.Observe()).
		Distinct(func(_ context.Context, i interface{}) interface{} {
			return i
		})

	prop.Update(1, 2, 2, 3, 4, 4, 5)
	prop.End()

	items := stream.ToSlice()
	require.Len(t, items, 5)

	for i := 1; i <= 5; i++ {
		require.EqualValues(t, i, items[i-1])
	}
}
