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
