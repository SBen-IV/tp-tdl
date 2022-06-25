package controller

import (
	"fmt"
	"net/http"
	"tp-tdl/model"

	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson"
)

type Auction model.Auction
type UserOffer model.UserOffer

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

func getAuction(auctions *AuctionDB, auction_id string) (Auction, error) {
	result := auctions.collection.FindOne(ctx, bson.M{"_id": auction_id})

	if result.Err() != nil {
		return Auction{}, result.Err()
	}

	var auction Auction
	result.Decode(&auction)

	return auction, nil
}

func deleteAuction(auctions *AuctionDB, auction_id string) {
	cur := auctions.collection.FindOneAndDelete(ctx, bson.M{"id": auction_id})
	fmt.Println(cur)
}

func createAuction(auctions *AuctionDB, auction *Auction) int {
	auctions.mu.Lock()
	auction.ID = ksuid.New().String()

	_, err := auctions.collection.InsertOne(ctx, auction)
	auctions.mu.Unlock()

	if err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

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
