package burrow

import (
	"burrow/lru"
	"sync"
)

type cache struct {
	mu       sync.Mutex
	lru      *lru.Cache
	capacity int
}

func (c *cache) put(key string, value lru.Value) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = lru.New(c.capacity)
	}

	c.lru.Put(key, value)
}

func (c *cache) delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return
	}

	c.lru.Delete(key)
}

func (c *cache) get(key string) (value lru.Value, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return
	}

	return c.lru.Get(key)
}
