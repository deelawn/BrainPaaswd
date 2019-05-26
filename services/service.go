package services

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/deelawn/BrainPaaswd/readers"
	"github.com/deelawn/BrainPaaswd/storage"
)

const sourceNotExist = "%s does not exist for source: %s\n"

// Service stores data that is shared between services; it is embedded within the services it serves
type Service struct {
	userSource     string
	groupSource    string
	readerBuilders map[string]func(source string) readers.Resource
	caches         map[string]storage.Cache
}

func (s Service) getReader(source string) (readers.Resource, error) {

	var readerBuilder func(source string) readers.Resource
	var exists bool

	if readerBuilder, exists = s.readerBuilders[source]; !exists {
		return nil, fmt.Errorf(sourceNotExist, "reader", source)
	}

	return readerBuilder(source), nil
}

// ReadData will return byte data that corresponds to a data source represented by a string value
func (s Service) ReadData(source string) ([]byte, error) {

	reader, err := s.getReader(source)

	if err != nil {
		return nil, err
	}

	// Initalize a new reader and pass it to ReadAll to read the data from the data source
	result, err := ioutil.ReadAll(reader)

	// Trim the null bytes left over in the file read buffer
	result = bytes.Trim(result, "\u0000")

	return result, err

}

// UserSource returns the user source string; not editable
func (s Service) UserSource() string {

	return s.userSource
}

// GroupSource returns the group source string; not editable
func (s Service) GroupSource() string {

	return s.groupSource
}

// GetCache returns cached data from the specified source
func (s Service) GetCache(source string) (storage.Cache, error) {

	var cache storage.Cache
	var exists bool

	// First check if the cache exists
	if cache, exists = s.caches[source]; !exists {
		return cache, fmt.Errorf(sourceNotExist, "cache", source)
	}

	// Okay it exists; now check if it can still be used
	reader, err := s.getReader(source)

	if err != nil {
		return cache, err
	}

	lastModified, err := reader.GetModifiedTime()

	if err != nil {
		return cache, err
	}

	if lastModified.After(cache.LastUpdated()) {
		return cache, errors.New("Cache is stale\n")
	}

	return cache, nil
}

// NewService constructs and returns a new instance of Service
func NewService(userSource, groupSource string,
	readerBuilders map[string]func(source string) readers.Resource) *Service {

	return &Service{
		userSource:     userSource,
		groupSource:    groupSource,
		readerBuilders: readerBuilders,
		caches: map[string]storage.Cache{
			userSource:  storage.NewLocalCache(),
			groupSource: storage.NewLocalCache(),
		},
	}
}
