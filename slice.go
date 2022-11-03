package hie

func Identity[T any](input T) T { return input }

type emptyIter[T any] struct{}

func (emptyIter[T]) HasNext() bool { return false }
func (emptyIter[T]) Next() T       { panic("next called on an empty iter") }

func EmptyIter[T any]() Iter[T] {
	return emptyIter[T]{}
}

type Iterator[T any] func(T) bool

type Iter[T any] interface {
	HasNext() bool
	Next() T
}

type SliceAsIter[T any] []T

func (s SliceAsIter[T]) AsIter() Iter[T] {
	return &sliceIter[T]{
		under: []T(s),
		idx:   -1,
	}
}

func (s SliceAsIter[T]) Clone() SliceAsIter[T] {
	return SliceAsIter[T](s)
}

func Slice[T any](value ...T) SliceAsIter[T] {
	return SliceAsIter[T](value)
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
