package controller

import (
	"fmt"
	"time"
	"tp-tdl/model"
)

type AuctionTime model.AuctionTime
type AuctionNoTime model.AuctionNoTime

func (aucTime *AuctionTime) HasEnded() bool {
	if aucTime.IsOver {
		return true
	}

	timePassed := time.Since(aucTime.Start)

	fmt.Println(timePassed)

	if timePassed < aucTime.Duration {
		aucTime.IsOver = true
	}

	return aucTime.IsOver
}

func (aucNoTime *AuctionNoTime) HasEnded() bool {
	return aucNoTime.IsOver
}
