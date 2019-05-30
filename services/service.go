package services

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/deelawn/BrainPaaswd/readers"
	"github.com/deelawn/BrainPaaswd/storage"
)

const sourceNotExist = "%s does not exist for source: %s\n"

// Service defines functionality shared between services
type Service struct {
	source        string
	cache         storage.Cache
	readerBuilder readers.ReaderBuilder
}

func (s *Service) Source() string {

	return s.source
}

func (s *Service) Cache() storage.Cache {

	return s.cache
}

func (s *Service) Reader() readers.Reader {

	return s.readerBuilder(s.source)
}

// ReadFromSource will read resource data from the source the reader points to and return it structured
func (s *Service) ReadFromSource(reader readers.Reader,
	parser ResourceParser) (interface{}, map[interface{}]interface{}, error) {

	// Read the data from the data source
	data, err := s.readData(reader)

	if err != nil {
		return nil, nil, err
	}

	// Define the map and list to return
	resourceMap := make(map[interface{}]interface{})
	resourceList := make([]interface{}, 0)

	// Now do the transformation
	records := strings.Split(string(data), "\n")
	for i, record := range records {
		newResource, id, err := parser(record)

		// Continue past any records that could not be parsed an log the error
		if err != nil {
			log.Printf("%s parsing error on line %d: %v; data: %s\n", s.Source(), i, err, record)
			continue
		}

		resourceMap[id] = newResource
		resourceList = append(resourceList, newResource)
	}

	// If this function is being executed, it means that the cache needs to updated
	s.cache.SetData(resourceList, resourceMap)

	return resourceList, resourceMap, nil
}

// readData will return byte data that corresponds to a data source that the reader points to
func (s *Service) readData(reader readers.Reader) ([]byte, error) {

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
func (s *Service) CacheIsStale(reader readers.Reader) bool {

	// Retrieves the last time the source was modified
	lastModified, err := reader.GetModifiedTime()

	if err != nil {
		return true
	}

	// If the source data was modified since the cache was last updated, the cache is stale
	if lastModified.After(s.cache.LastUpdated()) {
		return true
	}

	return false
}

// GetResources retrieves a list of all resources for a service
func (s *Service) GetResources(parser ResourceParser) (interface{}, error, int) {

	reader := s.Reader()
	resources := s.getCachedResources(reader)

	// The resources were found; return them
	if resources != nil {
		return resources, nil, http.StatusOK
	}

	// If it made it this far, then cached data can't be used
	resources, _, err := s.ReadFromSource(reader, parser)

	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return resources, nil, http.StatusOK
}

// getCachedResources attempts to retrieve a list of all resources for a service from the cache
func (s *Service) getCachedResources(reader readers.Reader) interface{} {

	if !s.CacheIsStale(reader) {
		resources, cacheErr := s.cache.Data()

		if cacheErr == nil {
			return resources
		}
	}

	return nil
}

// GetResource retrieves a resource for a service
func (s *Service) GetResource(id int64, parser ResourceParser) (interface{}, error, int) {

	reader := s.Reader()
	resource := s.getCachedResource(id, reader)

	// The resource was found; return it
	if resource != nil {
		return resource, nil, http.StatusOK
	}

	// If it made it this far, then cached data can't be used
	_, indexedResources, err := s.ReadFromSource(reader, parser)

	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	// Use the map that was returned to check if the resource exists
	var exists bool
	if resource, exists = indexedResources[id]; !exists {
		return nil, fmt.Errorf("%s not found\n", s.source), http.StatusNotFound
	}

	return resource, nil, http.StatusOK
}

// getCachedResource attempts to retrieve a resource for a service from the cache
func (s *Service) getCachedResource(id int64, reader readers.Reader) interface{} {

	if !s.CacheIsStale(reader) {
		resource, cacheErr := s.cache.IndexedData(id)

		if cacheErr == nil {
			return resource
		}
	}

	return nil
}

// NewService constructs and returns a new instance of Service
func NewService(source string, cache storage.Cache, readerBuilder readers.ReaderBuilder) *Service {

	return &Service{
		source:        source,
		cache:         cache,
		readerBuilder: readerBuilder,
	}
}
