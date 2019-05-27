package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/deelawn/BrainPaaswd/readers"
	"github.com/deelawn/BrainPaaswd/readers/file"
	"github.com/deelawn/BrainPaaswd/services"
	"github.com/deelawn/BrainPaaswd/services/groups"
	"github.com/deelawn/BrainPaaswd/services/users"
)

const (
	passwdPath = "passwd"
	groupPath  = "group"
)

var (
	baseService  *services.Service
	groupService *groups.Service
	userService  *users.Service
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/users", userService.List)
	r.HandleFunc("/users/{uid:[0-9]+}", userService.Read)
	r.HandleFunc("/groups", groupService.List)
	r.HandleFunc("/groups/{gid:[0-9]+}", groupService.Read)

	srv := &http.Server{
		Handler: r,
		Addr:    ":8000",
	}

	fmt.Println("Server started!")
	log.Fatal(srv.ListenAndServe())
}

func init() {

	// Initialize readers to read data from sources; in our case we are reading from files
	readers := map[string]func(source string) readers.Resource{
		passwdPath: file.NewReader,
		groupPath:  file.NewReader,
	}

	baseService = services.NewService(passwdPath, groupPath, readers)
	groupService = groups.NewService(*baseService)
	userService = users.NewService(*baseService)
}
