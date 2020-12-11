package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type User struct {
	ID        string `json:id`
	FirstName string `json:firstname`
	LastName  string `json:lastname`
}

var users []User

func initUsers() {
	users = []User{
		User{ID: "1", FirstName: "sakib", LastName: "alamin"},
		User{ID: "2", FirstName: "prangon", LastName: "majumder"},
		User{ID: "3", FirstName: "mehedi", LastName: "hasan"},
		User{ID: "4", FirstName: "sahadat", LastName: "hossain"},
	}
}

func getUsers(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(users)
}

func getUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for _, user := range users {
		if user.ID == params["id"] {
			json.NewEncoder(response).Encode(user)
			return
		}
	}
	json.NewEncoder(response).Encode(User{})
}

func addUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	// todo
}

func updateUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	// todo
}

func deleteUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	// todo
}

func handleRoutes(router *mux.Router) {
	router.HandleFunc("/api/users", getUsers).Methods("GET")
	router.HandleFunc("/api/user/{id}", getUser).Methods("GET")
	router.HandleFunc("/api/user/{id}", addUser).Methods("POST")
	router.HandleFunc("/api/user/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/api/user/{id}", deleteUser).Methods("DELETE")
}

func main() {
	fmt.Println("testing REST API!")
	initUsers()
	router := mux.NewRouter().StrictSlash(true)
	handleRoutes(router)
	http.ListenAndServe(":8080", router)
}
