package model

type Auction struct {
	ID           string `json:"id" bson:"_id"`
	Title        string `json:"title" bson:"title"`
	Description  string `json:"description" bson:"description"`
	Seller       User   `json:"seller" bson:"seller"`
	Participants []User `json:"participants" bson:"participants"`
	CurrentOffer int    `json:"currentoffer" bson:"current_offer"`
	IsTimed      bool   `json:"ispublic" bson:"is_public"`
	HasEnded     bool   `json:"-" bson:"has_started"`
}

type AuctionPageData struct {
	Auctions []Auction
}
