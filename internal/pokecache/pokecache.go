package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	data     map[string]cacheEntry
	interval time.Duration
	mu       sync.Mutex
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := cacheEntry{}
	entry.createdAt = time.Now()
	entry.val = val

	c.data[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.data[key]
	var val []byte

	if ok {
		val = entry.val
	}

	return val, ok
}

func (c *Cache) reapLoop() {

	ticker := time.NewTicker(c.interval)

	for range ticker.C {
		c.mu.Lock()

		for k, e := range c.data {
			if time.Since(e.createdAt) > c.interval {
				delete(c.data, k)
			}
		}
		c.mu.Unlock()
	}

}

func NewCache(interval time.Duration) Cache {
	c := Cache{data: make(map[string]cacheEntry)}
	c.interval = interval

	go c.reapLoop()
	return c
}
