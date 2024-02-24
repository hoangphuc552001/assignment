package util

import (
	"sync"
	"time"
)

type cacheEntry struct {
	Data         interface{}
	ExpiresAt    time.Time
	CacheTimeOut time.Duration
}

var (
	cache      = make(map[string]cacheEntry)
	cacheMutex sync.RWMutex
)

func SetCache(key string, data interface{}, timeout time.Duration) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	cache[key] = cacheEntry{
		Data:         data,
		ExpiresAt:    time.Now().Add(timeout),
		CacheTimeOut: timeout,
	}
}

func GetCache(key string) (interface{}, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	entry, ok := cache[key]
	if !ok || entry.ExpiresAt.Before(time.Now()) {
		return nil, false
	}
	return entry.Data, true
}
