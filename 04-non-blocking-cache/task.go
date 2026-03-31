package main

import "sync"

type Client interface {
	Get(address string) (string, error)
}

type task struct {
	body  string
	err   error
	ready chan struct{}
}

type Cache struct {
	client Client
	// You can add new fields if needed
	cache map[string]*task
	mu    sync.Mutex
}

// Don't update signature of NewCache
func NewCache(client Client) *Cache {
	// TODO: Implement
	return &Cache{client: client,
		cache: make(map[string]*task),
	}
}

// Cache Client.Get result
func (c *Cache) Get(address string) (string, error) {
	// TODO: Implement. Right now it doesn't cache
	c.mu.Lock()
	res := c.cache[address]
	if res == nil {
		res = &task{ready: make(chan struct{})}
		c.cache[address] = res
		c.mu.Unlock()

		res.body, res.err = c.client.Get(address)
		close(res.ready)
	} else {
		c.mu.Unlock()
		<-res.ready
	}
	return res.body, res.err
}
