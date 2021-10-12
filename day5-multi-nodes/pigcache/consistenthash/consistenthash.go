package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32.
// Hash 将字节数组映射成uint32。
type Hash func(data []byte) uint32

// Map contains all hashed keys.
// Map 包含所有的哈希键。
type Map struct {
	hash Hash
	replicas int
	keys []int // sorted
	hashMap map[int]string
}

// New creates a map instance.
// New 创建一个Map实例。
func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		keys:     nil,
		hashMap:  make(map[int]string),
	}

	if m.hash == nil {
		// 指定一个hash算法。
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

// Add adds some keys to the map.
// Add 添加一些键到map中。
func (m *Map) Add(keys ...string)  {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}

	sort.Ints(m.keys)
}

// Get gets the closest item in the hash to the provided key.
// Get 获得哈希环中与给定键最接近的节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// Binary search for appropriate replica
	// 二分查找合适的副本
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx % len(m.keys)]]
}

