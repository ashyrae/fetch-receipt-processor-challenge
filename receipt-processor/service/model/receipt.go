package model

import (
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	pb "github.com/ashyrae/fetch-receipt-processor-challenge/receipt-processor/api/proto"
)

type Receipt struct {
	Retailer string
	Date     string
	Time     string
	Total    string
	Items    []*Item
	Awarded  bool
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

	rec := Receipt{receipt.GetRetailer(), receipt.GetPurchaseDate(), receipt.GetPurchaseTime(), receipt.GetTotal(), receiptItems, false}

	// validate our fields
	if err := validateReceipt(receipt.GetRetailer(), receipt.GetPurchaseDate(), receipt.GetPurchaseTime(), receipt.GetTotal(), receiptItems); err != nil {
		log.Printf("Error encountered validating receipt: %s", err.Error())
		return Receipt{}, err
	} else {
		validated = rec
		return validated, nil
	}
}

func AwardPoints(r *Receipt) (awardPoints int64) {
	var pendingPts int64

	// Points for the Retailer field
	pendingPts = pendingPts + int64(len(alphanumeric_regexp.FindAllString(r.Retailer, -1)))

	// Points for the Purchase Date field
	split := strings.Split(r.Date, "-")
	if parsed, _ := strconv.ParseInt(split[2], 10, 64); parsed%2 == 0 {
		// 6 points if the day in the purchase date is odd.
		pendingPts = pendingPts + 6
	}

	// Points for the Purchase Time field
	if t, _ := time.Parse(time.TimeOnly, r.Time+":00"); t.Hour() < 16 && t.Hour() > 14 {
		// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
		pendingPts = pendingPts + 10
	}

	// Points for the Total field
	if split := strings.Split(r.Total, "."); split[1] == "00" {
		// 50 points if the total is a round dollar amount with no cents.
		pendingPts = pendingPts + 50
	} else if split[1] == "25" || split[1] == "50" || split[1] == "75" {
		// 25 points if the total is a multiple of 0.25.
		pendingPts = pendingPts + 25
	}

	// Points for the Items field
	for itemNum, item := range r.Items {
		// 5 points for every two items on the receipt.
		if itemNum%2 == 0 {
			pendingPts = pendingPts + 5
		}

		// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer.
		// The result is the number of points earned.
		if len(item.ShortDescription)%3 == 0 {
			// our data is sanitized, item prices conform to regex
			// since prices are decimals, parse as float64
			unadjusted, _ := strconv.ParseFloat(item.Price, 64)
			adjusted := int64(math.Round(unadjusted * 0.2))
			// ensure we account for any previously calculated dollars spent
			pendingPts = pendingPts + (adjusted * 10)
		} else {
			// our data is sanitized, item prices conform to regex
			// since prices are decimals, parse as float64
			unadjusted, _ := strconv.ParseFloat(item.Price, 64)
			// round to the nearest whole number,
			// and convert to int64 to conform to API spec
			award := int64(math.Round(unadjusted))
			// ensure we account for any previously calculated dollars spent
			pendingPts = pendingPts + (award * 10)
		}
	}

	// finalize our award amount
	awardPoints = awardPoints + pendingPts
	return

}
