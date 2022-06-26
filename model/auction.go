package model

import "time"

type Auction struct {
	ID          string    `json:"id" bson:"_id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	SellerID    string    `json:"seller" bson:"seller"`
	UserOffer   UserOffer `json:"user_offer" bson:"user_offer"`
	/* 	Time        AuctionTime `json:"time" bson:"time"`   // Para saber si la subasta es por tiempo (automatica o manual)
	   	HasEnded    bool        `json:"-" bson:"has_ended"` // true si la subasta terminó */
	Type     AuctionType `json:"type" bson:"type"`
	ImageURL string      `json:"image_url" bson:"image_url"`
}

type UserOffer struct {
	CurrentOffer int    `json:"current_offer" bson:"current_offer"`
	UserID       string `json:"id" bson:"id"`
	Username     string `json:"username" bson:"username"`
}

type AuctionTime struct {
	Start    time.Time     `json:"start_time" bson:"start_time"`
	Duration time.Duration `json:"duration" bson:"duration"`
	IsOver   bool          `json:"is_over" bson:"is_over"`
	IsTimed  bool          `json:"is_timed" bson:"is_timed"`
}

type AuctionNoTime struct {
	IsOver  bool `json:"is_over" bson:"is_over"`
	IsTimed bool `json:"is_timed" bson:"is_timed"`
}

type AuctionType interface {
	HasEnded() bool
}
