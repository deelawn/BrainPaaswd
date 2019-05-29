package storage

import (
	"time"
)

// Cache defines the interface any cache like mechanism
type Cache interface {
	Data() (interface{}, error)
	IndexedData(key interface{}) (interface{}, error)
	SetData(data interface{}, indexedData map[interface{}]interface{}) error
	LastUpdated() time.Time
}
