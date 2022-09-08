package lru

import "container/list"

// Cache is an LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes  int64
	usedBytes int64
	list      *list.List
	cache     map[Key]*list.Element
	// 某条记录被删除时的回调
	OnEvicted func(key Key, value Value)
}

// entry
type entry struct {
	key   Key
	value Value
}

// Key is comparable
// Key use Len to count how many bytes it takes
type Key interface {
	Len() int
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// New is the Constructor of Cache
func New(maxBytes int64, OnEvicted func(Key, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		usedBytes: 0,
		list:      list.New(),
		cache:     make(map[Key]*list.Element),
		OnEvicted: OnEvicted,
	}
}

// Get look ups a key's value
func (c *Cache) Get(key Key) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.list.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	if ele := c.list.Back(); ele != nil {
		c.list.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.usedBytes -= int64(kv.key.Len()) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(key Key, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.list.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.usedBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.list.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.usedBytes += int64(key.Len()) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.list.Len()
}
