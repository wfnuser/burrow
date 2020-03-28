package burrow

import (
	"burrow/consistent"
	"burrow/lru"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const defaultBasePath = "/burrow/"
const defaultReplicates = 3

// HTTPPool implements PeerPicker for a pool of HTTP peers.
type HTTPPool struct {
	basePath string
	serverID string   // server identification
	servers  []string // other servers
	hashRing *consistent.HashRing
}

// NewHTTPPool initializes an HTTP pool.
func NewHTTPPool(serverID string) *HTTPPool {
	return &HTTPPool{
		basePath: defaultBasePath,
		serverID: serverID,
	}
}

// NewHTTPPoolWithServers initializes an HTTP pool with servers.
func NewHTTPPoolWithServers(serverID string, servers []string) *HTTPPool {
	hashRing := consistent.New(defaultReplicates)
	for i := 0; i < len(servers); i++ {
		server := servers[i]
		hashRing.Add(server)
	}
	return &HTTPPool{
		basePath: defaultBasePath,
		serverID: serverID,
		servers:  servers,
		hashRing: hashRing,
	}
}

func fetchRemote(namespace string, key string, server string) (value lru.Value, ok bool) {
	fmt.Println("start fetch data")
	u := fmt.Sprintf(
		"http://%v%v%v/%v",
		server,
		defaultBasePath,
		url.QueryEscape(namespace),
		url.QueryEscape(key),
	)
	res, err := http.Get(u)
	if err != nil {
		return nil, false
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, false
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, false
	}

	//TODO: should use bytes instead of string...
	return string(bytes), true
}

// TODO: should remove w in fetch method
func (p *HTTPPool) fetch(key string, namespace string, w http.ResponseWriter) (value lru.Value, ok bool) {
	burrow := GetBurrow(namespace)

	if burrow == nil {
		http.Error(w, "no such burrow: "+namespace, http.StatusNotFound)
		return
	}

	server := p.hashRing.Get(key)
	fmt.Printf("%v\n", server)
	if server != p.serverID {
		value, ok = fetchRemote(namespace, key, server)
		fmt.Printf("%v\n", value)
		if !ok {
			http.Error(w, "remote server error", http.StatusInternalServerError)
			return
		}
		return value, ok
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

	value, _ := p.fetch(key, namespace, w)

	str := fmt.Sprint(value)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write([]byte(str))
}
