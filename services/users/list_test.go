package users

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gorilla/mux"
// 	"github.com/stretchr/testify/assert"

// 	"github.com/deelawn/BrainPaaswd/models"
// 	"github.com/deelawn/BrainPaaswd/services"
// )

// const url = "http://localhost"

// func getTestResponse(method, path string,
// 	handler func(http.ResponseWriter, *http.Request), params map[string]string) ([]byte, int) {

// 	req := httptest.NewRequest(method, fmt.Sprintf("%s/%s", url, path), nil)

// 	if params != nil {
// 		req = mux.SetURLVars(req, params)
// 	}

// 	w := httptest.NewRecorder()
// 	handler(w, req)

// 	resp := w.Result()
// 	data, err := ioutil.ReadAll(resp.Body)

// 	if err != nil {
// 		panic(fmt.Sprintf("Error retrieving data: %s\n", err.Error()))
// 	}

// 	return data, resp.StatusCode
// }

// func parseResponse(data []byte, target interface{}) {

// 	err := json.Unmarshal(data, target)

// 	if err != nil {
// 		panic(fmt.Sprintf("Error unmarshaling response: %s\n", err.Error()))
// 	}
// }

// var tests = []struct {
// 	name string
// 	size int
// 	status int
// 	index int
// 	username string
// 	uid int64
// 	gid int64
// 	comment string
// 	home string
// 	shell string
// 	err string
// }{
// 	{"full list", 98, http.StatusOK, 10, "_ces", 32, 32, "Certificate Enrollment Service", "/var/empty", "/usr/bin/false", ""},
// }

// func TestListUsers(t *testing.T) {

// 	s := services.NewService()

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			userList := make([]models.User, 0)
// 			data, status := getTestResponse(http.MethodGet, "users", List, nil)

// 			assert.EqualValues(t, tt.status, status)

// 			// Don't bother checking the rest if an error is expected
// 			if tt.err != "" {
// 				assert.EqualValues(t, tt.err, string(data))
// 				return
// 			}

// 			parseResponse(data, &userList)

// 			assert.EqualValues(t, tt.size, len(userList))
// 			assert.EqualValues(t, tt.username, userList[tt.index].Name)
// 			assert.EqualValues(t, tt.uid, userList[tt.index].UID)
// 			assert.EqualValues(t, tt.gid, userList[tt.index].GID)
// 			assert.EqualValues(t, tt.comment, userList[tt.index].Comment)
// 			assert.EqualValues(t, tt.home, userList[tt.index].Home)
// 			assert.EqualValues(t, tt.shell, userList[tt.index].Shell)
// 		})
// 	}
// }