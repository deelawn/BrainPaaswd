package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/deelawn/BrainPaaswd/readers/file"
	"github.com/deelawn/BrainPaaswd/services"
	"github.com/deelawn/BrainPaaswd/services/groups"
	"github.com/deelawn/BrainPaaswd/services/users"
	"github.com/deelawn/BrainPaaswd/storage"
)

const (
	passwdPath = "passwd"
	groupPath  = "group"
	port       = ":8000"
)

var (
	baseService  *services.Service
	groupService *groups.Service
	userService  *users.Service
)

func main() {

	r := mux.NewRouter()

	// Users service
	r.HandleFunc("/users", userService.List)
	r.HandleFunc("/users/{uid:[0-9]+}", userService.Read)
	r.HandleFunc("/users/{uid:[0-9]+}/groups", userService.Groups)
	r.HandleFunc("/users/query", userService.List)

	// Groups service
	r.HandleFunc("/groups", groupService.List)
	r.HandleFunc("/groups/{gid:[0-9]+}", groupService.Read)
	r.HandleFunc("/groups/query", groupService.List)

	srv := &http.Server{
		Handler: r,
		Addr:    port,
	}

	fmt.Println("Server started!")
	log.Fatal(srv.ListenAndServe())
}

func init() {

	baseService = services.NewService()
	groupService = groups.NewService(*baseService, groupPath, storage.NewLocalCache(), file.NewReader)
	userService = users.NewService(*baseService, passwdPath, storage.NewLocalCache(), file.NewReader)
}
