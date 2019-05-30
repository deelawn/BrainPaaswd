package groups

import (
	"errors"
	"fmt"
	"strings"

	"github.com/deelawn/convert"

	"github.com/deelawn/BrainPaaswd/models"
	"github.com/deelawn/BrainPaaswd/services"
)

const numRecordFields = 4

// Service defines a group service
type Service struct {
	*services.Service
}

var resourceParser = func(data string) (interface{}, int64, error) {

	// Skip blank lines
	if len(strings.TrimSpace(data)) == 0 {
		return nil, -1, errors.New("empty record")
	}

	fields := strings.Split(data, ":")

	// Validate the number of fields in the record
	if len(fields) != numRecordFields {
		return nil, -1, errors.New("group record is malformed")
	}

	// GID needs to be an integer
	gid, gidErr := convert.StringToInt64(fields[2])
	if gidErr != nil {
		return nil, -1, fmt.Errorf("invalid gid %s: %v\n", fields[2], gidErr)
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

	return newGroup, gid, nil
}

// NewService returns a new instance of the groups service
func NewService(service *services.Service) *Service {

	return &Service{
		Service: service,
	}
}
