package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func UsersService(writer http.ResponseWriter, request *http.Request) {

	writer.Write([]byte("Users response"))
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/users", UsersService)

	srv := &http.Server{
		Handler: r,
		Addr:    ":8000",
	}

	fmt.Println("Server started!")
	log.Fatal(srv.ListenAndServe())
}
