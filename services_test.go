package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/deelawn/BrainPaaswd/models"
	"github.com/deelawn/BrainPaaswd/readers/file"
	"github.com/deelawn/BrainPaaswd/services"
	"github.com/deelawn/BrainPaaswd/services/groups"
	"github.com/deelawn/BrainPaaswd/services/users"
	"github.com/deelawn/BrainPaaswd/storage"
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

// Initialization functions are handy

func initUserService() *users.Service {

	return users.NewService(services.NewService("passwd", storage.NewLocalCache(), file.NewReader, users.ResourceParser))
}

func initGroupService() *groups.Service {

	return groups.NewService(services.NewService("group", storage.NewLocalCache(), file.NewReader, groups.ResourceParser))
}

/****************************
*
* Tests begin here
*
****************************/

// Begin user tests

func TestListUsers(t *testing.T) {

	tests := []struct {
		name     string
		size     int
		status   int
		index    int
		username string
		uid      int64
		gid      int64
		comment  string
		home     string
		shell    string
		path     string
	}{
		{"full list one", 98, http.StatusOK, 10, "_ces", 32, 32, "Certificate Enrollment Service", "/var/empty", "/usr/bin/false", "users"},
		{"full list two", 98, http.StatusOK, 1, "root", 0, 0, "System Administrator", "/var/root", "/bin/sh", "users"},
		{"found by username", 1, http.StatusOK, 0, "_jabber", 84, 84, "Jabber XMPP Server", "/var/empty", "/usr/bin/false", "users/query?name=_jabber"},
		{"found by uid", 1, http.StatusOK, 0, "_netbios", 222, 222, "NetBIOS", "/var/empty", "/usr/bin/false", "users/query?uid=222"},
		{"found by gid", 1, http.StatusOK, 0, "_mysql", 74, 74, "MySQL Server", "/var/empty", "/usr/bin/false", "users/query?gid=74"},
		{"found by comment", 1, http.StatusOK, 0, "_unknown", 99, 99, "Unknown User", "/var/empty", "/usr/bin/false", "users/query?comment=Unknown%20User"},
		{"found by home", 1, http.StatusOK, 0, "_lp", 26, 26, "Printing Services", "/var/spool/cups", "/usr/bin/false", "users/query?home=%2Fvar%2Fspool%2Fcups"},
		{"found by shell", 95, http.StatusOK, 15, "_sandbox", 60, 60, "Seatbelt", "/var/empty", "/usr/bin/false", "users/query?shell=%2Fusr%2Fbin%2Ffalse"},
		{"found by home and shell", 68, http.StatusOK, 6, "_mcxalr", 54, 54, "MCX AppLaunch", "/var/empty", "/usr/bin/false", "users/query?home=%2Fvar%2Fempty&shell=%2Fusr%2Fbin%2Ffalse"},
		{"found by gid and comment", 1, http.StatusOK, 0, "_appserver", 79, 79, "Application Server", "/var/empty", "/usr/bin/false", "users/query?gid=79&comment=Application%20Server"},
		{"not found by uid and gid", 0, http.StatusOK, 0, "", 0, 0, "", "", "", "users/query?uid=211&gid=212"},
	}

	s := initUserService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userList := make([]models.User, 0)
			data, status := getTestResponse(http.MethodGet, tt.path, s.List, nil)

			assert.EqualValues(t, tt.status, status)

			parseResponse(data, &userList)

			assert.EqualValues(t, tt.size, len(userList))

			if len(userList) == 0 {
				return
			}

			assert.EqualValues(t, tt.username, userList[tt.index].Name)
			assert.EqualValues(t, tt.uid, userList[tt.index].UID)
			assert.EqualValues(t, tt.gid, userList[tt.index].GID)
			assert.EqualValues(t, tt.comment, userList[tt.index].Comment)
			assert.EqualValues(t, tt.home, userList[tt.index].Home)
			assert.EqualValues(t, tt.shell, userList[tt.index].Shell)
		})
	}
}

func TestReadUser(t *testing.T) {

	tests := []struct {
		name     string
		status   int
		username string
		uid      int64
		gid      int64
		comment  string
		home     string
		shell    string
	}{
		{"found user", http.StatusOK, "_mysql", 74, 74, "MySQL Server", "/var/empty", "/usr/bin/false"},
		{"user not found", http.StatusNotFound, "", 999, 999, "", "", ""},
	}

	s := initUserService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var user models.User
			uidStr := strconv.Itoa(int(tt.uid))
			data, status := getTestResponse(http.MethodGet, "users/"+uidStr, s.Read, map[string]string{"uid": uidStr})

			assert.EqualValues(t, tt.status, status)

			if status == http.StatusNotFound {
				return
			}

			parseResponse(data, &user)

			assert.EqualValues(t, tt.username, user.Name)
			assert.EqualValues(t, tt.uid, user.UID)
			assert.EqualValues(t, tt.gid, user.GID)
			assert.EqualValues(t, tt.comment, user.Comment)
			assert.EqualValues(t, tt.home, user.Home)
			assert.EqualValues(t, tt.shell, user.Shell)
		})
	}
}

// End user tests

// Begin group tests
func TestListGroups(t *testing.T) {

	tests := []struct {
		name      string
		size      int
		status    int
		index     int
		groupname string
		gid       int64
		members   []string
		path      string
	}{
		{"full list one", 125, http.StatusOK, 83, "_lda", 211, []string{}, "groups"},
		{"full list two", 125, http.StatusOK, 45, "_www", 70, []string{"_devicemgr", "_teamsserver"}, "groups"},
		{"found by groupname", 1, http.StatusOK, 0, "kmem", 2, []string{"root"}, "groups/query?name=kmem"},
		{"found by gid", 1, http.StatusOK, 0, "group", 16, []string{}, "groups/query?gid=16"},
		{"found by members", 2, http.StatusOK, 0, "certusers", 29, []string{"root", "_jabber", "_postfix", "_cyrus", "_calendar", "_dovecot"}, "groups/query?member=_jabber&member=_postfix"},
		{"found all memberless groups", 97, http.StatusOK, 96, "com.apple.access_ssh", 399, []string{}, "groups/query?member="},
		{"not found by gid", 0, http.StatusOK, 0, "", 0, nil, "groups/query?gid=9999"},
	}

	s := initGroupService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groupList := make([]models.Group, 0)
			data, status := getTestResponse(http.MethodGet, tt.path, s.List, nil)

			assert.EqualValues(t, tt.status, status)

			parseResponse(data, &groupList)

			assert.EqualValues(t, tt.size, len(groupList))

			if len(groupList) == 0 {
				return
			}

			assert.EqualValues(t, tt.groupname, groupList[tt.index].Name)
			assert.EqualValues(t, tt.gid, groupList[tt.index].GID)
			assert.Len(t, groupList[tt.index].Members, len(tt.members))

			for i, m := range tt.members {
				assert.EqualValues(t, m, groupList[tt.index].Members[i])
			}
		})
	}
}

func TestReadGroup(t *testing.T) {

	tests := []struct {
		name      string
		status    int
		groupname string
		gid       int64
		members   []string
	}{
		{"found group", http.StatusOK, "_detachedsig", 207, []string{"_locationd"}},
		{"group not found", http.StatusNotFound, "", 999, nil},
	}

	s := initGroupService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var group models.Group
			gidStr := strconv.Itoa(int(tt.gid))
			data, status := getTestResponse(http.MethodGet, "gid/"+gidStr, s.Read, map[string]string{"gid": gidStr})

			assert.EqualValues(t, tt.status, status)

			if status == http.StatusNotFound {
				return
			}

			parseResponse(data, &group)

			assert.EqualValues(t, tt.groupname, group.Name)
			assert.EqualValues(t, tt.gid, group.GID)
			assert.Len(t, group.Members, len(tt.members))

			for i, m := range tt.members {
				assert.EqualValues(t, m, group.Members[i])
			}
		})
	}
}

// End group tests
