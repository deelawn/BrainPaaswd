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

// List will return a list of all groups or a list specified by the provided query parameters
func (s *Service) List(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")

	// Retrieve the full list of users
	resources, err, statusCode := s.GetResources(ResourceParser)

	if err != nil {
		w.WriteHeader(statusCode)
		w.Write([]byte(`{"error":"could not read groups"}`))
		return
	}

	// Get the list and assert them to the proper type
	resourceList := resources.([]interface{})
	groupList := make([]models.Group, len(resourceList))
	for i, r := range resourceList {
		groupList[i] = r.(models.Group)
	}

	// Retrieve any query parameters
	params := r.URL.Query()

	// If the query URI is used and query parameters exist, then apply the filters
	if len(params) > 0 && strings.Index(r.RequestURI, "/groups/query?") != -1 {
		groupList, err = s.filter(groupList, params)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			return
		}
	}

	// Convert the list of groups to byte data
	respData, err := json.Marshal(groupList)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"group data error"}`))
	}

	// Write the response; no need to write the response code as 200 is the default
	_, err = w.Write(respData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"unknown group error"}`))
	}
}

// filter takes a list of groups and query parameters and returns the filtered list
func (s *Service) filter(groupList []models.Group, params url.Values) ([]models.Group, error) {

	filteredList := make([]models.Group, 0)

	// First check for gid; if it is provided but not found then the rest of the parameters don't need checking
	if gidStr, exists := params["gid"]; exists {
		gid, err := convert.StringToInt64(gidStr[0])

		if err == nil {
			// Retrieve the group using the provided GID
			resource, err, statusCode := s.GetResource(gid, ResourceParser)

			// Return empty list if no matching gid is found
			if statusCode == http.StatusNotFound {
				return filteredList, nil
			} else if err != nil {
				return nil, err
			}

			group, _ := resource.(models.Group)

			// If found by GID, this is the only group in the result list; no more can be added but this could be removed
			filteredList = append(filteredList, group)
		}
	} else {
		// Otherwise add all groups to the filtered list and remove them as necessary
		filteredList = append(filteredList, groupList...)
	}

	/*
		Now check for the rest of the query parameters. Each of the loops below will filter each parameter
		sequentially and remove items that do not match the provided parameter
	*/

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
				// Use the constructed group map to filter more quickly rather than yet another loop
				if !filteredList[i].ContainsMember(m) {
					filteredList = append(filteredList[:i], filteredList[i+1:]...)
					i--
				}
			}
		}
	}

	return filteredList, nil
}
