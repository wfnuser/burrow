package consistent

import (
	"fmt"
	"testing"
)

func TestAddAndGet(t *testing.T) {
	hashring := New(3)
	hashring.Add("192.168.0.1")
	hashring.Add("192.168.0.2")
	hashring.Add("192.168.0.3")

	fmt.Println(hashring.Get("6.824"))
	fmt.Println(hashring.Get("15213"))
	fmt.Println(hashring.Get("15445"))
}
