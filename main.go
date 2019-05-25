package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func UserService(writer http.ResponseWriter, request *http.Request) {

	writer.Write([]byte("Users response"))
}

const (
	passwdPath = "/etc/passwd"
	groupPath  = "/etc/groups"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/users", UserService)

	srv := &http.Server{
		Handler: r,
		Addr:    ":8000",
	}

	fmt.Println("Server started!")
	log.Fatal(srv.ListenAndServe())
}
