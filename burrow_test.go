package burrow

import (
	"fmt"
	"testing"
)

type String string

func TestGetBurrow(t *testing.T) {
	b := NewBurrow("test", 2)
	b.Put("key1", "value1")
	b.Put("key2", "value2")
	b.Put("key3", "value3")

	nb := GetBurrow("test")
	fmt.Printf("namespace is: %s \n", nb.namespace)
}

func increment(b *Burrow, key string) {
	v, ok := b.Get(key)

	if ok {
		b.Put(key, v.(int)+1)
	}
}

func TestPutBurrow(t *testing.T) {
	b := NewBurrow("test", 2)
	b.Put("test", 0)
	for i := 0; i < 100; i++ {
		go increment(b, "test")
	}
	v, ok := b.Get("test")

	fmt.Printf("%v %v", v, ok)
}
