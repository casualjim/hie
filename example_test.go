package hie_test

import (
	"fmt"

	"github.com/casualjim/hie"
	"github.com/casualjim/hie/iter"
	"github.com/casualjim/hie/iterable"
)

func ExampleFlatMap() {
	slices := hie.Slice(hie.Slice(1, 2).AsIter(), hie.Slice(3, 4).AsIter(), hie.Slice(5, 6).AsIter())

	it := iterable.FlatMap[hie.Iter[int]](slices, func(i hie.Iter[int]) hie.Iter[string] {
		return iter.Map(i, func(i int) string { return fmt.Sprintf("%d", i) })
	})

	fmt.Printf("%#v\n", iter.Collect(it))
	// Output: []string{"1", "2", "3", "4", "5", "6"}
}
