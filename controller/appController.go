package controller

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
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
	tmpl_layout                = "layout"
	tmpl_navbar                = "navbar"
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
	_, err := os.Stat(".env")

	if err != nil {
		fmt.Println("Loading envs from os")
	} else {
		fmt.Println("Loading envs from .env")
		godotenv.Load()
	}

	user_db, user_exist := os.LookupEnv("USER_DB")

	if !user_exist {
		panic("User does not exist")
	}

	password_db, pass_exist := os.LookupEnv("PASSWORD_DB")

	if !pass_exist {
		panic("Password does not exist")
	}

	uri_db, uri_exist := os.LookupEnv("URI_DB")

	if !uri_exist {
		panic("URI does not exist")
	}

	uri := fmt.Sprintf(uri_db, user_db, password_db)

	return mongo.Connect(ctx, options.Client().ApplyURI(uri))
}

/* Templates initialization */
func initializeTemplates() {
	navbar := "templates/" + tmpl_navbar + ".gohtml"
	for key := range templates {
		file_name := "templates/" + key + ".gohtml"
		templates[key] = template.Must(template.ParseFiles(navbar, file_name, "templates/css/styles.css"))
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

/* ================================================================== */
/* ========================= Auction ================================ */
/* ================================================================== */

func (app *AppController) GetAuctionForm(w http.ResponseWriter, r *http.Request) {
	templates[tmpl_new_auction].ExecuteTemplate(w, tmpl_layout, nil)
}

func (app *AppController) CreateAuction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	curr_offer, _ := strconv.Atoi(r.FormValue("currentoffer"))
	var is_timed model.AuctionType

	if r.FormValue("auctionType") == "timed" {
		var duration time.Duration

		switch r.FormValue("auctionLength") {
		case "6h":
			duration = time.Second * 30
		case "12h":
			duration = time.Minute
		case "24h":
			duration = time.Minute * 2
		}

		is_timed = &AuctionTime{
			Start:    time.Now(),
			Duration: duration,
			IsOver:   false,
			IsTimed:  true,
		}
	} else {
		is_timed = &AuctionNoTime{
			IsOver:  false,
			IsTimed: false,
		}
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
		Type:     is_timed,
		ImageURL: r.FormValue("imageurl"),
	}

	createAuction(app.db.auctionDB, &auction)

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
		templates[tmpl_auction_detail_seller].ExecuteTemplate(w, tmpl_layout, auction)
	} else {
		templates[tmpl_auction_detail].ExecuteTemplate(w, tmpl_layout, auction)
	}

	fmt.Println("Sending...")

}

func (app *AppController) GetAllAuctions(w http.ResponseWriter, r *http.Request) {
	auctions := getAllAuctions(app.db.auctionDB)

	templates[tmpl_main_hub].ExecuteTemplate(w, tmpl_layout, auctions)
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

	if !auction.Type.HasEnded() {
		offer, _ := strconv.Atoi(r.FormValue("offer"))

		user_offer := UserOffer{
			CurrentOffer: offer,
			UserID:       r.Header.Get("user_id"),
			Username:     r.Header.Get("username"),
		}

		updateAuctionOffer(app.db.auctionDB, &auction, user_offer)
	}

	http.Redirect(w, r, "/auctions/"+auction.ID, http.StatusSeeOther)
}

func (app *AppController) EndAuction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	auction_id := params["auction-id"]

	auction, err := getAuction(app.db.auctionDB, auction_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	endAuction(app.db.auctionDB, auction.ID)

	http.Redirect(w, r, "/auctions/"+auction.ID, http.StatusSeeOther)
}

func (app *AppController) DeleteAuction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	auction_id := params["auction-id"]

	deleteAuction(app.db.auctionDB, auction_id)

	w.WriteHeader(http.StatusOK)
}

/* =============================================================== */
/* ========================= User ================================ */
/* =============================================================== */

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

	templates[tmpl_user_profile].ExecuteTemplate(w, tmpl_layout, user)
}

func (app *AppController) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := token.Store.Get(r, "auth-token")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	session.Values["authorize"] = false

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *AppController) Profile(w http.ResponseWriter, r *http.Request) {
	user := User{
		Username: r.Header.Get("username"),
	}
	templates[tmpl_user_profile].ExecuteTemplate(w, tmpl_layout, user)
}

func (app *AppController) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	username := r.Header.Get("username")

	if len(username) > 0 {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	templates[tmpl_home].ExecuteTemplate(w, tmpl_layout, nil)
}
