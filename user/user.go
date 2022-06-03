package user

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type User struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

var users = []User{
	{
		ID:        "0",
		Firstname: "John",
		Lastname:  "Doe",
		Username:  "johnny_d",
		Password:  "1234",
	},
}

var lastUserID = 0

func userExists(username string) bool {
	for _, user := range users {
		if user.Username == username {
			return true
		}
	}

	return false
}

func validateUser(user User) bool {
	if len(user.Firstname) == 0 || len(user.Lastname) == 0 || len(user.Username) < 3 || len(user.Password) < 3 {
		return false
	}

	return !userExists(user.Username)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newUser User
	json.NewDecoder(r.Body).Decode(&newUser)

	if !validateUser(newUser) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newUser.ID = strconv.Itoa(lastUserID)
	lastUserID++
	//app.IncrementUserID()
	//app.AddNewUser(&user)
	users = append(users, newUser)
	json.NewEncoder(w).Encode(newUser)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// users := app.GetRegisteredUsers()
	json.NewEncoder(w).Encode(users)
}
