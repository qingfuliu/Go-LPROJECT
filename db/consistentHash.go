package db

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"sort"
	"strconv"
	"sync"
)

type sortedHashes []uint32

func GetNewConsistentNode() ConsistentHash {
	return &consistentHash{
		NodeInformation:  make(map[uint32]interface{}),
		CountVirtualNode: 20,
		SortedHashes:     make([]uint32, 20),
		mutex:            sync.Mutex{},
	}
}

func (s sortedHashes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortedHashes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortedHashes) Len() int {
	return len(s)
}

type consistentHash struct {
	NodeInformation  map[uint32]interface{}
	CountVirtualNode int
	SortedHashes     sortedHashes
	mutex            sync.Mutex
}

type ConsistentHash interface {
	//获取服务器信息，传入的参数可能是结构体等
	Get(key interface{}) interface{}
	//添加新的服务器节点
	Add(key string) error
	//删除节点
	Remove(key string)
	//hash算法
	hash(key string) uint32
	//为random生成随机key
	generateRandomKey(key string, i int) uint32
	//重新排序节点信息
	updateSortedHashes()
}

func (c *consistentHash) generateRandomKey(key string, i int) uint32 {
	var randomKey uint32
	for {
		randomKey = c.hash(key + strconv.Itoa(i))
		if _, ok := c.NodeInformation[randomKey]; ok == false {
			break
		}
	}
	return randomKey
}
func (c *consistentHash) hash(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *consistentHash) Get(key interface{}) interface{} {
	var str string
	switch key.(type) {
	case fmt.Stringer:
		str = key.(fmt.Stringer).String()
	case string:
		str = key.(string)
	default:
		for _, value := range c.NodeInformation {
			return value
		}
	}
	hashKey := c.hash(str)
	return c.searchNode(hashKey)
}

func (c *consistentHash) searchNode(hashkey uint32) interface{} {
	nodeIndex := sort.Search(c.CountVirtualNode, func(n int) bool {
		return c.SortedHashes[n] >= hashkey
	})
	if nodeIndex >= c.CountVirtualNode {
		nodeIndex = rand.Intn(c.CountVirtualNode)
	}
	return c.NodeInformation[c.SortedHashes[nodeIndex]]
}

func (c *consistentHash) Add(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.add(key)
}

func (c *consistentHash) Remove(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Remove(key)
}

func (c *consistentHash) add(key string) error {
	//需要进行节点的同步
	for index := 0; index < c.CountVirtualNode; index++ {
		randomKey := c.generateRandomKey(key, index)
		//node searchNode(randomKey) sync(node,key)
		c.NodeInformation[randomKey] = key
	}
	c.updateSortedHashes()
	return nil
}

func (c *consistentHash) updateSortedHashes() {
	if len(c.SortedHashes) > len(c.NodeInformation)*2 {
		c.SortedHashes = make([]uint32, 0, len(c.NodeInformation))
	}
	for key, _ := range c.NodeInformation {
		c.SortedHashes = append(c.SortedHashes, key)
	}
	sort.Sort(c.SortedHashes)
}

func (c *consistentHash) remove(key string) {
	for index := 0; index < c.CountVirtualNode; index++ {
		randomKey := c.generateRandomKey(key, index)
		//node searchNode(ranomKey) sync(node,key)
		delete(c.NodeInformation, randomKey)
	}
	c.updateSortedHashes()
}
