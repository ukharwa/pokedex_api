package internal

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	interval time.Duration
	mutex    sync.Mutex
	cache    map[string]cacheEntry
}

var cache Cache

func NewCache(interval time.Duration) (*Cache, error) {
	c := &Cache{
		interval: interval,
		cache:    make(map[string]cacheEntry),
	}

	go c.reapLoop()
	return c, nil
}

func (c *Cache) Add(key string, val []byte) error {
	c.mutex.Lock()
	c.cache[key] = cacheEntry{time.Now(), val}
	c.mutex.Unlock()
	return nil
}

func (c *Cache) Get(key string) ([]byte, bool, error) {
	c.mutex.Lock()
	val, exists := c.cache[key]
	c.mutex.Unlock()
	return val.val, exists, nil
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Cleaning up cache...")
			var count int
			c.mutex.Lock()
			for key, entry := range c.cache {
				if time.Now().Sub(entry.createdAt) > c.interval {
					delete(c.cache, key)
					count++
				}
			}
			c.mutex.Unlock()
			fmt.Printf("Clean up complete. (%d) outdated entries removed.\n", count)
		}
	}
}
