package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/deelawn/convert"

	"github.com/deelawn/BrainPaaswd/models"
)

// List will return a list of all users or a list specified by the provided query parameters
func (s *Service) List(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")

	// Retrieve the full list of users
	resources, err, statusCode := s.GetResources(ResourceParser)

	if err != nil {
		w.WriteHeader(statusCode)
		w.Write([]byte(`{"error":"could not read users"}`))
		return
	}

	// Get the list and assert them to the proper type
	resourceList := resources.([]interface{})
	userList := make([]models.User, len(resourceList))
	for i, r := range resourceList {
		userList[i] = r.(models.User)
	}

	// Retrieve any query parameters
	params := r.URL.Query()

	// If the query URI is used and query parameters exist, then apply the filters
	if len(params) > 0 && strings.Index(r.RequestURI, "/users/query?") != -1 {
		userList, err = s.filter(userList, params)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			return
		}
	}

	// Convert the list of users to byte data
	respData, err := json.Marshal(userList)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"user data error"}`))
	}

	// Write the response; no need to write the response code as 200 is the default
	_, err = w.Write(respData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"unknown user error"}`))
	}
}

// filter takes a list of users and query parameters and returns the filtered list
func (s *Service) filter(userList []models.User, params url.Values) ([]models.User, error) {

	filteredList := make([]models.User, 0)

	// First check for uid; if it is provided but not found then the rest of the parameters don't need checking
	if uidStr, exists := params["uid"]; exists {
		uid, err := convert.StringToInt64(uidStr[0])

		if err == nil {
			// Retrieve the user using the provided UID
			resource, err, statusCode := s.GetResource(uid, ResourceParser)

			// Return empty list if no matching uid is found
			if statusCode == http.StatusNotFound {
				return filteredList, nil
			} else if err != nil {
				return nil, err
			}

			user, _ := resource.(models.User)

			// If found by UID, this is the only user in the result list; no more can be added but this could be removed
			filteredList = append(filteredList, user)
		}
	} else {
		// Otherwise add all users to the filtered list and remove them as necessary
		filteredList = append(filteredList, userList...)
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

	if gidStr, exists := params["gid"]; exists {
		if gid, intErr := convert.StringToInt64(gidStr[0]); intErr == nil {
			for i := 0; i < len(filteredList); i++ {
				if filteredList[i].GID != gid {
					filteredList = append(filteredList[:i], filteredList[i+1:]...)
					i--
				}
			}
		}
	}

	if comment, exists := params["comment"]; exists {
		for i := 0; i < len(filteredList); i++ {
			if filteredList[i].Comment != comment[0] {
				filteredList = append(filteredList[:i], filteredList[i+1:]...)
				i--
			}
		}
	}

	if home, exists := params["home"]; exists {
		for i := 0; i < len(filteredList); i++ {
			if filteredList[i].Home != home[0] {
				filteredList = append(filteredList[:i], filteredList[i+1:]...)
				i--
			}
		}
	}

	if shell, exists := params["shell"]; exists {
		for i := 0; i < len(filteredList); i++ {
			if filteredList[i].Shell != shell[0] {
				filteredList = append(filteredList[:i], filteredList[i+1:]...)
				i--
			}
		}
	}

	// The remaining list contains only users that match all provided parameters
	return filteredList, nil
}
