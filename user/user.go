package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"tp-tdl/token"
)

type User struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

var users = map[string]User{
	"mike_j": {
		ID:        "0",
		Firstname: "Mike",
		Lastname:  "Johnson",
		Username:  "mike_j",
		Password:  "1234",
	},
	"mc_clown": {
		ID:        "1",
		Firstname: "Ronald",
		Lastname:  "McDonald",
		Username:  "mc_clown",
		Password:  "burger_king",
	},
}

var (
	lastUserID = 0
)

func validateUser(user User) bool {
	if len(user.Firstname) == 0 || len(user.Lastname) == 0 || len(user.Username) < 3 || len(user.Password) < 3 {
		return false
	}

	_, exist := users[user.Username]

	return !exist
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !validateUser(newUser) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newUser.ID = strconv.Itoa(lastUserID)
	lastUserID++
	users[newUser.Username] = newUser
	w.WriteHeader(http.StatusOK)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func Login(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Content-Type", "application/json")

	var user User

	json.NewDecoder(r.Body).Decode(&user)

	expectedUser, validUser := users[user.Username]

	if !validUser || expectedUser.Password != user.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenStr, err := token.CreateToken(expectedUser.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.Header().Set("Auth-Token", tokenStr)
	w.WriteHeader(http.StatusAccepted)
}

func Profile(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	w.Write([]byte(fmt.Sprintf("Profile of %v", username)))
}
