package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/deelawn/BrainPaaswd/services/users"
)

const url = "http://localhost"

/*
	For the two helper functions below:
		- getTestResponse
		- parseResponse

	Each of these panics on an error for two reasons:
		- No actual tests can be performed reliably if an error occurs
		- Removes the need for error checking in each of the test functions
*/

func getTestResponse(method, path string,
	handler func(http.ResponseWriter, *http.Request), params map[string]string) ([]byte, int) {

	req := httptest.NewRequest(method, fmt.Sprintf("%s/%s", url, path), nil)

	if params != nil {
		req = mux.SetURLVars(req, params)
	}

	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(fmt.Sprintf("Error retrieving data: %s\n", err.Error()))
	}

	return data, resp.StatusCode
}

func parseResponse(data []byte, target interface{}) {

	err := json.Unmarshal(data, target)

	if err != nil {
		panic(fmt.Sprintf("Error unmarshaling response: %s\n", err.Error()))
	}
}

func initUserService() *users.Service {

	return users.NewService(*baseService)
}

/****************************
*
* Tests begin here
*
****************************/

func TestListUsers(t *testing.T) {

	testUserService := initUserService()

	userList := make([]users.User, 0)
	data, status := getTestResponse(http.MethodGet, "users", testUserService.List, nil)
	parseResponse(data, &userList)

	assert.EqualValues(t, http.StatusOK, status)
	assert.Len(t, userList, 14)
}

func TestReadUser(t *testing.T) {

	testUserService := initUserService()

	var user users.User
	data, status := getTestResponse(http.MethodGet, "users/6", testUserService.Read, map[string]string{"uid": "6"})
	parseResponse(data, &user)

	assert.EqualValues(t, http.StatusOK, status)
	assert.EqualValues(t, "nuucp", user.Name)
	assert.EqualValues(t, 6, user.UID)
	assert.EqualValues(t, 5, user.GID)
	assert.EqualValues(t, "uucp login user", user.Comment)
	assert.EqualValues(t, "/var/spool/uucppublic", user.Home)
	assert.EqualValues(t, "/usr/sbin/uucp/uucico", user.Shell)
}

func TestReadUserNonexistent(t *testing.T) {

	testUserService := initUserService()

	data, status := getTestResponse(http.MethodGet, "users/75", testUserService.Read, map[string]string{"uid": "75"})

	assert.EqualValues(t, http.StatusNotFound, status)
	assert.Empty(t, data)
}
