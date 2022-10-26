# Hie

hasten ye collection operations

Hie contains a generic iterable with several functional combiner methods.

This library is built around the concept of an Iterable which is a very small interface that allows for 
lazy iteration over a sequence of values.

```go
type Iter[T any] interface {
 HasNext() bool
 Next() T
}
```

Due to the nature of Golangs generics we do need an entry point to allow for composability:

```go
type AsIter[T any] interface {
 AsIter() Iter[T]
}
```

You can find iter implementations for a slice and an option value.

## Combiner

* Map
* Filter
* FlatMap

## Terminators

* ForEach
* Fold
* Collect

## Option 

This library also contains an Option type that can be used with the same combinator.

## What's next

If I ever find time or the will to add

* an Either type
* a Result type (left biased either)
* function composition
* ...
