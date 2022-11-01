package iterable

import (
	"github.com/casualjim/hie"
	"github.com/casualjim/hie/iter"
	"github.com/casualjim/hie/option"
)

type AsIter[T any] interface {
	AsIter() hie.Iter[T]
}

func ForEach[T any](it AsIter[T], fn hie.Iterator[T]) {
	iter.ForEach(it.AsIter(), fn)
}

func FlattenAsIter[T any](it AsIter[AsIter[T]]) hie.Iter[T] {
	mapper := func(i AsIter[T]) hie.Iter[T] { return i.AsIter() }
	return iter.Flatten(iter.Map(it.AsIter(), mapper))
}

func MapAsIter[T any](it AsIter[AsIter[T]]) hie.Iter[hie.Iter[T]] {
	mapper := func(i AsIter[T]) hie.Iter[T] { return i.AsIter() }
	return iter.Map(it.AsIter(), mapper)
}

func Flatten[T any](it AsIter[hie.Iter[T]]) hie.Iter[T] {
	return iter.Flatten(it.AsIter())
}

func FilterMap[T, R any](it AsIter[T], fn iter.FilterMapper[T, R]) hie.Iter[R] {
	return iter.FilterMap(it.AsIter(), fn)
}

func Map[T, R any](it AsIter[T], fn iter.Mapper[T, R]) hie.Iter[R] {
	return iter.Map(it.AsIter(), fn)
}

func FlatMap[T, R any](it AsIter[T], fn iter.FlatMapper[T, R]) hie.Iter[R] {
	return iter.FlatMap(it.AsIter(), fn)
}

func Filter[T any](it AsIter[T], predicate iter.Predicate[T]) hie.Iter[T] {
	return iter.Filter(it.AsIter(), predicate)
}

func Collect[T any](it AsIter[T]) []T {
	return iter.Collect(it.AsIter())
}

func Concat[T any](left AsIter[T], right AsIter[T], others ...AsIter[T]) hie.Iter[T] {
	oa := make([]hie.Iter[T], len(others))
	for i, ai := range others {
		oa[i] = ai.AsIter()
	}
	return iter.Concat(left.AsIter(), right.AsIter(), oa...)
}

func Fold[A, T any](it AsIter[T], initialValue A, folder iter.AccumulatorLeft[A, T]) A {
	return iter.Fold(it.AsIter(), initialValue, folder)
}

func TakeN[T any](it AsIter[T], n int) hie.Iter[T] {
	return iter.TakeN(it.AsIter(), n)
}

func Find[T any](it AsIter[T], predicate iter.Predicate[T]) option.Option[T] {
	return iter.Find(it.AsIter(), predicate)
}
