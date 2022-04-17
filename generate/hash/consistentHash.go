package hash

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

const (
	minVirtualNodes = 100
	maxTopWeight    = 100
	prime           = "16777619"
)

type (
	Func func(str string) uint64

	ConsistentHash1 struct {
		hashFunc         Func
		keys             []uint64
		virtualNodeCount int
		virtualMaps      map[uint64][]interface{}
		nodes            map[string]struct{}
		mutex            sync.RWMutex
	}
)

func (c *ConsistentHash1) Add(node interface{}, virtualMaps int) {
	c.Remove(node)
	nodeStr := nodeToString(node)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for index := 0; index < virtualMaps; index++ {
		key := c.hashFunc(nodeStr + strconv.Itoa(index))
		c.keys = append(c.keys, key)
		c.virtualMaps[key] = append(c.virtualMaps[key], node)
	}
	c.addNode(nodeStr)
	sort.Slice(c.keys, func(i, j int) bool {
		return c.keys[i] < c.keys[j]
	})
}

func (c *ConsistentHash1) AddWithVirtualNodes(node interface{}, virtualMaps int) {
	if virtualMaps < minVirtualNodes {
		virtualMaps = minVirtualNodes
	}
	c.Add(node, virtualMaps)
}

func (c *ConsistentHash1) AddWithWrights(node interface{}, weights int) {
	virtualMaps := c.virtualNodeCount * weights / maxTopWeight
	c.AddWithVirtualNodes(node, virtualMaps)
}

func (c *ConsistentHash1) Get(v interface{}) (interface{}, bool) {
	if len(c.nodes) < 1 {
		return nil, false
	}
	vStr := nodeToString(v)
	vhashKey := c.hashFunc(vStr)
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	index := sort.Search(len(c.keys), func(i int) bool {
		return c.keys[i] >= vhashKey
	}) % len(c.keys)
	nodes := c.virtualMaps[c.keys[index]]
	switch len(nodes) {
	case 0:
		return nil, false
	case 1:
		return nodes[0], true
	default:
		vhashKey = c.hashFunc(vStr + prime)
		pos := int(vhashKey % uint64(len(nodes)))
		return nodes[pos], true
	}
}

func (c *ConsistentHash1) Remove(node interface{}) {
	nodeStr := nodeToString(node)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.nodes[nodeStr]; !ok {
		return
	}

	for i := 0; i < c.virtualNodeCount; i++ {
		key := c.hashFunc(nodeStr + strconv.Itoa(i))
		index := sort.Search(len(c.keys), func(i int) bool {
			return c.keys[i] >= key
		})
		if c.keys[index] == key {
			c.keys = append(c.keys[1:index], c.keys[index+1:]...)
		}

	}

	c.removeNode(nodeStr)
}

func (c *ConsistentHash1) removeVirtualNodes(key uint64,nodeStr string) {
	if nodes, ok := c.virtualMaps[key]; ok {
		temp := nodes[:0]
		for index,_:=range nodes{
			if nodeToString(nodes[index])!=nodeStr{
				temp=append(temp,nodes[index])
			}
		}
		if len(temp)>0{
			c.virtualMaps[key]=temp
		}else{
			delete(c.virtualMaps,key)
		}
	}
}

func nodeToString(node interface{}) string {
	switch node.(type) {
	case string:
		return node.(string)
	case fmt.Stringer:
		return node.(fmt.Stringer).String()
	}
	return ""
}

func (c *ConsistentHash1) removeNode(nodeStr string) {
	delete(c.nodes, nodeStr)
}

func (c *ConsistentHash1) addNode(nodeStr string) {
	c.nodes[nodeStr] = struct{}{}
}

