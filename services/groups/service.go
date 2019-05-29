package groups

import (
	"log"
	"strings"

	"github.com/deelawn/convert"

	"github.com/deelawn/BrainPaaswd/models"
	"github.com/deelawn/BrainPaaswd/readers"
	"github.com/deelawn/BrainPaaswd/services"
	"github.com/deelawn/BrainPaaswd/storage"
)

const numRecordFields = 4

// Service defines a group service
type Service struct {
	services.Service
	source        string
	cache         storage.Cache
	readerBuilder func(source string) readers.Reader
}

// readFromSource will read data from the source that reader points to, update the cache, and return the data
func (s *Service) readFromSource(cache storage.Cache,
	reader readers.Reader) ([]models.Group, map[int64]models.Group, error) {

	// Read the data from the data source
	data, err := s.ReadData(reader)

	if err != nil {
		return nil, nil, err
	}

	// The cache stores generic mapped data of any types
	results := make(map[interface{}]interface{})

	// Map and list are defined so that the calling routine has the option to use all data from the source or to
	// quickly access a single element
	groupMap := make(map[int64]models.Group)
	groupList := make([]models.Group, 0)

	// Now do the transformation
	records := strings.Split(string(data), "\n")
	for i, record := range records {

		// Skip blank lines
		if len(strings.TrimSpace(record)) == 0 {
			continue
		}

		fields := strings.Split(record, ":")

		/*
		 * I made a decision to just log errors and continue should it come accross anything
		 * that can't be parsed (invalid format). It is also possible that the desired behavior
		 * for this situation be to return an error. An arbitrary choice was made due to no
		 * further requirements or details being provided.
		 */

		// Validate the number of fields in the record
		if len(fields) != numRecordFields {
			log.Printf("group record on line %d is malformed: %s\n", i+1, record)
			continue
		}

		// GID needs to be an integer
		gid, gidErr := convert.StringToInt64(fields[2])
		if gidErr != nil {
			log.Printf("invalid gid %s on line %d: %s\n", fields[2], i+1, gidErr.Error())
			continue
		}

		newGroup := models.Group{
			Name: fields[0],
			GID:  gid,
		}

		// See the comment by members field in group.go for an explanation as to why Members isn't assigned explicity
		// to the value that is being store in the local members value on the line below.
		members := strings.Split(fields[3], ",")
		for _, m := range members {
			newGroup.AddMember(m)
		}

		// Add to the generic cache map and the result list
		results[gid] = newGroup
		groupList = append(groupList, newGroup)
	}

	// If this function is being executed, it means that the cache needs to updated
	if cache != nil {
		cache.SetData(groupList, results)
	}

	// Convert from generic map to the map the caller is expecting
	for k, v := range results {
		groupMap[k.(int64)] = v.(models.Group)
	}

	return groupList, groupMap, nil
}

// NewService returns a new instance of the groups service
func NewService(service services.Service, source string,
	cache storage.Cache, readerBuilder func(source string) readers.Reader) *Service {

	return &Service{
		Service:       service,
		source:        source,
		cache:         cache,
		readerBuilder: readerBuilder,
	}
}
