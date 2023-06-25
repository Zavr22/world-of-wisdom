package cache

import (
	clock "github.com/Zavr22/world-of-wisdom/internal/pkg/clock"
	"sync"
)

// InMemoryCache is an in-memory implementation of the Cache interface.
type InMemoryCache struct {
	dataMap map[int]inMemoryValue // A map of integer keys to in-memory cache values
	lock    *sync.Mutex           // A mutex to ensure thread-safety
	clock   clock.SystemClock     // A system clock for keeping track of time
}

type inMemoryValue struct {
	SetTime    int64 // The time at which the value was set
	Expiration int64 // The expiration time of the value, in seconds
}

// InitInMemoryCache initializes a new instance of InMemoryCache and returns a pointer to it.
func InitInMemoryCache(clock *clock.SystemClock) *InMemoryCache {
	return &InMemoryCache{
		dataMap: make(map[int]inMemoryValue),
		lock:    &sync.Mutex{},
		clock:   *clock,
	}
}

// AddCash adds a new cache entry to the InMemoryCache.
func (c *InMemoryCache) AddCash(key int, expiration int64) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.dataMap[key] = inMemoryValue{
		SetTime:    c.clock.Now().Unix(),
		Expiration: expiration,
	}
	return nil
}

// GetCash retrieves a cache entry from the InMemoryCache.
func (c *InMemoryCache) GetCash(key int) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	value, ok := c.dataMap[key]
	if ok && c.clock.Now().Unix()-value.SetTime > value.Expiration {
		// If the value has expired, return false and do not return the value.
		return false, nil
	}
	return ok, nil
}

// DeleteCash deletes a cache entry from the InMemoryCache.
func (c *InMemoryCache) DeleteCash(key int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.dataMap, key)
}
