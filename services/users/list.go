package users

import (
	"encoding/json"
	"net/http"

	"github.com/deelawn/BrainPaaswd/models"
)

func (s *Service) List(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	userList, err, statusCode := s.getUsers()

	if err != nil {
		w.WriteHeader(statusCode)
		w.Write([]byte(`{"error":"could not read users"}`))
		return
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

func (s *Service) getUsers() ([]models.User, error, int) {

	var userList []models.User

	cache, err := s.GetCache(s.UserSource())

	if err == nil {
		data, cacheErr := cache.Data()

		if cacheErr == nil {
			var ok bool
			if userList, ok = data.([]models.User); ok {
				return userList, nil, http.StatusOK
			}
		}
	}

	// If it made it this far, then cached data can't be used
	indexedUsers, err := s.readFromSource(cache)

	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	userList = make([]models.User, len(indexedUsers))
	idx := 0
	for _, user := range indexedUsers {
		userList[idx] = user.(models.User)
		idx++
	}

	return userList, nil, http.StatusOK
}
