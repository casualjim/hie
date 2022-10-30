package hie

type slice[T any] []T

func (s slice[T]) AsIter() Iter[T] {
	return &sliceIter[T]{
		under: []T(s),
		idx:   -1,
	}
}

func Slice[T any](value ...T) AsIter[T] {
	return slice[T](value)
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
