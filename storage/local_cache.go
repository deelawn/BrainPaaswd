package storage

import (
	"fmt"
	"sync"
	"time"
)

// LocalCache is a local storage mechanism that is thread safe
type LocalCache struct {
	mtx         *sync.RWMutex
	data        interface{}
	indexedData map[interface{}]interface{}
	lastUpdated time.Time
}

// Data reads data from the cache
func (l *LocalCache) Data() (interface{}, error) {

	// TODO: to use context here to define a timeout
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	return l.data, nil
}

// IndexedData tries reading and returning an indexed item using the provided key
func (l *LocalCache) IndexedData(key interface{}) (interface{}, error) {

	// TODO: to use context here to define a timeout
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	var item interface{}
	var exists bool

	if item, exists = l.indexedData[key]; !exists {
		return nil, fmt.Errorf("cached item could not be found for key: %v\n", key)
	}

	return item, nil
}

// SetData overwrites existing cached data and updates the timestamp
func (l *LocalCache) SetData(data interface{}, indexedData map[interface{}]interface{}) error {

	// TODO: to use context here to define a timeout
	l.mtx.Lock()
	defer l.mtx.Unlock()

	l.data = data
	l.indexedData = indexedData
	l.lastUpdated = time.Now()
	return nil
}

// LastUpdated returns the last updated time of the cache
func (l *LocalCache) LastUpdated() time.Time {

	return l.lastUpdated
}

// Initializes a new cache instance
func NewLocalCache() *LocalCache {

	return &LocalCache{
		mtx: &sync.RWMutex{},
	}
}
