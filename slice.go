package hie

func Slice[T any](value ...T) AsIter[T] {
	return &sliceIter[T]{
		under: value,
		idx:   -1,
	}
}

type sliceIter[T any] struct {
	under []T
	idx   int
}

func (s *sliceIter[T]) HasNext() bool {
	return s.idx < len(s.under)-1
}

func (s *sliceIter[T]) Next() T {
	if !s.HasNext() {
		panic("iterating beyond end")
	}

	s.idx++
	return s.under[s.idx]
}

func (s *sliceIter[T]) Collect() []T {
	result := make([]T, len(s.under))
	copy(result, s.under)
	return result
}

func (s *sliceIter[T]) AsIter() Iter[T] {
	return &sliceIter[T]{
		under: s.under,
		idx:   -1,
	}
}
