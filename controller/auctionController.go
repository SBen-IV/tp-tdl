package controller

import (
	"fmt"
	"tp-tdl/model"

	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson"
)

type Auction model.Auction
type UserOffer model.UserOffer

/*
Devuelve todas las subastas disponibles como un AuctionPageData
*/
func getAllAuctions(auctions *AuctionDB) AuctionPageData {
	cur, _ := auctions.collection.Find(ctx, bson.M{})
	var result AuctionPageData
	for cur.Next(ctx) {
		var auc Auction
		err := cur.Decode(&auc)

		if err != nil {
			fmt.Println(err)
			return result
		}

		result.Auctions = append(result.Auctions, auc)
	}

	return result
}

/*
Dado un auction_id devuelve la información completa de un auction y un error en caso de haberlo
*/
func getAuction(auctions *AuctionDB, auction_id string) (Auction, error) {
	auctions.mu.Lock()
	result := auctions.collection.FindOne(ctx, bson.M{"_id": auction_id})
	auctions.mu.Unlock()

	if result.Err() != nil {
		return Auction{}, result.Err()
	}

	var auction Auction
	result.Decode(&auction)
	return auction, nil
}

/*
Dado un auction_id, lo elimina de la base de datos auctions
*/
func deleteAuction(auctions *AuctionDB, auction_id string) {
	auctions.mu.Lock()
	cur := auctions.collection.FindOneAndDelete(ctx, bson.M{"id": auction_id})
	auctions.mu.Unlock()
	fmt.Println(cur)
}

/*
Dado un auction se guarda en la base de datos auctions.
Devuelve nil si todo salio bien, o un error en caso contrario
*/
func createAuction(auctions *AuctionDB, auction *Auction) error {
	auctions.mu.Lock()
	auction.ID = ksuid.New().String()

	_, err := auctions.collection.InsertOne(ctx, auction)
	auctions.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

/*
Dado un auction y una oferta, si esta es mayor a la oferta actual entonces se actualiza la oferta,
 en caso contrario no hace nada
*/
func updateAuctionOffer(auctions *AuctionDB, auction *Auction, user_offer UserOffer) {
	auctions.mu.Lock()
	if auction.UserOffer.CurrentOffer < user_offer.CurrentOffer {
		filter := bson.M{"_id": auction.ID}
		update := bson.M{"$set": bson.M{"user_offer": user_offer}}
		auctions.collection.UpdateOne(ctx, filter, update)
		fmt.Println("Update ", auction.ID)
	}

	fmt.Println("Curr offer", auction.UserOffer.CurrentOffer, "new offer", user_offer.CurrentOffer)
	auctions.mu.Unlock()
}

/*
Dado un auction_id setea el campo is_over en true indicando que la subasta terminó
*/
func endAuction(auctions *AuctionDB, auction_id string) {
	auctions.mu.Lock()
	auctions.collection.UpdateOne(ctx, bson.M{"_id": auction_id}, bson.M{"$set": bson.M{"is_over": true}})
	auctions.mu.Unlock()
}
