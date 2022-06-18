package main

import (
	"net/http"

	"tp-tdl/controller"
	"tp-tdl/middleware"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	app := controller.NewAppController()

	defer app.Disconnect()

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
	private.HandleFunc("/auctions/auctionid={auctionid}", app.DeleteAuction).Methods("DELETE")
	private.HandleFunc("/auctions/auctionid={auctionid}&userid={userid}", app.JoinAuction).Methods("PUT")
	private.HandleFunc("/auctions/auctionid={auctionid}&userid={userid}&newoffer={newoffer}", app.UpdateAuctionOffer).Methods("PUT")

	/* 	r.HandleFunc("/users", user.GetUsers).Methods("GET")

	 */

	http.ListenAndServe(":8000", r)
}
