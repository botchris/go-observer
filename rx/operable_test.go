package rx_test

import (
	"context"
	"testing"

	"github.com/botchris/observer"
	"github.com/botchris/observer/rx"
	"github.com/stretchr/testify/require"
)

func TestOperable_OnComplete(t *testing.T) {
	t.Run("GIVEN an operable and a done ctx WHEN converting to slice THEN onComplete callback is invoked once", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		prop := observer.NewProperty[int](-1)
		onCompleteCalls := 0
		operable := rx.MakeOperable[int](ctx, prop.Observe()).
			OnComplete(func() {
				onCompleteCalls++
			})

		operable.ToSlice()

		require.EqualValues(t, 1, onCompleteCalls)
	})

	t.Run("GIVEN an empty operable and a done ctx WHEN converting to map THEN onComplete callback is invoked once", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		prop := observer.NewProperty[int](-1)
		onCompleteCalls := 0
		operable := rx.MakeOperable[int](ctx, prop.Observe()).
			OnComplete(func() {
				onCompleteCalls++
			})

		operable.ToMap(func(ctx context.Context, i interface{}) interface{} {
			return i
		})

		require.EqualValues(t, 1, onCompleteCalls)
	})

	t.Run("GIVEN an operable WHEN emitter ends and context ends after THEN onComplete callback is invoked once", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		prop := observer.NewProperty[int](-1)
		onCompleteCalls := 0
		operable := rx.MakeOperable(ctx, prop.Observe()).
			OnComplete(func() {
				onCompleteCalls++
			})

		prop.Update(1, 2, 3)
		prop.End()

		operable.Start()
		cancel()

		operable.ToSlice()

		<-operable.Done()

		require.EqualValues(t, 1, onCompleteCalls)
	})

	t.Run("GIVEN an operable WHEN ctx ends suddenly while iterating THEN onComplete callback is invoked once", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		prop := observer.NewProperty[int](-1)
		onCompleteCalls := 0
		operable := rx.MakeOperable(ctx, prop.Observe()).
			OnComplete(func() {
				onCompleteCalls++
			})

		prop.Update(1, 2, 3, 4)

		cancel()

		for {
			<-operable.Changes()
			v := operable.Next()

			cancel()
			prop.End()

			if v == 0 {
				break
			}
		}

		require.EqualValues(t, 1, onCompleteCalls)
	})
}
