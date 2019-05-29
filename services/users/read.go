package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/deelawn/convert"
	"github.com/gorilla/mux"

	"github.com/deelawn/BrainPaaswd/models"
	"github.com/deelawn/BrainPaaswd/readers"
)

// Read will return a user that matches the provided UID or an error if one is not found
func (s *Service) Read(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")

	uid, err := convert.StringToInt64(mux.Vars(r)["uid"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"malformed uid"}`))
		return
	}

	// Retrieve the user
	foundUser, err, statusCode := s.getUser(uid)

	if err != nil {
		w.WriteHeader(statusCode)
		w.Write([]byte(`{"error":"could not read user"}`))
		return
	}

	// Convert the user to byte data
	respData, err := json.Marshal(&foundUser)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"user data error"}`))
		return
	}

	// Write the response; no need to write the response code as 200 is the default
	_, err = w.Write(respData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"unknown user error"}`))
	}
}

// getUser accepts a UID and returns the User that matches or an error if none match
func (s *Service) getUser(uid int64) (models.User, error, int) {

	reader := s.readerBuilder(s.source)
	user, err := s.getCachedUser(uid, reader)

	// No error so the user was found; return it
	if err == nil {
		return user, err, http.StatusOK
	}

	// If it made it this far, then cached data can't be used
	_, indexedUsers, err := s.readFromSource(s.cache, reader)

	if err != nil {
		return user, err, http.StatusInternalServerError
	}

	// Use the map that was returned to check if the user exists
	var exists bool
	if user, exists = indexedUsers[uid]; !exists {
		return user, errors.New("user not found"), http.StatusNotFound
	}

	return user, nil, http.StatusOK
}

// getCachedUser will check if a user is cached and return it if found
func (s *Service) getCachedUser(uid int64, reader readers.Reader) (models.User, error) {

	var user models.User

	if !s.CacheIsStale(s.source, s.cache, reader) {
		// Cache is okay, so retrieve the user from the cached map
		data, cacheErr := s.cache.IndexedData(uid)

		if cacheErr == nil {
			var ok bool
			if user, ok = data.(models.User); ok {
				// Found it! Return it
				return user, nil
			}
		}
	}

	return models.User{}, errors.New("user not in cache")
}
