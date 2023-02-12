# observer

NOTE: This is a fork from [imkira/go-observer](https://github.com/imkira/go-observer)
which adds Rx-like operators.

observer is a [Go](http://golang.org) package that aims to simplify the problem
of channel-based broadcasting of events from one or more publishers to one or
more observers.

# Problem

The typical quick-and-dirty approach to notifying a set of observers in go is
to use channels and call each in a for loop, like the following:

```go
for _, channel := range channels {
  channel <- value
}
```

There are two problems with this approach:

- The broadcaster blocks every time some channel is not ready to be written to.
- If the broadcaster blocks for some channel, the remaining channels will not
  be written to (and therefore not receive the event) until the blocking
  channel is finally ready.
- It is O(N). The more observers you have, the worse this loop will behave.

Of course, this could be solved by creating one goroutine for each channel so
the broadcaster doesn't block. Unfortunately, this is heavy and
resource-consuming. This is especially bad if you have events being raised
frequently and a considerable number of observers.

# Approach

The way observer package tackles this problem is very simple. For every event,
a state object containing information about the event, and a channel is
created. State objects are managed using a singly linked list structure: every
state points to the next. When a new event is raised, a new state object is
appended to the list and the channel of the previous state is closed (this
helps notify all observers that the previous state is outdated).

Package observer defines 2 concepts:

- Property: An object that is continuously updated by one or more publishers.
- Stream: The list of values a property is updated to. For every property
update, that value is appended to the list in the order they happen, and is
only discarded when you advance to the next value.

# Memory Usage

The amount of memory used for one property is not dependent on the number of
observers. It should be proportional to the number of value updates since the
value last obtained by the slowest observer. As long as you keep advancing all
your observers, garbage collection will take place and keep memory usage
stable.

# How to Use

First, you need to install the package:

```
go get -u github.com/christophercastro/go-observer
```

Then, you need to include it in your source:

```go
import "github.com/christophercastro/go-observer"
```

The package will be imported with ```observer``` as name.

The following example creates one property that is updated every second by one
or more publishers, and observed by one or more observers.

## Documentation

For advanced usage, make sure to check the
[available documentation here](http://godoc.org/github.com/christophercastro/go-observer).

## Example: Creating a Property

The following code creates a property with initial value ```1```.

```go
val := 1
prop := observer.NewProperty(val)
```

After creating the property, you can pass it around to publishers or
observers as you want.

## Example: Publisher

The following code represents a publisher that increments the value of the
property by one every second.

```go
val := 1
for {
  time.Sleep(time.Second)
  val += 1
  fmt.Printf("will publish value: %d\n", val)
  prop.Update(val)
}
```

Note:

- Property is goroutine safe: you can use it concurrently from multiple
goroutines.

## Example: Observer

The following code represents an observer that prints the initial value of a
property and waits indefinitely for changes to its value. When there is a
change, the stream is advanced and the current value of the property is
printed.

```go
stream := prop.Observe()
val := stream.Value().(int)
fmt.Printf("initial value: %d\n", val)
for {
  select {
    // wait for changes
    case <-stream.Changes():
      // advance to next value
      stream.Next()
      // new value
      val = stream.Value().(int)
      fmt.Printf("got new value: %d\n", val)
  }
}
```

Note:

- Stream is not goroutine safe: You must create one stream by calling
  ```Property.Observe()``` or ```Stream.Clone()``` if you want to have
  concurrent observers for the same property or stream.

## Example

Please check
[examples/multiple.go](https://github.com/christophercastro/go-observer/blob/master/examples/multiple.go)
for a simple example on how to use multiple observers with a single updater.

# Operators

- `All`: determine whether all items emitted meet some criteria.
- `Buffer`: periodically gather items emitted into bundles and emit these bundles rather than emitting the items one at a time.
- `Contains`: determine whether a particular item was emitted or not.
- `Debounce`: only emit an item if a particular timespan has passed without it emitting another item.
- `Filter`: emit only those items that pass a predicate test.
- `Map`: transform the items by applying a function to each item.
- `Concat`: emit the emissions from two or more source streams without interleaving them.
- `SkipWhile`: discard items until a specified condition becomes false.
- `IgnoreElements`: do not emit any items but mirror its termination notification.
- `Last`: emit only the last item.
- `LastOrDefault`: emit only the last item. If fails to emit any items, it emits a default value.
- `Max`: determines and emits the maximum-valued item according to a comparator.
- `Min`: determines and emits the minimum-valued item according to a comparator.
- `Timestamp`: attaches a timestamp to each item indicating when it was emitted.
- `ToSlice`: collects every emitted item until EOF is reached an returns an slice holding each collected item.
- `ToMap`: convert the sequence of emitted items into a map keyed by a specified key function.
- `GroupBy`: divides an Operable into a set of Operable that each emit a different group of items from the original Operable, organized by key.
- `Distinct`: suppresses duplicate items.
- `DistinctUntilChanged`: suppresses consecutive duplicate items.
- `Count`: counts the number of items emitted and emit only this value.

## Example

```go
stream := observe.MakeOperable(ctx, prop.Observe()).
    Filter(func (item interface) bool {
        return item.(int)%2 == 0
    }).
    Filter(func (item interface) bool {
        return item.(int)%8 == 0
    })

for i := 1; 1 <= 100; i++ {
   prop.Update(i)
}

prop.End()

for {
  select {
    case <-stream.Changes():
      val = stream.Next().(int)
      if val == io.EOF {
        break
      }

      fmt.Printf("even value divisible by 8: %d\n", val)
  }
}
```