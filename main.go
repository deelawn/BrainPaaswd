package main

import (
	"flag"
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
	defaultPasswdPath = "passwd"
	defaultGroupPath  = "group"
	defaultPort       = "8000"
)

func main() {

	passwdPath := flag.String("passwd", defaultPasswdPath, "path to passwd file")
	groupPath := flag.String("group", defaultGroupPath, "path to group file")
	port := flag.String("port", defaultPort, "the port to run the web server on")
	flag.Parse()

	groupService := groups.NewService(services.NewService(*groupPath, storage.NewLocalCache(), file.NewReader))
	userService := users.NewService(services.NewService(*passwdPath, storage.NewLocalCache(), file.NewReader))

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
		Addr:    ":" + *port,
	}

	fmt.Printf("Server started on port %s!\n", *port)
	log.Fatal(srv.ListenAndServe())
}
