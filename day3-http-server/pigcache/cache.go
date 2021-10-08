package pigcache

import (
	"Pigcache/day1-lru/pigcache/lru"
	"sync"
)

type cache struct {
	mu 			sync.Mutex
	lru 		*lru.Cache
	cacheBytes 	int64
}

// add adds key-value to lru-cache
// add 函数添加键值对数据到lru缓存中
func (c *cache) add(key string, value ByteView)  {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

// get gets key-value from lru-cache
// get从lru缓存获取键值对数据
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return // Lazy Initialization well or bad?
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}

