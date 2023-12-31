package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("test data"),
		},
		{
			key: "https://example.com/testpath",
			val: []byte("test data 2"),
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", c), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)

			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find a value for key %v", c.key)
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value %v for key %v, got %v instead", c.val, c.key, val)
				return
			}
		})
	}
}

func TestCleanLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = 15 * time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("test data"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find a value for key https://example.com")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected value for key https://example.com to be removed from cache after %v", waitTime)
		return
	}
}

func TestUpdateKey(t *testing.T) {
	cache := NewCache(5 * time.Second)

	cache.Add("https://example.com", []byte("test data"))
	cache.Add("https://example.com", []byte("test data 2"))

	val, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find a value for key https://example.com")
		return
	}

	if string(val) != "test data 2" {
		t.Errorf("expected to find value test data 2 for key https://example.com, got %v instead", val)
		return
	}
}
