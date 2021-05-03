package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/botchris/observer"
	"github.com/botchris/observer/rx"
	"github.com/stretchr/testify/require"
)

func TestOperable_Max(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	prop := observer.NewProperty(nil)
	stream := rx.MakeOperable(ctx, prop.Observe()).
		Max(func(_ context.Context, i1 interface{}, i2 interface{}) int {
			return i1.(int) - i2.(int)
		})

	prop.Update(2, 5, 1, 6, 3, 4)
	prop.End()

	items := stream.ToSlice()

	require.Len(t, items, 1)
	require.EqualValues(t, 6, items[0].(int))
}
