package main_test

import (
	"context"
	"fmt"
	"io"
	"math"
	"testing"
	"time"

	"github.com/botchris/observer"
	"github.com/botchris/observer/rx"
	"github.com/stretchr/testify/require"
)

func TestOperable_Integration(t *testing.T) {
	t.Run("GIVEN an observable property with multiple operators WHEN context is cancelled THEN resulting stream ends with io.EOF", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		prop := observer.NewProperty(nil)
		operable := rx.MakeOperable(ctx, prop.Observe())

		for i := 1; i <= 100; i++ {
			prop.Update(i)
		}

		stream := operable.
			Filter(func(_ context.Context, v interface{}) bool {
				return v.(int)%2 == 0
			}).
			Filter(func(_ context.Context, v interface{}) bool {
				return v.(int)%8 == 0
			}).
			Map(func(_ context.Context, i interface{}) interface{} {
				sqrt := math.Sqrt(float64(i.(int)))

				return fmt.Sprintf("sqrt(%d) = %f", i, sqrt)
			})

		<-ctx.Done()

		read := make([]interface{}, 0)
		for {
			<-stream.Changes()
			v := stream.Next()
			read = append(read, v)

			if v == io.EOF {
				break
			}
		}

		require.NotEmpty(t, read)
		last := read[len(read)-1]

		err, ok := last.(error)
		require.True(t, ok, "last message is expected to be io.EOF")
		require.Error(t, err)
		require.EqualError(t, err, io.EOF.Error())
	})

	t.Run("GIVEN an observable with multiple operators that produce numbers lower than 100 WHEN applying an All operator THEN true is emitted", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		prop := observer.NewProperty(nil)
		operable := rx.MakeOperable(ctx, prop.Observe())

		stream := operable.
			Filter(func(_ context.Context, v interface{}) bool {
				return v.(int)%2 == 0
			}).
			Filter(func(_ context.Context, v interface{}) bool {
				return v.(int)%8 == 0
			}).
			Map(func(_ context.Context, i interface{}) interface{} {
				sqrt := math.Sqrt(float64(i.(int)))

				return sqrt
			}).
			All(func(_ context.Context, v interface{}) bool {
				return v.(float64) < 100
			})

		for i := 1; i <= 100; i++ {
			prop.Update(i)
		}
		prop.End()

		<-stream.Changes()
		value := stream.Next()
		require.True(t, value.(bool))

		<-stream.Changes()
		value = stream.Next()
		require.EqualValues(t, value, io.EOF)
	})

	t.Run("GIVEN two filtered observables WHEN concat them THEN concat emits correct values", func(t *testing.T) {
		/*
			Range(1..100) -> Property1 -(filter:even)-> Stream1
			Range(1..100) -> Property2 -(filter:odd)-> Stream2
			Concat(Stream1, Stream2) -(filter:n>50)-> Stream3

			Stream3 should see a sequence of even values between [50, 100], followed
			by a sequence of odd values between [50, 100].
		*/
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		p1 := observer.NewProperty(nil)
		o1 := rx.MakeOperable(ctx, p1.Observe()).
			Filter(func(_ context.Context, v interface{}) bool {
				// even
				return v.(int)%2 == 0
			})

		p2 := observer.NewProperty(nil)
		o2 := rx.MakeOperable(ctx, p1.Observe()).
			Filter(func(_ context.Context, v interface{}) bool {
				// odd
				return v.(int)%2 != 0
			})

		for i := 1; i <= 100; i++ {
			p1.Update(i)
			p2.Update(i)
		}

		p1.End()
		p2.End()

		concat := rx.Concat(ctx, []observer.Stream{o1, o2}).
			Filter(func(_ context.Context, v interface{}) bool {
				return v.(int) > 50
			})

		for {
			<-concat.Changes()
			v := concat.Next()
			if v == io.EOF {
				break
			}

			require.NotEmpty(t, v)
			require.Greater(t, v.(int), 50)
		}
	})
}

func TestOperable_ToMap(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	prop := observer.NewProperty(nil)
	stream := rx.MakeOperable(ctx, prop.Observe())

	prop.Update(1, 2, 3, io.EOF)

	results := stream.ToMap(func(_ context.Context, i interface{}) interface{} {
		return i.(int) * 10
	})

	require.Len(t, results, 3)
	for k, v := range results {
		require.EqualValues(t, v.(int)*10, k.(int))
	}
}
