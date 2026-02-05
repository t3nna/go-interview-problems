package main

import (
	"context"
	"errors"
	"sync"
)

type Getter interface {
	Get(ctx context.Context, address, key string) (string, error)
}

// Call `Getter.Get()` for each address in parallel.
// Returns the first successful response.
// If all requests fail, returns an error.
func Get(ctx context.Context, getter Getter, addresses []string, key string) (string, error) {
	if len(addresses) == 0 {
		return "", nil
	}

	// Create a context that we can cancel as soon as we get a success.
	// This propagates the cancellation to all other ongoing requests.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Use a buffered channel with a size equal to the number of addresses.
	// This prevents goroutine leaks by ensuring 'send' operations are non-blocking.
	resCh := make(chan string, len(addresses))

	// Use a WaitGroup to track all goroutines.
	var wg sync.WaitGroup
	wg.Add(len(addresses))

	for _, address := range addresses {
		// Fix: Pass 'address' as a parameter to avoid the common loop variable bug.
		go func(addr string) {
			defer wg.Done()

			val, err := getter.Get(ctx, addr, key)
			if err != nil {
				// On error, we just finish. If ALL fail, the wg.Wait() goroutine handles it.
				return
			}

			// First success wins.
			select {
			case resCh <- val:
				// Successfully sent result; cancel remaining requests.
				cancel()
			case <-ctx.Done():
				// Already finished or cancelled elsewhere.
			}
		}(address)
	}

	// This goroutine implements the "plan to stop"[cite: 658].
	// It waits for all workers to finish (either success or fail) and closes the channel.
	allDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(allDone)
	}()

	select {
	case result := <-resCh:
		return result, nil
	case <-allDone:
		// If we reach here, it means all goroutines called wg.Done()
		// without any of them sending a success to resCh.
		return "", errors.New("all requests failed")
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
