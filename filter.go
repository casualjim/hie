package hie

// Find returns the first match in the iterator, or none if it can't be found
func Find[T any](iter Iter[T], predicate Predicate[T]) Option[T] {
	i := iter
	for i.HasNext() {
		elem := i.Next()
		if predicate(elem) {
			return Some(elem)
		}
	}
	return None[T]()
}

func First[T any](iter Iter[T]) Option[T] {
	for iter.HasNext() {
		return Some(iter.Next())
	}
	return None[T]()
}
