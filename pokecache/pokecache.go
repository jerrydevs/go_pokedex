package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c cacheEntry) String() string {
	return fmt.Sprintf("createdAt: %v, val: %v", c.createdAt, string(c.val))
}

type Cache struct {
	interval time.Duration
	entries  map[string]cacheEntry
	lock     sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	if interval == 0 {
		interval = 5 * time.Second
	}

	ticker := time.NewTicker(interval)
	newCache := &Cache{
		interval: interval,
		entries:  make(map[string]cacheEntry),
		lock:     sync.Mutex{},
	}

	go newCache.cleanLoop(ticker)

	return newCache
}

func (c *Cache) Add(key string, val []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) cleanLoop(ticker *time.Ticker) {
	for {
		<-ticker.C

		fmt.Println("Cleaning cache!!!")
		c.lock.Lock()
		for key, entry := range c.entries {
			if entry.createdAt.Add(c.interval).Before(time.Now()) {
				delete(c.entries, key)
			}
		}
		c.lock.Unlock()
	}
}
