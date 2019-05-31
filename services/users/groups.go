package users

import (
	"io/ioutil"
	"net/http"

	"github.com/deelawn/BrainPaaswd/models"
	"github.com/deelawn/convert"
	"github.com/gorilla/mux"
)

const groupsQuery = "http://localhost:8000/groups/query?member="

// Group returns all groups that a user belongs to
func (s *Service) Groups(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")

	uid, err := convert.StringToInt64(mux.Vars(r)["uid"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"malformed uid"}`))
		return
	}

	// Get the user so that the name can be provided to the groups query
	resource, err, statusCode := s.GetResource(uid, ResourceParser)

	if err != nil {
		w.WriteHeader(statusCode)
		w.Write([]byte(`{"error":"could not read user"}`))
		return
	}

	user, _ := resource.(models.User)
	query := groupsQuery + user.Name

	// Query the groups service over HTTP to fully separate concerns of services
	resp, err := http.Get(query)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not read user groups"}`))
		return
	}

	defer resp.Body.Close()
	groups, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not parse user groups"}`))
		return
	}

	// groups is already []byte type, so no need to Unmarshal and Marshal the response from the groups service
	w.Write(groups)
}
