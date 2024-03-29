package hie

import (
	"io"
)

func Identity[T any](input T) T { return input }

func Close(closeable io.Closer) error { return closeable.Close() }

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
	return SliceAsIter[T](append([]T(nil), []T(s)...))
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

func (s *sliceIter[T]) Clone() Iter[T] {
	// et := reflect.TypeOf(s.under).Elem()
	// under := s.under
	// if et.AssignableTo(clonableIterT) {

	// }
	return &sliceIter[T]{
		under: s.under,
		idx:   -1,
	}
}

// var clonableIterT reflect.Type

// func init() {
// 	clonableIterT = reflect.TypeOf([0]ClonableIter[any]{}).Elem()
// }
