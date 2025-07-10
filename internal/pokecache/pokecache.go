package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache    map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	c.cache[key] = cacheEntry{createdAt: time.Now(), val: value}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	val, ok := c.cache[key]
	c.mu.Unlock()
	if ok {
		return val.val, ok
	}
	return nil, false
}

func (c *Cache) reapLoop() {
	interval := c.interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for k, v := range c.cache {
			if time.Since(v.createdAt) > interval {
				delete(c.cache, k)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		interval: interval,
		mu:       sync.Mutex{},
		cache:    make(map[string]cacheEntry),
	}

	go cache.reapLoop()

	return &cache
}
