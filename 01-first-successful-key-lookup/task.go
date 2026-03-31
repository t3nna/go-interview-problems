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
	if ctx.Err() != nil {
		return "", errors.New("")
	}

	ctxCancel, cancel := context.WithCancel(ctx)

	defer cancel()

	errorCh := make(chan error, len(addresses))
	valCh := make(chan string, len(addresses))

	var wg sync.WaitGroup
	go func() {
		wg.Wait()
		close(errorCh)
		close(valCh)
	}()

	for _, address := range addresses {
		wg.Add(1)
		go func(addr string) {

			defer wg.Done()

			val, err := getter.Get(ctxCancel, addr, key)

			if err != nil {
				errorCh <- err
				return
			}

			valCh <- val

		}(address)
	}

	errCount := 0
	for {
		select {
		case val := <-valCh:
			cancel()
			return val, nil
		case <-errorCh:
			errCount++
			if errCount == len(addresses) {
				return "", errors.New("not found")
			}
		case <-ctxCancel.Done():
			return "", ctxCancel.Err()
		}

	}
}
