package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	lock  sync.RWMutex
}

func (c *Cache) Add(key string, val []byte) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, exists := c.cache[key]; exists {
		return false
	}

	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	return true
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval * time.Second)
	for range ticker.C {
		c.lock.Lock()
		for key, entry := range c.cache {
			if time.Since(entry.createdAt) > interval {
				delete(c.cache, key)
			}
		}
		c.lock.Unlock()
	}
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache: make(map[string]cacheEntry),
		lock:  sync.RWMutex{},
	}
	go c.reapLoop(interval)
	return c
}
