package burrow

import (
	"burrow/lru"
	"sync"
)

// Burrow is the cache type
type Burrow struct {
	namespace string
	burrow    cache
	getter    Getter
}

// Getter load data from data source
type Getter interface {
	Get(key string) (lru.Value, bool)
}

// FuncGetter implements Getter
type FuncGetter func(key string) (lru.Value, bool)

// Get implements Getter's Get Method
func (f FuncGetter) Get(key string) (lru.Value, bool) {
	return f(key)
}

var (
	burrows = make(map[string]*Burrow)
	mu      sync.Mutex
)

// NewBurrow used to create a burrow with namespace
func NewBurrow(namespace string, capacity int, getter Getter) *Burrow {
	mu.Lock()
	defer mu.Unlock()

	b := &Burrow{
		namespace: namespace,
		getter:    getter,
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
	v, ok := b.burrow.get(key)
	if ok {
		return v, ok
	}

	originValue, originOk := b.getter.Get(key)
	if originOk {
		b.burrow.put(key, originValue)
		return originValue, originOk
	} else {
		return nil, false
	}

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
