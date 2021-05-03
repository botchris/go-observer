package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/botchris/observer"
	"github.com/botchris/observer/rx"
	"github.com/stretchr/testify/require"
)

func TestOperable_Last(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	prop := observer.NewProperty(nil)
	stream := rx.MakeOperable(ctx, prop.Observe()).Last()

	for i := 1; i <= 10; i++ {
		prop.Update(i)
	}
	prop.End()

	results := stream.ToSlice()

	require.Len(t, results, 1)
	require.EqualValues(t, 10, results[0].(int))
}
