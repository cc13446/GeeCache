package GeeCache

import (
	"fmt"
	"log"
	"sync"
)

// A Getter loads data for a key
type Getter interface {
	Get(key KeyView) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(key KeyView) ([]byte, error)

// Get implements Getter interface function
// 函数类型实现一个接口，叫接口型函数，方便使用者在调用时既能够传入函数作为参数，也能够传入实现了该接口的结构体作为参数。
func (f GetterFunc) Get(key KeyView) ([]byte, error) {
	return f(key)
}

// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mutex  sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup create a new instance of Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mutex.Lock()
	defer mutex.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group.
func GetGroup(name string) *Group {
	mutex.RLock()
	defer mutex.RUnlock()
	g := groups[name]
	return g
}

// Get value for a key from cache
func (g *Group) Get(key KeyView) (ValueView, error) {
	if key.String() == "" {
		return ValueView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key KeyView) (value ValueView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key KeyView) (ValueView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ValueView{}, err

	}
	value := ValueView{data: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key KeyView, value ValueView) {
	g.mainCache.add(key, value)
}
