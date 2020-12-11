package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// User struct
type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
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

// Credentials struct
type Credentials struct {
	Username string
	Password string
}

var authentication Credentials

func initAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	authentication.Username = os.Getenv("Username")
	authentication.Password = os.Getenv("Password")
	fmt.Println("Username: ", authentication.Username, "Password: ", authentication.Password)
}

func checkBasicAuthentication(currentUsername string, currentPassword string, isAuthenticated bool) bool {
	if !isAuthenticated {
		fmt.Println("No Authentication credentials! Username or Password provided!")
		return false
	}

	if authentication.Username != currentUsername || authentication.Password != currentPassword {
		fmt.Println("Wrong Username or Password! Not Authenticated!")
		return false
	}
	return true
}

func getUsers(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	if authenticated := checkBasicAuthentication(request.BasicAuth()); !authenticated {
		fmt.Println("Authentication Failed!")
		json.NewEncoder(response).Encode(User{})
		return
	}

	fmt.Println("Authentication successful!")
	json.NewEncoder(response).Encode(users)
}

func getUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	if authenticated := checkBasicAuthentication(request.BasicAuth()); !authenticated {
		fmt.Println("Authentication Failed!")
		json.NewEncoder(response).Encode(User{})
		return
	}
	fmt.Println("Authentication successful!")

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

	if authenticated := checkBasicAuthentication(request.BasicAuth()); !authenticated {
		fmt.Println("Authentication Failed!")
		json.NewEncoder(response).Encode(User{})
		return
	}
	fmt.Println("Authentication successful!")

	user := User{}
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", user)
	users = append(users, user)
	json.NewEncoder(response).Encode(users)
}

func updateUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	if authenticated := checkBasicAuthentication(request.BasicAuth()); !authenticated {
		fmt.Println("Authentication Failed!")
		json.NewEncoder(response).Encode(User{})
		return
	}
	fmt.Println("Authentication successful!")

	newUser := User{}
	if err := json.NewDecoder(request.Body).Decode(&newUser); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", newUser)
	params := mux.Vars(request)
	for index, user := range users {
		if user.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			users = append(users, newUser)
			break
		}
	}
	json.NewEncoder(response).Encode(users)
}

func deleteUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	if authenticated := checkBasicAuthentication(request.BasicAuth()); !authenticated {
		fmt.Println("Authentication Failed!")
		json.NewEncoder(response).Encode(User{})
		return
	}
	fmt.Println("Authentication successful!")

	params := mux.Vars(request)
	for index, user := range users {
		if user.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			break
		}
	}
	json.NewEncoder(response).Encode(users)
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
	initAuth()

	router := mux.NewRouter().StrictSlash(true)
	handleRoutes(router)
	log.Fatal(http.ListenAndServe(":8080", router))
}
