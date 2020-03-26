package burrow

import (
	"burrow/lru"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
)

type String string

var db = map[string]string{
	"6.824": "MIT",
	"15213": "CMU",
	"15445": "CMU",
}

func TestGetBurrow(t *testing.T) {
	loadCounts := make(map[string]int, len(db))

	b := NewBurrow("test", 2, FuncGetter(
		func(key string) (lru.Value, bool) {
			log.Println("Fetch data from datasource by: ", key)
			if v, ok := db[key]; ok {
				loadCounts[key]++
				return v, true
			}
			return nil, false
		}))

	b.Get("6.824")
	b.Get("6.824")
	b.Get("15213")
	b.Get("15445")
	// nb := GetBurrow("test")
	// fmt.Printf("namespace is: %s \n", nb.namespace)
}

func increment(b *Burrow, key string) {
	v, ok := b.Get(key)

	if ok {
		b.Put(key, v.(int)+1)
	}
}

// won't fail; you should compare it to lru_test.go
func TestPutBurrow(t *testing.T) {
	b := NewBurrow("test", 2, FuncGetter(
		func(key string) (value lru.Value, ok bool) { return }))
	b.Put("test", 0)
	for i := 0; i < 10000; i++ {
		go b.Put(strconv.Itoa(rand.Intn(10)), i)
		go b.Get(strconv.Itoa(rand.Intn(10)))
	}
	v, ok := b.Get("test")

	fmt.Printf("%v %v", v, ok)
}
