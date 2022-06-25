package model

type Auction struct {
	ID          string    `json:"id" bson:"_id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	SellerID    string    `json:"seller" bson:"seller"`
	UserOffer   UserOffer `json:"user_offer" bson:"user_offer"`
	IsTimed     bool      `json:"is_timed" bson:"is_timed"` // Para saber si la subasta es por tiempo (automatica o manual)
	HasEnded    bool      `json:"-" bson:"has_ended"`       // true si la subasta terminó
	ImageURL    string    `json:"image_url" bson:"image_url"`
}

type UserOffer struct {
	CurrentOffer int    `json:"current_offer" bson:"current_offer"`
	UserID       string `json:"id" bson:"id"`
	Username     string `json:"username" bson:"username"`
}
