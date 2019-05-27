package users

import (
	"log"
	"strings"

	"github.com/deelawn/convert"

	"github.com/deelawn/BrainPaaswd/models"
	"github.com/deelawn/BrainPaaswd/services"
	"github.com/deelawn/BrainPaaswd/storage"
)

const numRecordFields = 7

type Service struct {
	services.Service
}

func (s *Service) readFromSource(cache storage.Cache) ([]models.User, map[int64]models.User, error) {

	// Read the data from the data source
	data, err := s.ReadData(s.UserSource())

	if err != nil {
		return nil, nil, err
	}

	results := make(map[interface{}]interface{})
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

		// Check that it is valid
		if len(fields) != numRecordFields {
			log.Printf("user record on line %d is malformed: %s\n", i+1, record)
			continue
		}

		uid, uidErr := convert.StringToInt64(fields[2])
		if uidErr != nil {
			log.Printf("invalid uid %s on line %d\n", fields[2], i+1)
			continue
		}

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

		results[uid] = newUser
		userList = append(userList, newUser)
	}

	// If this function is being executed, it means that the cache needs to updated
	if cache != nil {
		cache.SetData(userList, results)
	}

	for k, v := range results {
		userMap[k.(int64)] = v.(models.User)
	}

	return userList, userMap, nil
}

func NewService(service services.Service) *Service {

	return &Service{service}
}
