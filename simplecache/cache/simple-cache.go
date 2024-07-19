package cache

import (
	"encoding/gob"
	"os"
	"sync"
	"time"
)

// CacheItem represents a single cache item
type CacheItem struct {
	Value      interface{}
	Expiration int64
}

// Cache represents the cache structure
type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

// NewCache creates a new cache
func NewCache() *Cache {
	cache := &Cache{
		items: make(map[string]CacheItem),
	}
	go cache.cleanup()
	return cache
}

// SetCache adds an item to the cache
func (c *Cache) SetCache(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := CacheItem{
		Value:      value,
		Expiration: time.Now().Add(duration).UnixNano(),
	}

	c.items[key] = item
}

// GetCache retrieves an item from the cache
func (c *Cache) GetCache(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if time.Now().UnixNano() > item.Expiration {
		delete(c.items, key) // Remove expired item
		return nil, false
	}
	return item.Value, true
}

// DeleteCache removes an item from the cache
func (c *Cache) DeleteCache(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// cleanup periodically removes expired items from the cache
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

func (c *Cache) SaveToFile(filename string) error {

	//Thread safe
	c.mu.Lock()
	defer c.mu.Unlock()
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(c.items)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) LoadCache(filename string) error {

	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	decode := gob.NewDecoder(file)
	err = decode.Decode(&c.items)
	if err != nil {
		return err
	}
	return nil
}
