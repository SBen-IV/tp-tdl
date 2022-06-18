package controller

import (
	"fmt"
	"tp-tdl/model"

	"go.mongodb.org/mongo-driver/bson"
)

type Auction model.Auction

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

	fmt.Println(result)

	return result
}

func deleteAuction(auctions *AuctionDB, auction_id string) {
	cur := auctions.collection.FindOneAndDelete(ctx, bson.M{"id": auction_id})
	fmt.Println(cur)
}

func createAuction(auctions *AuctionDB, auction Auction) {
	// Validar campos

	cur, err := auctions.collection.InsertOne(ctx, auction)
	fmt.Println(err)
	fmt.Println(cur)
}
