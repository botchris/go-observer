package observer_test

import (
	"context"
	"testing"
	"time"

	"github.com/botchris/observer"
	"github.com/stretchr/testify/require"
)

func TestOperable_Contains(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	{
		prop := observer.NewProperty(nil)
		stream := observer.MakeOperable(ctx, prop.Observe()).
			Contains(func(v interface{}) bool {
				return v == 22
			})

		for i := 1; i <= 25; i++ {
			prop.Update(i)
		}
		prop.End()

		// wait for stream to complete
		<-stream.Done()

		<-stream.Changes()
		value := stream.Next()
		require.True(t, value.(bool))
	}

	{
		prop := observer.NewProperty(nil)
		stream := observer.MakeOperable(ctx, prop.Observe()).
			Contains(func(v interface{}) bool {
				return v == 27
			})

		for i := 1; i <= 25; i++ {
			prop.Update(i)
		}
		prop.End()

		// wait for stream to complete
		<-stream.Done()

		<-stream.Changes()
		value := stream.Next()
		require.False(t, value.(bool))
	}
}
