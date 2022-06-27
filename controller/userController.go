package controller

import (
	"errors"
	"fmt"
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

func validateAndInsertUser(users *UserDB, newUser User) error {
	users.mu.Lock()
	result := users.collection.FindOne(ctx, bson.M{"username": newUser.Username})

	// User no existe, por lo tanto es válido
	if result.Err() != mongo.ErrNoDocuments {
		return errors.New("username")
	}

	newUser.ID = ksuid.New().String()

	_, err := users.collection.InsertOne(ctx, newUser)
	users.mu.Unlock()

	if err != nil {
		fmt.Println(err)
		return errors.New("server")
	}

	return nil
}

func addNewUser(users *UserDB, newUser User) error {
	if !isValidUser(newUser) {
		return errors.New("parameters")
	}

	if err := validateAndInsertUser(users, newUser); err != nil {
		return err
	}

	return nil
}

func loginUser(users *UserDB, user User) string {
	var expectedUser User
	users.mu.Lock()
	validUser := users.collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&expectedUser)
	users.mu.Unlock()

	if validUser != nil || expectedUser.Password != user.Password {
		return ""
	}

	return expectedUser.ID
}
