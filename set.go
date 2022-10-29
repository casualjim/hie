package hie

// Contains returns true if an element is present in the collection
func Contains[T comparable](iter AsIter[T], elem T) bool {
	i := iter.AsIter()
	for i.HasNext() {
		if elem == i.Next() {
			return true
		}
	}
	return false
}

// Exists returns true if predicate function return true.
func Exists[T any](iter AsIter[T], predicate Predicate[T]) bool {
	i := iter.AsIter()
	for i.HasNext() {
		if predicate(i.Next()) {
			return true
		}
	}
	return false
}

// All returns true if the predicate returns true for all of the elements in the collection or if the collection is empty.
func All[T any](iter AsIter[T], predicate Predicate[T]) bool {
	i := iter.AsIter()
	for i.HasNext() {
		if !predicate(i.Next()) {
			return false
		}
	}
	return true
}

// NoneExist returns true if the predicate returns false for all the elements in the collection or if the collection is empty
func NoneExist[T any](iter AsIter[T], predicate Predicate[T]) bool {
	return All(iter, func(elem T) bool { return !predicate(elem) })
}

// IsSubset returns true if all elements of a subset are contained into a collection or if the subset is empty.
func IsSubset[T comparable](iter AsIter[T], subset AsIter[T]) bool {
	si := subset.AsIter()
	cached := Collect(iter)
	rewindable := Slice(cached...)
	for si.HasNext() {
		if !Contains(rewindable, si.Next()) {
			return false
		}
	}
	return true
}

// IsDisjoint returns true if none of the elements of the subset set are contained by the superset
func IsDisjoint[T comparable](iter AsIter[T], other AsIter[T]) bool {
	oi := other.AsIter()
	cached := Slice(Collect(iter)...)
	for oi.HasNext() {
		if Contains(cached, oi.Next()) {
			return false
		}
	}
	return true
}

// Intersect returns the intersection between two collections.
func Intersect[T comparable](first AsIter[T], second AsIter[T]) AsIter[T] {
	seen := make(map[T]struct{})
	firsti := first.AsIter()
	for firsti.HasNext() {
		seen[firsti.Next()] = struct{}{}
	}

	return &intersectingIter[T]{
		seen:        seen,
		second:      second.AsIter(),
		secondMatch: None[T](),
	}
}

type intersectingIter[T comparable] struct {
	seen        map[T]struct{}
	second      Iter[T]
	secondMatch Option[T]
}

func (i *intersectingIter[T]) HasNext() bool {
	for i.second.HasNext() {
		val := i.second.Next()
		if _, seen := i.seen[val]; seen {
			i.secondMatch = Some(val)
			return true
		}
	}
	return false
}

func (i *intersectingIter[T]) Next() T {
	val := i.secondMatch
	i.secondMatch = None[T]()
	i.seen[val.Value()] = struct{}{}
	return val.Value()
}

func (i *intersectingIter[T]) AsIter() Iter[T] {
	return i
}

// Difference returns the difference between two collections.
// The returned slice are the elements that are in the first collection but not in the second
func Difference[T comparable](first AsIter[T], second AsIter[T]) AsIter[T] {
	var result []T
	seen := make(map[T]struct{})

	left := first.AsIter()
	right := second.AsIter()

	for right.HasNext() {
		elem := right.Next()
		seen[elem] = struct{}{}
	}

	for left.HasNext() {
		elem := left.Next()
		if _, ok := seen[elem]; !ok {
			result = append(result, elem)
		}
	}
	return Slice(result...)
}

// SymmetricDifference removes the overlap between the 2 collections
// The returned slice are the elements that are in either the first or the second collection,
// but not in both
func SymmetricDifference[T comparable](first AsIter[T], second AsIter[T]) AsIter[T] {
	var result []T
	var leftbuf []T
	var rightbuf []T

	seenLeft := make(map[T]struct{})
	seenRight := make(map[T]struct{})

	left := first.AsIter()
	right := second.AsIter()

	for left.HasNext() {
		elem := left.Next()
		leftbuf = append(leftbuf, elem)
		seenLeft[elem] = struct{}{}
	}

	for right.HasNext() {
		elem := right.Next()
		rightbuf = append(rightbuf, elem)
		seenRight[elem] = struct{}{}
	}

	for _, v := range leftbuf {
		if _, ok := seenRight[v]; !ok {
			result = append(result, v)
		}
	}

	for _, v := range rightbuf {
		if _, ok := seenLeft[v]; !ok {
			result = append(result, v)
		}
	}

	return Slice(result...)
}

// Union returns all distinct elements from given collections.
// result returns will not change the order of elements relatively.
func Union[T comparable](left AsIter[T], right AsIter[T], others ...AsIter[T]) AsIter[T] {
	return &unionIter[T]{
		seen:    make(map[T]struct{}),
		current: Concat(left, right, others...),
	}
}

type unionIter[T comparable] struct {
	seen    map[T]struct{}
	current Iter[T]
	notSeen Option[T]
}

func (i *unionIter[T]) HasNext() bool {
	for i.current.HasNext() {
		elem := i.current.Next()
		if _, seen := i.seen[elem]; !seen {
			i.seen[elem] = struct{}{}
			i.notSeen = Some(elem)
			return true
		}
	}
	return false
}

func (i *unionIter[T]) Next() T {
	res := i.notSeen
	i.notSeen = None[T]()
	return res.Value()
}

func (i *unionIter[T]) AsIter() Iter[T] {
	return i
}
