package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/botchris/go-observer"
	"github.com/botchris/go-observer/rx"
	"github.com/stretchr/testify/require"
)

func TestOperable_DistinctUntilChanged(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	prop := observer.NewProperty(nil)
	stream := rx.MakeOperable(ctx, prop.Observe()).
		DistinctUntilChanged(func(_ context.Context, i interface{}) interface{} {
			return i
		})

	prop.Update(1, 2, 2, 1, 1, 3)
	prop.End()

	items := stream.ToSlice()
	require.Len(t, items, 4)

	require.EqualValues(t, 1, items[0])
	require.EqualValues(t, 2, items[1])
	require.EqualValues(t, 1, items[2])
	require.EqualValues(t, 3, items[3])
}
