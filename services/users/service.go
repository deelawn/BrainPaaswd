package users

import (
	"log"
	"strings"

	"github.com/deelawn/convert"

	"github.com/deelawn/BrainPaaswd/models"
	"github.com/deelawn/BrainPaaswd/readers"
	"github.com/deelawn/BrainPaaswd/services"
	"github.com/deelawn/BrainPaaswd/storage"
)

const numRecordFields = 7

// Service defines a user service
type Service struct {
	services.Service
	source        string
	cache         storage.Cache
	readerBuilder func(source string) readers.Reader
}

// readFromSource will read data from the source that reader points to, update the cache, and return the data
func (s *Service) readFromSource(cache storage.Cache,
	reader readers.Reader) ([]models.User, map[int64]models.User, error) {

	// Read the data from the data source
	data, err := s.ReadData(reader)

	if err != nil {
		return nil, nil, err
	}

	// The cache stores generic mapped data of any types
	results := make(map[interface{}]interface{})

	// Map and list are defined so that the calling routine has the option to use all data from the source or to
	// quickly access a single element
	userMap := make(map[int64]models.User)
	userList := make([]models.User, 0)

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
			log.Printf("user record on line %d is malformed: %s\n", i+1, record)
			continue
		}

		// UID needs to be an integer
		uid, uidErr := convert.StringToInt64(fields[2])
		if uidErr != nil {
			log.Printf("invalid uid %s on line %d\n", fields[2], i+1)
			continue
		}

		// GID needs to be an integer
		gid, gidErr := convert.StringToInt64(fields[3])
		if gidErr != nil {
			log.Printf("invalid gid %s on line %d\n", fields[3], i+1)
			continue
		}

		newUser := models.User{
			Name:    fields[0],
			UID:     uid,
			GID:     gid,
			Comment: fields[4],
			Home:    fields[5],
			Shell:   fields[6],
		}

		// Add to the generic cache map and the result list
		results[uid] = newUser
		userList = append(userList, newUser)
	}

	// If this function is being executed, it means that the cache needs to updated
	if cache != nil {
		cache.SetData(userList, results)
	}

	// Convert from generic map to the map the caller is expecting
	for k, v := range results {
		userMap[k.(int64)] = v.(models.User)
	}

	return userList, userMap, nil
}

// NewService returns a new instance of the users service
func NewService(service services.Service, source string,
	cache storage.Cache, readerBuilder func(source string) readers.Reader) *Service {

	return &Service{
		Service:       service,
		source:        source,
		cache:         cache,
		readerBuilder: readerBuilder,
	}
}
