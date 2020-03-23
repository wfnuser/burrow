package main

import (
	"burrow"
	"net/http"
)

func main() {

	server := burrow.NewHTTPPool()
	http.ListenAndServe("localhost:9999", server)
}
