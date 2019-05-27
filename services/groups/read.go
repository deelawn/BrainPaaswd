package groups

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

	gid, err := convert.StringToInt64(mux.Vars(r)["gid"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"malformed gid"}`))
		return
	}

	foundGroup, err, statusCode := s.getGroup(gid)

	if err != nil {
		w.WriteHeader(statusCode)
		w.Write([]byte(`{"error":"could not read group"}`))
		return
	}

	respData, err := json.Marshal(&foundGroup)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"group data error"}`))
		return
	}

	_, err = w.Write(respData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"unknown group error"}`))
	}
}

func (s *Service) getGroup(gid int64) (Group, error, int) {

	var group Group

	cache, err := s.GetCache(s.GroupSource())

	if err == nil {
		data, cacheErr := cache.IndexedData(gid)

		if cacheErr == nil {
			var ok bool
			if group, ok = data.(Group); ok {
				return group, nil, http.StatusOK
			}
		}
	}

	// If it made it this far, then cached data can't be used
	indexedGroups, err := s.readFromSource(cache)

	if err != nil {
		return group, err, http.StatusInternalServerError
	}

	var exists bool
	if group, exists = indexedGroups[gid].(Group); !exists {
		return group, errors.New("user not found"), http.StatusNotFound
	}

	return group, nil, http.StatusOK
}
