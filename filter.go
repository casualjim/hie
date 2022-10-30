package hie

// Find returns the first match in the iterator, or none if it can't be found
func Find[T any](iter AsIter[T], predicate Predicate[T]) Option[T] {
	i := iter.AsIter()
	for i.HasNext() {
		elem := i.Next()
		if predicate(elem) {
			return Some(elem)
		}
	}
	return None[T]()
}
