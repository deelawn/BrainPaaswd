package groups

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/deelawn/convert"

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

	// Retrieve any query parameters
	params := r.URL.Query()

	// If the query URI is used and query parameters exist, then apply the filters
	if len(params) > 0 && strings.Index(r.RequestURI, "/groups/query?") == 0 {
		groupList, err = s.filter(groupList, params)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			return
		}
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

func (s *Service) filter(groupList []models.Group, params url.Values) ([]models.Group, error) {

	filteredList := make([]models.Group, 0)

	// First check for gid; if it is provided but not found then the rest of the parameters don't need checking
	if gidStr, exists := params["gid"]; exists {
		gid, err := convert.StringToInt64(gidStr[0])

		if err == nil {
			group, err, statusCode := s.getGroup(gid)

			// Return empty list if no matching gid is found
			if statusCode == http.StatusNotFound {
				return filteredList, nil
			} else if err != nil {
				return nil, err
			}

			filteredList = append(filteredList, group)
		}
	} else {
		// Otherwise add all groups to the filtered list and remove them as necessary
		filteredList = append(filteredList, groupList...)
	}

	// Now check for the rest of the query parameters
	if name, exists := params["name"]; exists {
		for i := 0; i < len(filteredList); i++ {
			if filteredList[i].Name != name[0] {
				filteredList = append(filteredList[:i], filteredList[i+1:]...)
				i--
			}
		}
	}

	if member, exists := params["member"]; exists {
		for _, m := range member {
			for i := 0; i < len(filteredList); i++ {
				if !filteredList[i].ContainsMember(m) {
					filteredList = append(filteredList[:i], filteredList[i+1:]...)
					i--
				}
			}
		}
	}

	return filteredList, nil
}
