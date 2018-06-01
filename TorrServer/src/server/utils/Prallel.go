package utils

import (
	"sync"
)

func ParallelFor(f, t int, fn func(i int)) {
	var wg sync.WaitGroup
	wg.Add(t - f)
	for i := f; i < t; i++ {
		go func(i int) {
			fn(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
