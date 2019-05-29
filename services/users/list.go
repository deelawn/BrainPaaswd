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

// may be needed if I write a query parser
// var validParams = map[string]interface{} {
// 	"name": nil,
// 	"uid": nil,
// 	"gid": nil,
// 	"comment": nil,
// 	"home": nil,
// 	"shell": nil,
// }

func (s *Service) List(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	userList, err, statusCode := s.getUsers()

	if err != nil {
		w.WriteHeader(statusCode)
		w.Write([]byte(`{"error":"could not read users"}`))
		return
	}

	// Retrieve any query parameters
	params := r.URL.Query()

	// If the query URI is used and query parameters exist, then apply the filters
	if len(params) > 0 && strings.Index(r.RequestURI, "/users/query?") == 0 {
		userList, err = s.filter(userList, params)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			return
		}
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
	userList, _, err = s.readFromSource(cache)

	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return userList, nil, http.StatusOK
}

func (s *Service) filter(userList []models.User, params url.Values) ([]models.User, error) {

	filteredList := make([]models.User, 0)

	// First check for uid; if it is provided but not found then the rest of the parameters don't need checking
	if uidStr, exists := params["uid"]; exists {
		uid, err := convert.StringToInt64(uidStr[0])

		if err == nil {
			user, err, statusCode := s.getUser(uid)

			// Return empty list if no matching uid is found
			if statusCode == http.StatusNotFound {
				return filteredList, nil
			} else if err != nil {
				return nil, err
			}

			filteredList = append(filteredList, user)
		}
	} else {
		// Otherwise add all users to the filtered list and remove them as necessary
		filteredList = append(filteredList, userList...)
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

	return filteredList, nil
}
