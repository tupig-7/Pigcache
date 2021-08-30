package lru

import (
	"container/list"
	"fmt"
)

// Cache is a LRU cache. It is not safe for concurrent access
// 缓存是LRU缓存，对于并发进程是不安全的
type Cache struct {
	// Maximum memory allowed
	// 允许使用的最大内存
	maxBytes int64
	// Used memory
	// 已经使用的内存
	nbytes int64
	ll     *list.List
	cache  map[string]*list.Element
	// optional and executed when an entry is purged
	// 可选项，在清除记录时执行，节点移除时的回调函数
	OnEvicted func(key string, value Value)
}

// Data type of bidirectional linked list node
// 双向链表节点的数据类型
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
// Value 使用Len去记录使用了多少字节
type Value interface {
	Len() int
}

// New is the Constructor of Cache
// New() 是Cache的构造函数
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get look ups a key's value from the cache
// 从cache中获取一个键的值
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {

		//Move the most recently visited element to the end of the list
		// 将最近访问的元素移动到链表末尾
		c.ll.MoveToBack(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item
// RemoveOldest 移除最近最少访问的节点（链表尾）
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add adds a value to the cache
// 添加一个值到cache
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)

		// Using new value cover old value
		// 使用新值覆盖旧值
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	// Remove the node until the memory used is less than or equal to the maximum memory
	// 移除节点直到已经使用内存小于等于最大内存
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
// Len cache中entry的个数
func (c *Cache) Len() int {
	return c.ll.Len()
}
