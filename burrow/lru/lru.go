package lru

import (
	"container/list"
)

type Cache struct {
	l        *list.List
	cache    map[string]*list.Element
	capacity int
}

type Value interface{}

type entry struct {
	key   string
	value Value
}

// New is the cache constructor
func New(capacity int) *Cache {
	return &Cache{
		l:        list.New(),
		cache:    make(map[string]*list.Element),
		capacity: capacity,
	}
}

// Put put kv to cache
func (c *Cache) Put(key string, value Value) {
	if element, ok := c.cache[key]; ok {
		c.l.MoveToFront(element)
		element.Value.(*entry).value = value
	} else {
		element := c.l.PushFront(&entry{key, value})
		c.cache[key] = element
	}

	if c.capacity < c.l.Len() {
		back := c.l.Back()
		c.l.Remove(back)
		delete(c.cache, back.Value.(*entry).key)
	}
}

// Get get key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		kv := element.Value.(*entry)
		return kv.value, true
	}
	return
}

// Delete delete key's value
func (c *Cache) Delete(key string) {
	if element, ok := c.cache[key]; ok {
		kv := element.Value.(*entry)
		c.l.Remove(element)
		delete(c.cache, kv.key)
		return
	}
	return
}

// GetCapacity get capacity
func (c *Cache) GetCapacity() int {
	return c.capacity
}

// GetEntriesNumber get number of entries
func (c *Cache) GetEntriesNumber() int {
	return c.l.Len()
}
