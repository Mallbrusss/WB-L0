package cache

import (
	"sync"
	"time"
)

type Cache struct {
	data  map[string]interface{}
	mutex sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, exists := c.data[key]
	return value, exists
}

func (c *Cache) Set(key string, value interface{}, expiration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
	if expiration > 0 {
		go c.expireKey(key, expiration)
	}
}

func (c *Cache) expireKey(key string, expiration time.Duration) {
	<-time.After(expiration)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.data, key)
}
