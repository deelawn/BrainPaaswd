package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/deelawn/BrainPaaswd/services"
	"github.com/deelawn/BrainPaaswd/services/users"
)

const (
	passwdPath = "passwd"
	groupPath  = "groups"
)

var (
	baseService *services.Service
	userService *users.Service
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/users", userService.List)
	r.HandleFunc("/users/{uid:[0-9]+}", userService.Read)

	srv := &http.Server{
		Handler: r,
		Addr:    ":8000",
	}

	fmt.Println("Server started!")
	log.Fatal(srv.ListenAndServe())
}

func init() {

	baseService = services.NewService(passwdPath, groupPath)
	userService = users.NewService(*baseService)
}
