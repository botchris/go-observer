package rx_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/botchris/observer"
	"github.com/botchris/observer/rx"
	"github.com/stretchr/testify/require"
)

func TestOperable_GroupBy(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	prop := observer.NewProperty[int](-1)
	count := 3
	stream := rx.MakeOperable[int](ctx, prop.Observe()).
		GroupBy(count, func(item int) int {
			return item % count
		})

	max := 10
	for i := 0; i <= max; i++ {
		prop.Update(i)
	}
	prop.End()

	groups := stream.ToSlice()
	require.Len(t, groups, count)

	output := make([]int, 0)
	for _, group := range groups {
		require.IsType(t, &rx.Operable[int]{}, group)

		results := group.(*rx.Operable[int]).ToSlice()
		fmt.Printf("new operable: %+v\n", results)
		require.NotEmpty(t, results)

		for _, n := range results {
			output = append(output, n.(int))
		}
	}

	require.Len(t, output, max+1)
}
