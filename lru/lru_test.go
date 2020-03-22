package lru

import (
	"fmt"
	"testing"
)

type String string

func TestPutAndGet(t *testing.T) {
	lru := New(10)
	lru.Put("key", "test value")

	if v, ok := lru.Get("key"); !ok || string(v.(string)) != "test value" {
		t.Fatalf("cache hit key=key failed")
	}
}

func TestPutAndGetAfterCapacityFull(t *testing.T) {
	lru := New(2)
	lru.Put("key", "test value")
	fmt.Printf("%v\n", lru.GetEntriesNumber())
	lru.Put("key1", "test value1")
	fmt.Printf("%v\n", lru.GetEntriesNumber())
	lru.Put("key2", "test value2")
	fmt.Printf("%v\n", lru.GetEntriesNumber())

	if _, ok := lru.Get("key"); ok {
		t.Fatalf("cache hit key=key failed")
	}
	if v, ok := lru.Get("key1"); !ok || string(v.(string)) != "test value1" {
		t.Fatalf("cache hit key=key1 failed")
	}

	lru.Delete("key1")
	fmt.Printf("%v\n", lru.GetEntriesNumber())
}
