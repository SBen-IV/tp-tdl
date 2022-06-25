package main

import (
	"bufio"
	"net/http"
	"os"

	"tp-tdl/controller"
	"tp-tdl/middleware"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	app := controller.NewAppController()

	defer func() {
		err := app.Disconnect()
		if err != nil {
			panic(err)
		}
	}()

	r.PathPrefix("/templates/css/").Handler(http.StripPrefix("/templates/css/", http.FileServer(http.Dir("templates/css/"))))
	//http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./templates/css"))))

	public := r.NewRoute().Subrouter()
	private := r.NewRoute().Subrouter()

	// Public endpoints
	public.HandleFunc("/", app.Home).Methods("GET")
	public.HandleFunc("/users", app.CreateUser).Methods("POST")
	public.HandleFunc("/login", app.Login).Methods("POST")

	// Private enpoints
	private.Use(middleware.AuthUser)

	// Users
	private.HandleFunc("/profile", app.Profile).Methods("GET")
	private.HandleFunc("/logout", app.Logout).Methods("POST")

	// Auctions
	private.HandleFunc("/create-auction", app.GetAuctionForm).Methods("GET")
	private.HandleFunc("/create-auction", app.CreateAuction).Methods("POST")
	private.HandleFunc("/auctions", app.GetAllAuctions).Methods("GET")
	private.HandleFunc("/auctions/{auction-id}", app.GetAuction).Methods("GET")
	private.HandleFunc("/auctions/{auction-id}", app.DeleteAuction).Methods("DELETE")
	private.HandleFunc("/auctions/{auction-id}", app.UpdateAuctionOffer).Methods("POST")

	go http.ListenAndServe(":8000", r)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "q" {
			break
		}
	}
}
