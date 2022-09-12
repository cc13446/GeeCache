package GeeCache

import (
	"fmt"
	"log"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	gee := NewGroup("scores", 2<<10, GetterFunc(
		func(key KeyView) ([]byte, error) {
			keyString := key.String()
			log.Println("[SlowDB] search key", keyString)
			if v, ok := db[keyString]; ok {
				if _, ok := loadCounts[keyString]; !ok {
					loadCounts[keyString] = 0
				}
				loadCounts[keyString] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		// load from callback function
		if view, err := gee.Get(FromString(k)); err != nil || view.String() != v {
			t.Fatal("failed to get value of Tom")
		}

		// cache hit
		if _, err := gee.Get(FromString(k)); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}

	if view, err := gee.Get(FromString("unknown")); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
