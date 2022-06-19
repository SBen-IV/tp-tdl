package controller

import (
	"fmt"
	"net/http"
	"tp-tdl/model"

	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User model.User

func isValidUser(user User) bool {
	if len(user.Firstname) == 0 || len(user.Lastname) == 0 || len(user.Username) < 3 || len(user.Password) < 3 {
		return false
	}

	return true
}

func validateAndInsertUser(users *UserDB, newUser User) (int, string) {
	// users lock
	result := users.collection.FindOne(ctx, bson.M{"username": newUser.Username})

	// User no existe, por lo tanto es válido
	if result.Err() != mongo.ErrNoDocuments {
		return http.StatusBadRequest, "Username ya existe"
	}

	newUser.ID = ksuid.New().String()

	_, err := users.collection.InsertOne(ctx, newUser)
	// users unlock

	if err != nil {
		fmt.Println(err)
		return http.StatusInternalServerError, "Internal error"
	}

	return http.StatusOK, ""
}

func addNewUser(users *UserDB, newUser User) (int, string) {
	if !isValidUser(newUser) {
		return http.StatusBadRequest, "Datos inválidos"
	}

	if status, err := validateAndInsertUser(users, newUser); err != "" {
		return status, err
	}

	return http.StatusOK, "OK"
}

func loginUser(users *UserDB, user User) bool {
	var expectedUser User
	// users lock
	validUser := users.collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&expectedUser)
	// users unlock

	if validUser != nil || expectedUser.Password != user.Password {
		return false
	}

	return true
}
