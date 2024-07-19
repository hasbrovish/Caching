package cache

import (
	"sync"
	"time"
)

// Create cache function set , get , delete
// Need to make caching operation thread safe to prevent from Race Condtion
// Mutex locks are use to prevent Race Condition
// Need to make a goroutine or periodic function which will run in every minute and remove all the expired cache
type CacheItem struct {
	Value      interface{}
	Expiration int64
}

type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

func init() {
	cache := NewCache()
	done := make(chan bool)

	go func() {
		for i := 0; i < 1000; i++ {
			cache.setCache("Key1", "Value1", 10*time.Second)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			cache.getCache("Key1")
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			cache.deleteCache("Key1")
		}
		done <- true
	}()
	<-done
	<-done
	<-done
}

func (c *Cache) cleanup() {

	for {
		time.Sleep(time.Minute)
		c.mu.Lock()
		for key, item := range c.items {
			if time.Now().UnixNano() > item.Expiration {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}

}

func NewCache() *Cache {

	cache := &Cache{
		items: make(map[string]CacheItem),
	}
	go cache.cleanup()
	return cache
}

func (c *Cache) setCache(key string, value interface{}, duration time.Duration) {

	c.mu.Lock()
	defer c.mu.Unlock()

	item := CacheItem{
		Value:      value,
		Expiration: time.Now().Add(duration).UnixNano(),
	}

	c.items[key] = item
}

func (c *Cache) getCache(key string) (interface{}, bool) {

	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if time.Now().UnixNano() > c.items[key].Expiration {
		return nil, false
	}
	return item.Value, true
}

func (c *Cache) deleteCache(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}
