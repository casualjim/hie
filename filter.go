package hie

// Find returns the first match in the iterator, or none if it can't be found
func Find[T any](iter Iter[T], predicate Predicate[T]) Option[T] {
	for iter.HasNext() {
		elem := iter.Next()
		if predicate(elem) {
			return Some(elem)
		}
	}
	return None[T]()
}
