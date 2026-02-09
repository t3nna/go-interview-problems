package main

import "sync"

func Merge(channels ...<-chan int) <-chan int {

	res := make(chan int)

	var wg sync.WaitGroup
	for _, ch := range channels {
		wg.Add(1)
		go func(ch <-chan int) {
			for intChan := range ch {
				res <- intChan
			}
			wg.Done()
		}(ch)
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	return res
}
