package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"tp-tdl/user"

	"github.com/gorilla/mux"
)

type User user.User

type AppController struct {
	users      []*User
	auctions   []*Auction
	lastUserID int
}

func (app *AppController) authUser(currentUser *User) bool {
	//To be implemented
	return true
}

// Same code in FindUserByID & FindAuctionByID
func (app *AppController) FindUserByID(id string) *User {
	var foundUser *User
	users := app.GetRegisteredUsers()
	for _, user := range users {
		if user.ID == id {
			foundUser = user
			break
		} else {
			foundUser = nil
		}
	}
	if foundUser == nil {
		panic("User is invalid")
	}
	return foundUser
}

func (app *AppController) FindAuctionByID(id string) *Auction {
	auctions := app.GetRegisteredAuctions()
	var foundAuction *Auction
	for _, auction := range auctions {
		if auction.ID == id {
			foundAuction = auction
			break
		} else {
			foundAuction = nil
		}
	}
	if foundAuction == nil {
		panic("Auction is invalid")
	}
	return foundAuction
}

func (app *AppController) GetCurrentUserID() int {
	return app.lastUserID
}

func (app *AppController) IncrementUserID() {
	app.lastUserID++
}

func (app *AppController) GetRegisteredUsers() []*User {
	return app.users
}

func (app *AppController) AddNewUser(newUser *User) {
	app.users = append(app.users, newUser)
}

func (app *AppController) GetRegisteredAuctions() []*Auction {
	return app.auctions
}

type Auction struct {
	ID           string  `json:"id"`
	Seller       *User   `json:"seller"`
	Participants []*User `json:"participants"`
	currentOffer int     `json:"currentoffer"`
	isPublic     bool    `json:"ispublic"`
	hasStarted   bool    `json:"-"`
}

type AuctionType interface {
	updateOffer() int
	addParticipant() int
}

func createAuction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var auction Auction
	json.NewDecoder(r.Body).Decode(&auction)
	if auction.Participants == nil {
		auction.Participants = []*User{auction.Seller}
	}
	auctions := app.GetRegisteredAuctions()
	json.NewEncoder(w).Encode(auctions)
}

func getAuctions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	auctions := app.GetRegisteredAuctions()
	json.NewEncoder(w).Encode(auctions)
}

func joinAuction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	foundAuction := app.FindAuctionByID(params["auctionid"])
	foundUser := app.FindUserByID(params["userid"])
	foundAuction.Participants = append(foundAuction.Participants, foundUser)
}

func updateAuctionOffer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	auctionId := params["auctionid"]
	rawOffer := params["newoffer"]
	newOffer, err := strconv.Atoi(rawOffer)
	if err != nil {
		panic("Received offer is invalid")
	}
	auctionToUpdate := app.FindAuctionByID(auctionId)
	if auctionToUpdate.currentOffer < newOffer {
		auctionToUpdate.currentOffer = newOffer
	}
}

var app AppController

func main() {
	r := mux.NewRouter()

	app = AppController{[]*User{}, []*Auction{}, 0}

	r.HandleFunc("/users", user.CreateUser).Methods("POST")
	r.HandleFunc("/users", user.GetUsers).Methods("GET")
	r.HandleFunc("/auctions", createAuction).Methods("POST")
	r.HandleFunc("/auctions", getAuctions).Methods("GET")
	r.HandleFunc("/auctions/auctionid={auctionid}&userid={userid}", joinAuction).Methods("PUT")
	r.HandleFunc("/auctions/auctionid={auctionid}&userid={userid}&newoffer={newoffer}", updateAuctionOffer).Methods("PUT")
	r.HandleFunc("/auctions", createAuction).Methods("POST")

	http.ListenAndServe(":8000", r)
}
