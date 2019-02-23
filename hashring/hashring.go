// Package hashring implement a consistent hashing algorithm.
package hashring

import (
	"hash"
	"sync"
)

type node struct {
	Node
	addr []byte
	key  uint64

	replicas *[]*node
}

// HashRing implement a consistent hashing, allowing to add and remove
// nodes while maintaining the routing for the other nodes.
type HashRing struct {
	m sync.Mutex

	numReplicas uint16
	hasher      hash.Hash64
	ring        Tree
}

// New instantiate the HashRing, increasing the number of replicas will increase
// the consistency of the returned value from Get.
func New(tree Tree, numReplicas uint16, hasher hash.Hash64) *HashRing {
	return &HashRing{
		hasher:      hasher,
		numReplicas: numReplicas,
		ring:        tree,
	}
}

func (h *HashRing) hash(data []byte) uint64 {
	h.m.Lock()
	defer h.m.Unlock()

	h.hasher.Reset()
	h.hasher.Write(data)
	return h.hasher.Sum64()
}

// Add a node in the cluster
func (h *HashRing) Add(addr []byte) {
	var replicas []*node
	for i := uint16(0); i < h.numReplicas; i++ {
		data := append(addr)
		if i > 0 {
			h := uint8(i >> 8)
			t := uint8(i & 0xff)
			data = append(data, h, t)
		}
		node := &node{
			addr:     addr,
			key:      h.hash(data),
			replicas: &replicas,
		}
		h.ring.Put(node.key, node)
		replicas = append(replicas, node)
	}
}

// Remove a node from the cluster
func (h *HashRing) Remove(addr []byte) {
	key := h.hash(addr)
	n := h.ring.Get(key)
	if n == nil {
		return
	}

	node := n.(*node)
	for _, replica := range *node.replicas {
		h.ring.Remove(replica.key)
	}
}

// Get will consistently retrieve an added node for the same given key,
// as long as the node hasn't been removed.
func (h *HashRing) Get(key []byte) ([]byte, error) {
	k := h.hash(key)
	root := h.ring.Root()
	if closestNode := search(k, root); closestNode != nil {
		v := closestNode.Value().(*node)
		return v.addr, nil
	}
	r := root.Value().(*node)
	return r.addr, nil
}
