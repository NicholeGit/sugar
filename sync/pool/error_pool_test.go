package pool

import (
	"fmt"

	"github.com/NicholeGit/sugar/errors"
)

func ExampleErrorPool() {
	p := New().WithErrors()
	for i := 0; i < 3; i++ {
		i := i
		p.Go(func() error {
			if i == 2 {
				return errors.New("oh no!")
			}
			return nil
		})
	}
	err := p.Wait()
	fmt.Println(err)
	// Output:
	// [[ExampleErrorPool.func1] (error_pool_test.go#15) oh no!]
}
