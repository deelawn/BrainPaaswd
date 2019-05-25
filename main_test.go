package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

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

func getTestResponse(method, path string, handler func(http.ResponseWriter, *http.Request)) []byte {

	req := httptest.NewRequest(method, fmt.Sprintf("%s/%s", url, path), nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(fmt.Sprintf("Error retrieving data: %s\n", err.Error()))
	}

	return data
}

func parseResponse(data []byte, target interface{}) {

	err := json.Unmarshal(data, target)

	if err != nil {
		panic(fmt.Sprintf("Error unmarshaling response: %s\n", err.Error()))
	}
}

/****************************
*
* Tests begin here
*
****************************/

func TestListUsers(t *testing.T) {

	data := getTestResponse(http.MethodGet, "users", UserService)
	userList := make([]users.User, 0)
	parseResponse(data, &userList)

	assert.Equal(t, "Users response", string(data))
}
