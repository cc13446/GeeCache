package GeeCache

import (
	"GeeCache/lru"
	"sync"
)

type cache struct {
	mutex      sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key KeyView, value ValueView) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key KeyView) (value ValueView, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ValueView), ok
	}

	return
}
