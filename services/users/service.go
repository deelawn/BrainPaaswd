package users

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/deelawn/BrainPaaswd/services"
	"github.com/deelawn/convert"
)

const numRecordFields = 7

type Service struct {
	services.Service
}

func (s *Service) readUsers() (map[int64]User, error) {

	// Read the data from the data source
	data, err := s.ReadData(s.PasswdPath)

	if err != nil {
		return nil, err
	}

	results := make(map[int64]User)

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

		gid, gidErr := convert.StringToInt64(fields[3])
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

		results[uid] = newUser
	}

	return results, nil
}

func (s *Service) List(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	results, err := s.readUsers()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not read users"}`))
		return
	}

	userList := make([]User, len(results))
	idx := 0
	for _, user := range results {
		userList[idx] = user
		idx++
	}

	respData, err := json.Marshal(userList)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"user data error"}`))
	}

	_, err = w.Write(respData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"unknown user error"}`))
	}
}

func (s *Service) Read(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")

	// Even though using the regex in the router, still need to error check in case of potential overflow
	uid, err := strconv.ParseInt(mux.Vars(r)["uid"], 10, 0)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"malformed uid"}`))
		return
	}

	results, err := s.readUsers()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not read users"}`))
		return
	}

	var foundUser User
	var exists bool

	if foundUser, exists = results[uid]; !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	respData, err := json.Marshal(&foundUser)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"user data error"}`))
		return
	}

	_, err = w.Write(respData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"unknown user error"}`))
	}
}

func NewService(service services.Service) *Service {

	return &Service{service}
}
