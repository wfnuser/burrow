package consistent

import (
	"hash/crc32"
	"sort"
)

// HashRing contains the hosts and virtual hosts
type HashRing struct {
	nodes      map[uint32]string
	replicates int
	keys       []uint32
}

// New hashring
func New(replicates int) *HashRing {
	hashRing := &HashRing{
		replicates: replicates,
		nodes:      make(map[uint32]string),
	}

	return hashRing
}

// Add add nodes and virtual nodes with key to hash ring
func (hashRing *HashRing) Add(key string) {
	for i := 0; i < hashRing.replicates; i++ {
		hash := crc32.ChecksumIEEE([]byte(key + "-" + string(i)))
		hashRing.keys = append(hashRing.keys, hash)
		hashRing.nodes[hash] = key
	}
	sort.Slice(hashRing.keys, func(i, j int) bool { return hashRing.keys[i] < hashRing.keys[j] })
}

// Get the closest node to key on the hash ring
func (hashRing *HashRing) Get(key string) string {
	hash := crc32.ChecksumIEEE([]byte(key))
	idx := sort.Search(len(hashRing.keys), func(i int) bool { return hashRing.keys[i] >= hash })
	if idx == len(hashRing.keys) {
		idx = 0
	}

	return hashRing.nodes[hashRing.keys[idx]]
}
