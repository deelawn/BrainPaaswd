package groups

import (
	"encoding/json"
	"net/http"

	"github.com/deelawn/BrainPaaswd/models"
)

func (s *Service) List(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	groupList, err, statusCode := s.getGroups()

	if err != nil {
		w.WriteHeader(statusCode)
		w.Write([]byte(`{"error":"could not read groups"}`))
		return
	}

	respData, err := json.Marshal(groupList)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"group data error"}`))
	}

	_, err = w.Write(respData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"unknown group error"}`))
	}
}

func (s *Service) getGroups() ([]models.Group, error, int) {

	var groupList []models.Group

	cache, err := s.GetCache(s.GroupSource())

	if err == nil {
		data, cacheErr := cache.Data()

		if cacheErr == nil {
			var ok bool
			if groupList, ok = data.([]models.Group); ok {
				return groupList, nil, http.StatusOK
			}
		}
	}

	// If it made it this far, then cached data can't be used
	groupList, _, err = s.readFromSource(cache)

	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return groupList, nil, http.StatusOK
}
