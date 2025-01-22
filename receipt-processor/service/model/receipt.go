package model

import (
	"math"
	"strconv"

	pb "github.com/ashyrae/fetch-receipt-processor-challenge/receipt-processor/api/proto"
)

const (
	basePoints = 25
)

type Receipt struct {
	Retailer string
	Date     string
	Time     string
	Total    string
	Items    []*Item
}

type Item struct {
	ShortDescription string
	Price            string
}

type Points int64

func ProcessReceipt(receipt *pb.Receipt) (validated Receipt, err error) {
	// parse receipt items
	receiptItems := make([]*Item, 0)
	for _, item := range receipt.GetItems() {
		parsed := Item{item.GetShortDescription(), item.GetPrice()}
		receiptItems = append(receiptItems, &parsed)
	}

	rec := Receipt{receipt.GetRetailer(), receipt.GetPurchaseDate(), receipt.GetPurchaseTime(), receipt.GetTotal(), receiptItems}

	// validate our fields
	if err := validateReceipt(receipt.GetRetailer(), receipt.GetPurchaseDate(), receipt.GetPurchaseTime(), receipt.GetTotal(), receiptItems); err != nil {
		return Receipt{}, err
	} else {
		validated = rec
		return validated, nil
	}
}

func AwardPoints(r *Receipt) (awardPoints int64) {
	// Fetch rewards 25 points minimum per valid receipt,
	// even if there are no matched offers.
	// For simplicity, we will start at 25 points,
	// & award 10 extra per US dollar spent.
	var pending int64
	for _, item := range r.Items {
		// our data is sanitized, item prices conform to regex
		// since prices are decimals, parse as float64
		unadjusted, _ := strconv.ParseFloat(item.Price, 64)
		// round to the nearest whole number,
		// and convert to int64 to conform to API spec
		award := int64(math.Round(unadjusted))
		// ensure we account for any previously calculated dollars spent
		pending = pending + (award * 10)
		awardPoints = awardPoints + pending
	}
	if basePoints > awardPoints {
		return basePoints
	} else {
		return awardPoints
	}

}
