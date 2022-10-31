# Hie

hasten ye collection operations

Hie contains a generic iterable with several functional combiner methods.

This library is built around the concept of an Iterable which is a very small interface that allows for 
lazy iteration over a sequence of values.  With this you can combine many filter, map and flatmap invocations and they will only be evaluated when the first `Next` method is called.

You can interrupt both ForEach and Fold so you can stop iterating and throw away the rest of the results.

```go
type Iter[T any] interface {
 HasNext() bool
 Next() T
}
```

You can find iter implementations for a slice and an option value.

## Combiner

* Map
* Filter
* FilterMap
* FlatMap
* Union: return the unique elements
* Intersect: return the intersection of 2 iterators
* Concat: combine several iterators into 1

## Terminators

* ForEach
* Fold
* Collect
* Difference
* Symmetric Difference
* Find: return an option with the first matching value
* First: return the first element of a collection as an option

## Option 

This library also contains an Option type that can be used with the same combinators.

## What's next

If I ever find time or the will to add

* an Either type
* a Result type (left biased either)
* function composition
* ...
