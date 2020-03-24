package main

import (
	"burrow"
	"burrow/lru"
	"log"
	"net/http"
)

var db = map[string]string{
	"6.824": "MIT",
	"15213": "CMU",
	"15445": "CMU",
}

func main() {
	burrow.NewBurrow("test", 5, burrow.FuncGetter(
		func(key string) (lru.Value, bool) {
			log.Println("Fetch data from datasource by: ", key)
			if v, ok := db[key]; ok {
				return v, true
			}
			return nil, false
		}))

	server := burrow.NewHTTPPool()
	http.ListenAndServe("localhost:9999", server)
}
