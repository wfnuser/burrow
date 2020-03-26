package burrow

import (
	"burrow/consistent"
	"burrow/lru"
	"fmt"
	"net/http"
	"strings"
)

const defaultBasePath = "/burrow/"
const defaultReplicates = 3

// HTTPPool implements PeerPicker for a pool of HTTP peers.
type HTTPPool struct {
	basePath string
	serverID string   // server identification
	servers  []string // other servers
	hashRing consistent.HashRing
}

// NewHTTPPool initializes an HTTP pool of peers.
func NewHTTPPool(serverID string) *HTTPPool {
	return &HTTPPool{
		basePath: defaultBasePath,
		serverID: serverID,
	}
}

// // NewHTTPPool initializes an HTTP pool of peers.
// func NewHTTPPool(serverID string, servers string[]) *HTTPPool {
// 	return &HTTPPool{
// 		basePath: defaultBasePath,
// 		serverID: serverID,
// 	}
// }

// TODO: should remove w in fetch method
func fetch(key string, namespace string, w http.ResponseWriter) (value lru.Value, ok bool) {
	burrow := GetBurrow(namespace)
	if burrow == nil {
		http.Error(w, "no such burrow: "+namespace, http.StatusNotFound)
		return
	}

	value, ok = burrow.Get(key)
	if !ok {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	return value, ok
}

// ServeHTTP handle all http requests
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}

	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	namespace := parts[0]
	key := parts[1]

	value, _ := fetch(key, namespace, w)

	str := fmt.Sprint(value)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write([]byte(str))
}
