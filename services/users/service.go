package users

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/deelawn/BrainPaaswd/services"
	"github.com/deelawn/convert"
)

const numRecordFields = 7

type Service struct {
	services.Service
}

func (s *Service) readUsersFromFile() ([]User, error) {

	// Read the data from the file
	data, err := ioutil.ReadFile(s.PasswdPath)

	if err != nil {
		return nil, err
	}

	results := make([]User, 0)

	// Now do the transformation
	records := strings.Split(string(data), "\n")
	for i, record := range records {
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

		gid, gidErr := convert.StringToInt64(fields[2])
		if gidErr != nil {
			log.Printf("invalid gid %s on line %d\n", fields[3], i+1)
			continue
		}

		newUser := User{
			Name:    fields[0],
			UID:     uid,
			GID:     gid,
			Comment: fields[4],
			Home:    fields[5],
			Shell:   fields[6],
		}

		results = append(results, newUser)
	}

	return results, nil
}

func (s *Service) ListUsers(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	results, err := s.readUsersFromFile()

	if err != nil {
		w.Write([]byte(`{"error":"could not read users"}`))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respData, err := json.Marshal(results)

	if err != nil {
		w.Write([]byte(`{"error":"user data error"}`))
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err = w.Write(respData)

	if err != nil {
		w.Write([]byte(`{"error":"unknown user error"}`))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func NewService(service services.Service) *Service {

	return &Service{service}
}
