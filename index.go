package main

import (
	"log"
	"net/http"
	"server/users"

	"github.com/gorilla/mux"
)

func main() {

	port := ":8080"
	routes := mux.NewRouter()

	routes.HandleFunc("/users", users.CreateUser).Methods(http.MethodPost)
	routes.HandleFunc("/users", users.GetUsers).Methods(http.MethodGet)
	routes.HandleFunc("/users/{id}", users.GetUserById).Methods(http.MethodGet)
	routes.HandleFunc("/users/{id}", users.UpdateUser).Methods(http.MethodPut)
	routes.HandleFunc("/users/{id}", users.DeleteUser).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(port, routes))
}
