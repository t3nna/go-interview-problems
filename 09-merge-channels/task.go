package main

import "sync"

func Merge(channels ...<-chan int) <-chan int {
	res := make(chan int)

	var wg sync.WaitGroup

	wg.Add(len(channels))
	for _, ch := range channels {
		go func(ch <-chan int) {
			defer wg.Done()
			for val := range ch {
				res <- val
			}

		}(ch)
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	return res
}
