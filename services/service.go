package services

import (
	"bytes"
	"io/ioutil"

	"github.com/deelawn/BrainPaaswd/readers"
	"github.com/deelawn/BrainPaaswd/storage"
)

const sourceNotExist = "%s does not exist for source: %s\n"

// Service stores data that is shared between services; it is embedded within the services it serves
type Service struct{}

// ReadData will return byte data that corresponds to a data source that the reader points to
func (s Service) ReadData(reader readers.Reader) ([]byte, error) {

	// Initalize a new reader and pass it to ReadAll to read the data from the data source
	result, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	// Trim the null bytes left over in the file read buffer
	result = bytes.Trim(result, "\u0000")

	return result, nil

}

// CacheIsStale returns true if a cache is stale
func (s Service) CacheIsStale(source string, cache storage.Cache, reader readers.Reader) bool {

	// Retrieves the last time the source was modified
	lastModified, err := reader.GetModifiedTime()

	if err != nil {
		return true
	}

	// If the source data was modified since the cache was last updated, the cache is stale
	if lastModified.After(cache.LastUpdated()) {
		return true
	}

	return false
}

// NewService constructs and returns a new instance of Service
func NewService() *Service {

	return &Service{}
}
