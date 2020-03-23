package burrow

import (
	"fmt"
	"net/http"
	"strings"
)

const defaultBasePath = "/burrow/"

// HTTPPool implements PeerPicker for a pool of HTTP peers.
type HTTPPool struct {
	basePath string
}

// NewHTTPPool initializes an HTTP pool of peers.
func NewHTTPPool() *HTTPPool {
	return &HTTPPool{
		basePath: defaultBasePath,
	}
}

// ServeHTTP handle all http requests
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}

	fmt.Println("%s %s", r.Method, r.URL.Path)

	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	namespace := parts[0]
	key := parts[1]

	burrow := GetBurrow(namespace)
	if burrow == nil {
		http.Error(w, "no such burrow: "+namespace, http.StatusNotFound)
		return
	}

	value, ok := burrow.Get(key)
	if !ok {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(value.([]byte))
}
