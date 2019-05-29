package groups

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/deelawn/convert"
	"github.com/gorilla/mux"

	"github.com/deelawn/BrainPaaswd/models"
	"github.com/deelawn/BrainPaaswd/readers"
)

// Read will return a group that matches the provided GID or an error if one is not found
func (s *Service) Read(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")

	gid, err := convert.StringToInt64(mux.Vars(r)["gid"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"malformed gid"}`))
		return
	}

	// Retrieve the group
	foundGroup, err, statusCode := s.getGroup(gid)

	if err != nil {
		w.WriteHeader(statusCode)
		w.Write([]byte(`{"error":"could not read group"}`))
		return
	}

	// Convert the group to byte data
	respData, err := json.Marshal(&foundGroup)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"group data error"}`))
		return
	}

	// Write the response; no need to write the response code as 200 is the default
	_, err = w.Write(respData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"unknown group error"}`))
	}
}

// getGroup accepts a GID and returns the Group that matches or an error if none match
func (s *Service) getGroup(gid int64) (models.Group, error, int) {

	reader := s.readerBuilder(s.source)
	group, err := s.getCachedGroup(gid, reader)

	// No error so the group was found; return it
	if err == nil {
		return group, err, http.StatusOK
	}

	// If it made it this far, then cached data can't be used
	_, indexedGroups, err := s.readFromSource(s.cache, reader)

	if err != nil {
		return group, err, http.StatusInternalServerError
	}

	// Use the map that was returned to check if the group exists
	var exists bool
	if group, exists = indexedGroups[gid]; !exists {
		return group, errors.New("group not found"), http.StatusNotFound
	}

	return group, nil, http.StatusOK
}

// getCachedGroup will check if a group is cached and return it if found
func (s *Service) getCachedGroup(gid int64, reader readers.Reader) (models.Group, error) {

	var group models.Group

	if !s.CacheIsStale(s.source, s.cache, reader) {
		// Cache is okay, so retrieve the group from the cached map
		data, cacheErr := s.cache.IndexedData(gid)

		if cacheErr == nil {
			var ok bool
			if group, ok = data.(models.Group); ok {
				// Found it! Return it
				return group, nil
			}
		}
	}

	return models.Group{}, errors.New("group not in cache")
}
