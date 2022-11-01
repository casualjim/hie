package iter

import (
	"github.com/casualjim/hie"
	"github.com/casualjim/hie/option"
)

// Find returns the first match in the iterator, or none if it can't be found
func Find[T any](iter hie.Iter[T], predicate Predicate[T]) option.Option[T] {
	i := iter
	for i.HasNext() {
		elem := i.Next()
		if predicate(elem) {
			return option.Some(elem)
		}
	}
	return option.None[T]()
}

func First[T any](iter hie.Iter[T]) option.Option[T] {
	for iter.HasNext() {
		return option.Some(iter.Next())
	}
	return option.None[T]()
}
