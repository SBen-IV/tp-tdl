package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type Auction struct {
	ID           string  `json:"id"`
	Participants []*User `json:"participants"`
	currentOffer int     `json:"currentoffer"`
	isPublic     bool    `json:"ispublic"`
}

type AppController struct {
	users      []*User
	lastUserID int
}

type AuctionType interface {
	updateOffer() int
	addParticipant() int
}

func (*AppController) authUser(currentUser *User) bool {
	return true
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/user/auth", authUser).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	http.ListenAndServe(":8000", r)
}
