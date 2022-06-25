package controller

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"sync"
	"tp-tdl/model"
	"tp-tdl/token"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()

type UserDB struct {
	mu         sync.Mutex
	collection *mongo.Collection
}

type AuctionDB struct {
	mu         sync.Mutex
	collection *mongo.Collection
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
	tmpl_home                  = "index"
	tmpl_main_hub              = "allAuctions"
	tmpl_new_auction           = "newAuction"
	tmpl_auction_detail        = "auctionDetail"
	tmpl_user_profile          = "userProfile"
	tmpl_auction_detail_seller = "auctionDetailSeller"
)

var templates = map[string]*template.Template{
	tmpl_home:                  nil,
	tmpl_main_hub:              nil,
	tmpl_new_auction:           nil,
	tmpl_auction_detail:        nil,
	tmpl_user_profile:          nil,
	tmpl_auction_detail_seller: nil,
}

/* Database initialization */
func connectToDB() (*mongo.Client, error) {
	err := godotenv.Load()

	if err != nil {
		panic("Can't load .env")
	}

	user_db, password_db, uri_db := os.Getenv("USER_DB"), os.Getenv("PASSWORD_DB"), os.Getenv("URI_DB")
	uri := fmt.Sprintf(uri_db, user_db, password_db)

	return mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
}

/* Templates initialization */
func initializeTemplates() {
	for key := range templates {
		file_name := "templates/" + key + ".html"
		templates[key] = template.Must(template.ParseFiles(file_name, "templates/css/styles.css"))
	}
}

/* App controller initalization */
func NewAppController() *AppController {
	initializeTemplates()

	client, err := connectToDB()

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to DB")

	app := AppController{
		db: &Database{
			client: client,
			userDB: &UserDB{
				collection: client.Database("the_blues").Collection("users"),
				// mutex does not need to be initialized
			},
			auctionDB: &AuctionDB{
				collection: client.Database("the_blues").Collection("auctions"),
				// mutex does not need to be initialized
			},
		},
	}

	return &app
}

/* Disconnects the database */
func (app *AppController) Disconnect() error {
	fmt.Println("Disconnecting DB...")
	err := app.db.client.Disconnect(ctx)
	fmt.Println("DB disconnected.")
	return err
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
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		SellerID:    r.Header.Get("user_id"),
		UserOffer: model.UserOffer{
			CurrentOffer: curr_offer,
			UserID:       r.Header.Get("user_id"),
			Username:     r.Header.Get("username"),
		},
		IsTimed:  is_timed,
		HasEnded: false,
		ImageURL: r.FormValue("imageurl"),
	}

	createAuction(app.db.auctionDB, &auction)

	// w.WriteHeader(status)
	http.Redirect(w, r, "/auctions/"+auction.ID, http.StatusSeeOther)
}

func (app *AppController) GetAuction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	auction_id := params["auction-id"]

	fmt.Println("Auction id:", auction_id)

	auction, err := getAuction(app.db.auctionDB, auction_id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if auction.SellerID == r.Header.Get("user_id") {
		templates[tmpl_auction_detail_seller].Execute(w, auction)
	} else {
		templates[tmpl_auction_detail].Execute(w, auction)
	}

	fmt.Println("Sending...")

}

func (app *AppController) GetAllAuctions(w http.ResponseWriter, r *http.Request) {
	auctions := getAllAuctions(app.db.auctionDB)

	templates[tmpl_main_hub].Execute(w, auctions)
}

func (app *AppController) UpdateAuctionOffer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	auction_id := params["auction-id"]

	auction, err := getAuction(app.db.auctionDB, auction_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	offer, _ := strconv.Atoi(r.FormValue("offer"))

	user_offer := UserOffer{
		CurrentOffer: offer,
		UserID:       r.Header.Get("user_id"),
		Username:     r.Header.Get("username"),
	}

	updateAuctionOffer(app.db.auctionDB, &auction, user_offer)
	http.Redirect(w, r, "/auctions/"+auction.ID, http.StatusSeeOther)
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
		Firstname: r.FormValue("firstName"),
		Lastname:  r.FormValue("lastName"),
		Username:  r.FormValue("username"),
		Password:  r.FormValue("password"),
	}

	fmt.Println(newUser.Firstname, newUser.Lastname, newUser.Username, newUser.Password)

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

	session, err := token.Store.Get(r, "auth-token")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	session.Values["authorize"] = true
	session.Values["user_id"] = user_id
	session.Values["username"] = user.Username

	session.Save(r, w)

	w.WriteHeader(http.StatusAccepted)

	templates[tmpl_user_profile].Execute(w, nil)
}

func (app *AppController) Profile(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	w.Write([]byte(fmt.Sprintf("Profile of %v", username)))
}

func (app *AppController) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	username := r.Header.Get("username")

	if len(username) > 0 {
		w.Write([]byte(fmt.Sprintf("Welcome %v", username)))
		return
	}

	templates[tmpl_home].Execute(w, nil)
}
