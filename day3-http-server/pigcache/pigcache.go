package pigcache

import (
	"fmt"
	"log"
	"sync"
)

// A Getter loads data for a key
// Getter 使用加载指定Key的Value
type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function
// GetterFunc 函数类型实现了Getter接口
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
// Get函数实现Getter接口
func (f GetterFunc) Get(key string) ([]byte, error)  {
	return f(key)
}

// A Group is a cache namespace and associated data loaded spread over
// Group是一个缓存的命名空间
type Group struct {
	name string
	getter Getter
	mainCache cache
}

var (
	mu sync.Mutex
	groups = make(map[string]*Group)
)

// NewGroup create a new instance of Group
// NewGroup 创建一个Group的实例
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}

	groups[name] = g
	return g
}

// GetGroup returns the named group previously created with NewGroup, or nil if there's no such group.
// GetGroup 返回指定名字的group
func GetGroup(name string) *Group {
	mu.Lock()
	defer mu.Unlock()
	g := groups[name]
	return g
}


// Get value for a key from cache
// Get从cache中获取指定key的value
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[PigCache hit]")
		return v, nil
	}

	return g.load(key)
}

// load value for key from local database
// load从本地数据库获取指定key的value
func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

//getLocally get for key from local
func (g *Group) getLocally(key string)  (ByteView, error) {
	println("从数据库加载数据")
	bytes, err := g.getter.Get(key) // 调用回调函数从数据库数据
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// populateCache add key-value to mainCache
func (g *Group) populateCache(key string, value ByteView)  {
	g.mainCache.add(key, value)
}

