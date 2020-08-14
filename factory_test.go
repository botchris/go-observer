package observer_test

import (
	"context"
	"io"
	"sync"
	"testing"

	"github.com/christophercastro/observer"
	"github.com/stretchr/testify/require"
)

func TestFactory_Concat(t *testing.T) {
	ctx := context.Background()
	p1 := observer.NewProperty(nil)
	p2 := observer.NewProperty(nil)

	stream := observer.Concat(ctx, []observer.Stream{p1.Observe(), p2.Observe()})

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 1; i <= 10; i++ {
			p1.Update(i)
		}

		p1.End()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 11; i <= 20; i++ {
			p2.Update(i)
		}

		p2.End()
	}()

	// wait for emitters to complete
	wg.Wait()

	rcv := make([]int, 0)
	for {
		<-stream.Changes()
		v := stream.Next()

		if v == io.EOF {
			break
		}

		rcv = append(rcv, v.(int))
	}

	require.Len(t, rcv, 20)
	prev := 0

	for _, v := range rcv {
		require.EqualValues(t, prev+1, v)
		prev = v
	}
}
