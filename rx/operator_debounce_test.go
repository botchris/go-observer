package rx_test

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/botchris/go-observer"
	"github.com/botchris/go-observer/rx"
	"github.com/stretchr/testify/require"
)

func TestOperable_Debounce(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	prop := observer.NewProperty(nil)
	stream := rx.MakeOperable(ctx, prop.Observe()).Debounce(100 * time.Millisecond)

	var wg sync.WaitGroup

	rcv := make([]int, 0)
	consumerReady := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		close(consumerReady)

		for {
			<-stream.Changes()
			v := stream.Next()
			if v == io.EOF {
				break
			}

			rcv = append(rcv, v.(int))
		}
	}()

	<-consumerReady

	sent := make([]int, 0)
	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 1

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(50 * time.Millisecond):
				sent = append(sent, i)
				prop.Update(i)
				i++
			}
		}
	}()

	wg.Wait()

	require.Greater(t, len(sent), 0)
	require.Greater(t, len(rcv), 0)
	require.Greater(t, len(sent), len(rcv))
}
