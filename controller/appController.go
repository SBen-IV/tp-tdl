package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"
	"tp-tdl/token"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()

type UserDB struct {
	collection *mongo.Collection
	// Add mutex
}

type AuctionDB struct {
	collection *mongo.Collection
	// Add mutex
}

type Database struct {
	client    *mongo.Client
	userDB    *UserDB
	auctionDB *AuctionDB
}

type AppController struct {
	db *Database
}

type AuctionPageData struct {
	Auctions []Auction
}

const (
	tmpl_home        = "index"
	tmpl_main_hub    = "mainHub"
	tmpl_new_auction = "newAuction"
)

var templates = map[string]*template.Template{
	tmpl_home:        nil,
	tmpl_main_hub:    nil,
	tmpl_new_auction: nil,
}

func connectToDB() (*mongo.Client, error) {
	err := godotenv.Load()

	if err != nil {
		panic("Can't load .env")
	}

	user_db, password_db, uri_db := os.Getenv("USER_DB"), os.Getenv("PASSWORD_DB"), os.Getenv("URI_DB")
	uri := fmt.Sprintf(uri_db, user_db, password_db)

	return mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
}

func initializeTemplates() {
	for key := range templates {
		file_name := "templates/" + key + ".html"
		templates[key] = template.Must(template.ParseFiles(file_name))
	}
}

func NewAppController() *AppController {
	initializeTemplates()

	client, err := connectToDB()

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to DB")

	/* 	coll := client.Database("the_blues").Collection("users")
	   	var result User
	   	coll.FindOne(context.TODO(), bson.M{"username": "mark"}).Decode(&result)
	   	fmt.Println(result)
	   	title := "Back to the Future"
	   	var result bson.M
	   	// err = coll.FindOne(context.TODO(), bson.D{{"title", title}}).Decode(&result)
	   	if err == mongo.ErrNoDocuments {
	   		//fmt.Printf("No document was found with the title %s\n", title)
	   		panic("fail")
	   	}
	   	if err != nil {
	   		panic(err)
	   	}
	   	jsonData, err := json.MarshalIndent(result, "", "    ")
	   	if err != nil {
	   		panic(err)
	   	}
	   	fmt.Printf("%s\n", jsonData) */

	app := AppController{
		db: &Database{
			client: client,
			userDB: &UserDB{
				collection: client.Database("the_blues").Collection("users"),
				// Add mutex
			},
			auctionDB: &AuctionDB{
				collection: client.Database("the_blues").Collection("auctions"),
				// Add mutex
			},
		},
	}

	return &app
}

// Same code in FindUserByID & FindAuctionByID
func (app *AppController) FindUserByID(id string) *User {
	var foundUser *User
	/* 	users := app.GetRegisteredUsers()
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
	   	} */
	return foundUser
}

func (app *AppController) FindAuctionByID(id string) *Auction {
	/* 	auctions := app.GetRegisteredAuctions() */
	var foundAuction *Auction
	/* 	for _, auction := range auctions {
	   		if auction.ID == id {
	   			foundAuction = auction
	   			break
	   		} else {
	   			foundAuction = nil
	   		}
	   	}
	   	if foundAuction == nil {
	   		panic("Auction is invalid")
	   	} */
	return foundAuction
}

type AuctionType interface {
	updateOffer() int
	addParticipant() int
}

/* ============================================================================================= */
/* ============================================================================================= */
/* ============================================================================================= */

func (app *AppController) CreateAuction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var auction Auction
	json.NewDecoder(r.Body).Decode(&auction)

	createAuction(app.db.auctionDB, auction)
	/* 	if auction.Participants == nil {
		auction.Participants = []*User{auction.Seller}
	} */
	/* 	auctions := app.GetRegisteredAuctions() */
	/* 	json.NewEncoder(w).Encode(auctions) */
}

func (app *AppController) GetAuctions(w http.ResponseWriter, r *http.Request) {
	auctions := getAllAuctions(app.db.auctionDB)

	tmpl, err := template.ParseFiles("templates/mainHub.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, auctions)
}

func (app *AppController) JoinAuction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	/* 	params := mux.Vars(r)
	 	foundAuction := App.FindAuctionByID(params["auctionid"])
		foundUser := App.FindUserByID(params["userid"])
		foundAuction.Participants = append(foundAuction.Participants, foundUser) */
}

func (app *AppController) UpdateAuctionOffer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	/* 	params := mux.Vars(r)
	   	auctionId := params["auctionid"]
	   	rawOffer := params["newoffer"]
	   	newOffer, err := strconv.Atoi(rawOffer)
	   	if err != nil {
	   		panic("Received offer is invalid")
	   	}
	   	auctionToUpdate := App.FindAuctionByID(auctionId) */
	/* 	if auctionToUpdate.currentOffer < newOffer {
		auctionToUpdate.currentOffer = newOffer
	} */
}

func (app *AppController) DeleteAuction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	auction_id := params["auctionid"]

	deleteAuction(app.db.auctionDB, auction_id)

	w.WriteHeader(http.StatusOK)
}

/*
	User
*/

func (app *AppController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	status, msg := addNewUser(app.db.userDB, newUser)

	w.WriteHeader(status)
	w.Write([]byte(msg))
}

func (app *AppController) Login(w http.ResponseWriter, r *http.Request) {
	var user = User{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	if !loginUser(app.db.userDB, user) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenStr, err := token.CreateToken(user.Username)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Auth-Token", tokenStr)
	w.WriteHeader(http.StatusAccepted)
	// http.redirect("/profile")
}

func (app *AppController) Profile(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	w.Write([]byte(fmt.Sprintf("Profile of %v", username)))
}

func (app *AppController) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Println("Hello! world!")
	/* 	t, err := template.ParseFiles("templates/index.html")

	   	if err != nil {
	   		return
	   	}

	   	t.Execute(w, nil) */
	templates[tmpl_home].Execute(w, nil)
}

func (app *AppController) Disconnect() error {
	fmt.Println("Disconnecting DB...")
	err := app.db.client.Disconnect(ctx)
	fmt.Println("DB disconnected.")
	return err
}
