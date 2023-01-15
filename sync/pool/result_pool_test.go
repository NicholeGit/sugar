package pool

import (
	"fmt"
	"sort"
)

func ExampleResultPool() {
	p := NewWithResults[int]()
	for i := 0; i < 10; i++ {
		i := i
		p.Go(func() int {
			return i * 2
		})
	}
	res := p.Wait()
	// Result order is nondeterministic, so sort them first
	sort.Ints(res)
	fmt.Println(res)

	// Output:
	// [0 2 4 6 8 10 12 14 16 18]
}
