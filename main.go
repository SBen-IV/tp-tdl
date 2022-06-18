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

	r.PathPrefix("/templates/css/").Handler(http.StripPrefix("/templates/css/", http.FileServer(http.Dir("./templates/css/"))))

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

	// Auctions
	private.HandleFunc("/auctions", app.CreateAuction).Methods("POST")
	private.HandleFunc("/auctions", app.GetAuctions).Methods("GET")
	private.HandleFunc("/auctions/auctionid=", app.DeleteAuction).Methods("DELETE")
	private.HandleFunc("/auctions/auctionid={auctionid}", app.DeleteAuction).Methods("DELETE")
	private.HandleFunc("/auctions/auctionid={auctionid}&userid={userid}", app.JoinAuction).Methods("PUT")
	private.HandleFunc("/auctions/auctionid={auctionid}&userid={userid}&newoffer={newoffer}", app.UpdateAuctionOffer).Methods("PUT")

	go http.ListenAndServe(":8000", r)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "q" {
			break
		}
	}
}
