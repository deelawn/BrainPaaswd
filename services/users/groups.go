package users

import (
	"io/ioutil"
	"net/http"

	"github.com/deelawn/convert"
	"github.com/gorilla/mux"
)

const groupsQuery = "http://localhost:8000/groups/query?member="

func (s *Service) Groups(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")

	uid, err := convert.StringToInt64(mux.Vars(r)["uid"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"malformed uid"}`))
		return
	}

	foundUser, err, statusCode := s.getUser(uid)

	if err != nil {
		w.WriteHeader(statusCode)
		w.Write([]byte(`{"error":"could not read user"}`))
		return
	}

	query := groupsQuery + foundUser.Name

	// Query the groups service over HTTP to fully separate concerns of services
	resp, err := http.Get(query)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not read user groups"}`))
		return
	}

	groups, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"could not parse user groups"}`))
		return
	}

	w.Write(groups)
}
