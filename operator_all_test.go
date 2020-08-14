package observer_test

import (
	"context"
	"testing"
	"time"

	"github.com/christophercastro/observer"
	"github.com/stretchr/testify/require"
)

func TestOperable_All(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	{
		prop := observer.NewProperty(nil)
		stream := observer.MakeOperable(ctx, prop.Observe()).
			All(func(v interface{}) bool {
				return v == "yolo"
			})

		for i := 1; i <= 10; i++ {
			prop.Update("yolo")
		}
		prop.End()

		<-stream.Changes()
		value := stream.Next()

		require.True(t, value.(bool))
	}

	{
		prop := observer.NewProperty(nil)
		stream := observer.MakeOperable(ctx, prop.Observe()).
			All(func(v interface{}) bool {
				return v == "yolo"
			})

		for i := 1; i <= 10; i++ {
			prop.Update("no yolo")
		}
		prop.End()

		<-stream.Changes()
		value := stream.Next()

		require.False(t, value.(bool))
	}
}
