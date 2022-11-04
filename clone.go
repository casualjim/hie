package hie

type ClonableIter[T any] interface {
	Clone() Iter[T]
}

type Clonable[T any] interface {
	Clone() T
}

func Clone[T Clonable[T]](value T) T { return value.Clone() }
