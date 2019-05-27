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

	"github.com/deelawn/BrainPaaswd/services/groups"
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

func initGroupService() *groups.Service {

	return groups.NewService(*baseService)
}

/****************************
*
* Tests begin here
*
****************************/

// Begin user tests

func TestListUsers(t *testing.T) {

	testUserService := initUserService()

	userList := make([]users.User, 0)
	data, status := getTestResponse(http.MethodGet, "users", testUserService.List, nil)
	parseResponse(data, &userList)

	assert.EqualValues(t, http.StatusOK, status)
	assert.Len(t, userList, 98)
}

func TestReadUser(t *testing.T) {

	testUserService := initUserService()

	var user users.User
	data, status := getTestResponse(http.MethodGet, "users/1", testUserService.Read, map[string]string{"uid": "1"})
	parseResponse(data, &user)

	assert.EqualValues(t, http.StatusOK, status)
	assert.EqualValues(t, "daemon", user.Name)
	assert.EqualValues(t, 1, user.UID)
	assert.EqualValues(t, 1, user.GID)
	assert.EqualValues(t, "System Services", user.Comment)
	assert.EqualValues(t, "/var/root", user.Home)
	assert.EqualValues(t, "/usr/bin/false", user.Shell)
}

func TestReadUserNonexistent(t *testing.T) {

	testUserService := initUserService()

	data, status := getTestResponse(http.MethodGet, "users/1000", testUserService.Read, map[string]string{"uid": "1000"})

	assert.EqualValues(t, http.StatusNotFound, status)
	assert.EqualValues(t, `{"error":"could not read user"}`, string(data))
}

// End user tests

// Begin group tests
func TestListGroups(t *testing.T) {

	testGroupService := initGroupService()

	groupList := make([]groups.Group, 0)
	data, status := getTestResponse(http.MethodGet, "groups", testGroupService.List, nil)
	parseResponse(data, &groupList)

	assert.EqualValues(t, http.StatusOK, status)
	assert.Len(t, groupList, 125)
}

func TestReadGroup(t *testing.T) {

	testGroupService := initGroupService()

	var group groups.Group
	data, status := getTestResponse(http.MethodGet, "groups/29", testGroupService.Read, map[string]string{"gid": "29"})
	parseResponse(data, &group)

	assert.EqualValues(t, http.StatusOK, status)
	assert.EqualValues(t, "certusers", group.Name)
	assert.EqualValues(t, 29, group.GID)
	assert.Len(t, group.Members, 6)
	assert.EqualValues(t, "root", group.Members[0])
}

func TestReadGroupNonExistent(t *testing.T) {

	testGroupService := initGroupService()
	data, status := getTestResponse(http.MethodGet, "groups/29", testGroupService.Read, map[string]string{"gid": "9999"})
	
	assert.EqualValues(t, http.StatusNotFound, status)
	assert.EqualValues(t, `{"error":"could not read group"}`, string(data))
}

// End group tests