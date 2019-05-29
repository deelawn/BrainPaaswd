package groups

import (
	"log"
	"strings"

	"github.com/deelawn/convert"

	"github.com/deelawn/BrainPaaswd/models"
	"github.com/deelawn/BrainPaaswd/services"
	"github.com/deelawn/BrainPaaswd/storage"
)

const numRecordFields = 4

type Service struct {
	services.Service
}

func (s *Service) readFromSource(cache storage.Cache) ([]models.Group, map[int64]models.Group, error) {

	// Read the data from the data source
	data, err := s.ReadData(s.GroupSource())

	if err != nil {
		return nil, nil, err
	}

	results := make(map[interface{}]interface{})
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

		// Check that it is valid
		if len(fields) != numRecordFields {
			log.Printf("group record on line %d is malformed: %s\n", i+1, record)
			continue
		}

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

		results[gid] = newGroup
		groupList = append(groupList, newGroup)
	}

	// If this function is being executed, it means that the cache needs to updated
	if cache != nil {
		cache.SetData(groupList, results)
	}

	for k, v := range results {
		groupMap[k.(int64)] = v.(models.Group)
	}

	return groupList, groupMap, nil
}

func NewService(service services.Service) *Service {

	return &Service{service}
}
