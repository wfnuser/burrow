package burrow

import (
	"burrow/lru"
	"sync"
)

type Burrow struct {
	namespace string
	burrow    cache
}

var (
	burrows = make(map[string]*Burrow)
	mu      sync.Mutex
)

// NewBurrow used to create a burrow with namespace
func NewBurrow(namespace string, capacity int) *Burrow {
	mu.Lock()
	defer mu.Unlock()

	b := &Burrow{
		namespace: namespace,
		burrow:    cache{capacity: capacity},
	}

	burrows[namespace] = b
	return b
}

// GetBurrow used to get the burrow by namespace
func GetBurrow(namespace string) *Burrow {
	mu.Lock()
	defer mu.Unlock()

	b := burrows[namespace]
	return b
}

// Get used to get key-value from burrow
func (b *Burrow) Get(key string) (value lru.Value, ok bool) {
	return b.burrow.get(key)
}

// Delete used to delete key-value from burrow
func (b *Burrow) Delete(key string) {
	b.burrow.delete(key)
	return
}

// Put used to put key-value to burrow
func (b *Burrow) Put(key string, value lru.Value) {
	b.burrow.put(key, value)
	return
}
