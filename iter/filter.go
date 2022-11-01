package iter

import (
	"github.com/casualjim/hie"
	"github.com/casualjim/hie/opt"
)

// Find returns the first match in the iterator, or none if it can't be found
func Find[T any](iter hie.Iter[T], predicate Predicate[T]) opt.Option[T] {
	i := iter
	for i.HasNext() {
		elem := i.Next()
		if predicate(elem) {
			return opt.Some(elem)
		}
	}
	return opt.None[T]()
}

func First[T any](iter hie.Iter[T]) opt.Option[T] {
	for iter.HasNext() {
		return opt.Some(iter.Next())
	}
	return opt.None[T]()
}
