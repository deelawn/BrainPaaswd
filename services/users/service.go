package users

import (
	"errors"
	"fmt"
	"strings"

	"github.com/deelawn/convert"

	"github.com/deelawn/BrainPaaswd/models"
	"github.com/deelawn/BrainPaaswd/services"
)

const numRecordFields = 7

// Service defines a user service
type Service struct {
	*services.Service
}

var ResourceParser = func(data string) (interface{}, int64, error) {

	// Skip blank lines
	if len(strings.TrimSpace(data)) == 0 {
		return nil, -1, errors.New("empty record")
	}

	fields := strings.Split(data, ":")

	// Validate the number of fields in the record
	if len(fields) != numRecordFields {
		return nil, -1, errors.New("user record is malformed")
	}

	// UID needs to be an integer
	uid, uidErr := convert.StringToInt64(fields[2])
	if uidErr != nil {
		return nil, -1, fmt.Errorf("invalid uid %s : %v\n", fields[2], uidErr)
	}

	// GID needs to be an integer
	gid, gidErr := convert.StringToInt64(fields[3])
	if gidErr != nil {
		return nil, -1, fmt.Errorf("invalid gid %s : %v\n", fields[3], gidErr)
	}

	newUser := models.User{
		Name:    fields[0],
		UID:     uid,
		GID:     gid,
		Comment: fields[4],
		Home:    fields[5],
		Shell:   fields[6],
	}

	return newUser, uid, nil
}

// NewService returns a new instance of the users service
func NewService(service *services.Service) *Service {

	return &Service{
		Service: service,
	}
}
