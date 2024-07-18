package main

// import (
// 	"time"
// )

// type CacheItem struct {
// 	Value      interface{}
// 	Expiration int64
// }

// type Cache struct {
// 	items map[string]CacheItem
// }

// func NewCache() *Cache {
// 	cache := &Cache{
// 		items: make(map[string]CacheItem),
// 	}
// 	go cache.cleanup()
// 	return cache
// }

// func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
// 	c.items[key] = CacheItem{
// 		Value:      value,
// 		Expiration: time.Now().Add(duration).UnixNano(),
// 	}
// }

// func (c *Cache) Get(key string) (interface{}, bool) {
// 	item, found := c.items[key]
// 	if !found {
// 		return nil, false
// 	}
// 	if time.Now().UnixNano() > item.Expiration {
// 		return nil, false
// 	}
// 	return item.Value, true
// }

// func (c *Cache) Delete(key string) {
// 	delete(c.items, key)
// }

// func (c *Cache) cleanup() {
// 	for {
// 		time.Sleep(time.Minute)
// 		for key, item := range c.items {
// 			if time.Now().UnixNano() > item.Expiration {
// 				delete(c.items, key)
// 			}
// 		}
// 	}
// }

// func main() {
// 	cache := NewCache()
// 	done := make(chan bool)

// 	// Goroutine 1: Set key1
// 	go func() {
// 		for i := 0; i < 1000; i++ {
// 			cache.Set("key1", "value1", 10*time.Second)
// 		}
// 		done <- true
// 	}()

// 	// Goroutine 2: Get key1
// 	go func() {
// 		for i := 0; i < 1000; i++ {
// 			cache.Get("key1")
// 		}
// 		done <- true
// 	}()

// 	// Goroutine 3: Delete key1
// 	go func() {
// 		for i := 0; i < 1000; i++ {
// 			cache.Delete("key1")
// 		}
// 		done <- true
// 	}()

// 	<-done
// 	<-done
// 	<-done
// }
