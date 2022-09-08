package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add(String("key1"), String("1234"))
	if v, ok := lru.Get(String("key1")); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get(String("key2")); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	max := len(k1 + k2 + v1 + v2)
	lru := New(int64(max), nil)
	lru.Add(String(k1), String(v1))
	lru.Add(String(k2), String(v2))
	lru.Add(String(k3), String(v3))

	if _, ok := lru.Get(String("key1")); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]Key, 0)
	callback := func(key Key, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add(String("k1"), String("12345678"))
	lru.Add(String("k2"), String("k2"))
	lru.Add(String("k3"), String("k3"))
	lru.Add(String("k4"), String("k4"))

	expect := []Key{String("k1"), String("k2")}
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s, but now we get %s", expect, keys)
	}
}
