package controller

import (
	"fmt"
	"time"
	"tp-tdl/model"
)

type AuctionTime model.AuctionTime

func (aucTime *AuctionTime) HasEnded() bool {
	timePassed := time.Since(aucTime.Start)

	fmt.Println(timePassed)

	return timePassed < aucTime.Duration
}
