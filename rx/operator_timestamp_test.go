package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/botchris/observer"
	"github.com/botchris/observer/rx"
	"github.com/stretchr/testify/require"
)

func TestOperable_Timestamp(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	prop := observer.NewProperty(nil)
	stream := rx.MakeOperable(ctx, prop.Observe()).Timestamp()

	prop.Update(2, 5, 1, 6, 3, 4)
	prop.End()

	items := stream.ToSlice()
	require.Greater(t, len(items), 0)

	for _, item := range items {
		require.IsType(t, rx.TimestampItem{}, item)

		ts := item.(rx.TimestampItem)
		require.NotZero(t, ts.Timestamp)
		require.NotEmpty(t, ts.Item)
	}
}
