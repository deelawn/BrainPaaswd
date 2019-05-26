package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/deelawn/BrainPaaswd/readers/file"
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

	// Initialize readers to read data from sources; in our case we are reading from files
	readers := map[string]func(source string) io.Reader{
		passwdPath: file.NewReader,
		groupPath:  file.NewReader,
	}

	baseService = services.NewService(passwdPath, groupPath, readers)
	userService = users.NewService(*baseService)
}
