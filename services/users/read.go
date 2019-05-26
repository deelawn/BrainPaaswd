package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/deelawn/convert"
	"github.com/gorilla/mux"
)

func (s *Service) Read(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")

	// Even though using the regex in the router, still need to error check in case of potential overflow
	uid, err := convert.StringToInt64(mux.Vars(r)["uid"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"malformed uid"}`))
		return
	}

	foundUser, err, statusCode := s.getUser(uid)

	if err != nil {
		w.WriteHeader(statusCode)
		w.Write([]byte(`{"error":"could not read users"}`))
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

func (s *Service) getUser(uid int64) (User, error, int) {

	var user User

	cache, err := s.GetCache(s.UserSource())

	if err == nil {
		data, cacheErr := cache.IndexedData(uid)

		if cacheErr == nil {
			var ok bool
			if user, ok = data.(User); ok {
				return user, nil, http.StatusOK
			}
		}
	}

	// If it made it this far, then cached data can't be used
	indexedUsers, err := s.readFromSource(cache)

	if err != nil {
		return user, err, http.StatusInternalServerError
	}

	var exists bool
	if user, exists = indexedUsers[uid].(User); !exists {
		return user, errors.New("user not found"), http.StatusNotFound
	}

	return user, nil, http.StatusOK
}
