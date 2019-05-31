package groups

import (
	"encoding/json"
	"net/http"

	"github.com/deelawn/convert"
	"github.com/gorilla/mux"
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
	group, err, statusCode := s.GetResource(gid, ResourceParser)

	if err != nil {
		w.WriteHeader(statusCode)
		w.Write([]byte(`{"error":"could not read group"}`))
		return
	}

	// Convert the group to byte data
	respData, err := json.Marshal(group)

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
