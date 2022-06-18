package controller

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func getAllAuctions(auctions AuctionDB) []Auction {
	cur, _ := auctions.collection.Find(ctx, bson.M{})
	var result []Auction
	for cur.Next(ctx) {
		var auc Auction
		err := cur.Decode(&auc)

		if err != nil {
			fmt.Println(err)
			return result
		}

		result = append(result, auc)
	}

	fmt.Println(result)

	return result
}

func deleteAuction(auctions AuctionDB, auction_id string) {
	cur := auctions.collection.FindOneAndDelete(ctx, bson.M{"id": auction_id})
	fmt.Println(cur)
}

func createAuction(auctions AuctionDB, auction Auction) {
	// Validar campos

	cur, err := auctions.collection.InsertOne(ctx, auction)
	fmt.Println(err)
	fmt.Println(cur)
}
