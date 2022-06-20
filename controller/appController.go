package controller

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
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
	tmpl_home           = "index"
	tmpl_main_hub       = "mainHub"
	tmpl_new_auction    = "newAuction"
	tmpl_auction_detail = "auctionDetail"
)

var templates = map[string]*template.Template{
	tmpl_home:           nil,
	tmpl_main_hub:       nil,
	tmpl_new_auction:    nil,
	tmpl_auction_detail: nil,
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
		templates[key] = template.Must(template.ParseFiles(file_name, "templates/css/styles.css"))
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

type AuctionType interface {
	updateOffer() int
	addParticipant() int
}

/* ============================================================================================= */
/* ============================================================================================= */
/* ============================================================================================= */

func (app *AppController) GetAuctionForm(w http.ResponseWriter, r *http.Request) {
	templates[tmpl_new_auction].Execute(w, nil)
}

func (app *AppController) CreateAuction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	curr_offer, _ := strconv.Atoi(r.FormValue("currentoffer"))
	var is_timed = false

	if r.FormValue("auctionType") == "timed" {
		is_timed = true
	}

	var auction = Auction{
		Title:        r.FormValue("title"),
		Description:  r.FormValue("description"),
		SellerID:     r.Header.Get("user_id"),
		CurrentOffer: curr_offer,
		IsTimed:      is_timed,
		HasEnded:     false,
		ImageURL:     r.FormValue("imageurl"),
	}

	status := createAuction(app.db.auctionDB, &auction)

	w.WriteHeader(status)
}

func (app *AppController) GetAuction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	auction_id := params["auction-id"]

	auction, err := getAuction(app.db.auctionDB, auction_id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	templates[tmpl_auction_detail].Execute(w, auction)
}

func (app *AppController) GetAllAuctions(w http.ResponseWriter, r *http.Request) {
	auctions := getAllAuctions(app.db.auctionDB)

	templates[tmpl_main_hub].Execute(w, auctions)
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
	auction_id := params["auction-id"]

	deleteAuction(app.db.auctionDB, auction_id)

	w.WriteHeader(http.StatusOK)
}

/*
	User
*/

func (app *AppController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser = User{
		Firstname: r.FormValue("firstname"),
		Lastname:  r.FormValue("lastname"),
		Username:  r.FormValue("username"),
		Password:  r.FormValue("password"),
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

	user_id := loginUser(app.db.userDB, user)

	if len(user_id) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenStr, err := token.CreateToken(user_id)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Auth-Token", tokenStr)
	w.WriteHeader(http.StatusAccepted)
	app.GetAllAuctions(w, r)
}

func (app *AppController) Profile(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	w.Write([]byte(fmt.Sprintf("Profile of %v", username)))
}

func (app *AppController) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Println("Hello! world!")

	templates[tmpl_home].Execute(w, nil)
}

func (app *AppController) Disconnect() error {
	fmt.Println("Disconnecting DB...")
	err := app.db.client.Disconnect(ctx)
	fmt.Println("DB disconnected.")
	return err
}
