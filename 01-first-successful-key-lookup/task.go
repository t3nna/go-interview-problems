package main

import (
	"context"
)

type Getter interface {
	Get(ctx context.Context, address, key string) (string, error)
}

// Call `Getter.Get()` for each address in parallel.
// Returns the first successful response.
// If all requests fail, returns an error.
func Get(ctx context.Context, getter Getter, addresses []string, key string) (string, error) {
	res := make(chan string)
	errorSig := make(chan error)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if len(addresses) == 0 {
		return "", nil
	}
	for idx, address := range addresses {
		go func() {
			val, err := getter.Get(ctx, address, key)
			if err != nil {
				if idx == len(addresses)-1 {
					errorSig <- err
				}
				return
			}
			res <- val
		}()
	}

	select {
	case err := <-errorSig:
		return "", err
	case success := <-res:
		return success, nil
	case <-ctx.Done():
		return "", context.Canceled
	}
}
