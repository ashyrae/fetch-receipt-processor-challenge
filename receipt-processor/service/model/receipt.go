package model

import (
	"regexp"

	pb "github.com/ashyrae/fetch-receipt-processor-challenge/receipt-processor/api/proto"
)

type constraint = string

const (
	// Receipt
	id_regexp       constraint = "^\\S+$"
	retailer_regexp constraint = "^[\\w\\s\\-&]+$"
	total_regexp    constraint = "^\\d+\\.\\d{2}$"

	// Item
	shortDesc_regexp constraint = "^[\\w\\s\\-]+$"
	price_regexp     constraint = "^\\d+\\.\\d{2}$"
)

type Receipt struct {
	Retailer string
	Date     string
	Time     string
	Total    string
	Items    []Item
}

type Item struct {
	ShortDescription string
	Price            string
}

type Points int32

func ProcessReceipt(receipt *pb.Receipt) (validated Receipt, err error) {
	// parse receipt items
	receiptItems := make([]Item, 0)
	for _, item := range receipt.Items {
		parsed := Item{item.ShortDescription, item.Price}
		receiptItems = append(receiptItems, parsed)
	}

	// validate our fields
	if err := validateReceipt(receipt.Retailer, receipt.PurchaseDate, receipt.PurchaseTime, receipt.Total, receiptItems); err != nil {
		return Receipt{}, err
	} else {
		validated = Receipt{receipt.Retailer, receipt.PurchaseDate, receipt.PurchaseTime, receipt.Total, receiptItems}
		return validated, nil
	}
}

func validateReceipt(
	retailer string,
	date string,
	time string,
	total string,
	items []Item,
) (err error) {
	rv, err := validateStringField(retailer, retailer_regexp)
	if err != nil {
		return err
	} else if !rv {
		return ErrInternalServer("Receipt Retailer was invalid, but validator did not return an error")
	}

	tv, err := validateStringField(total, total_regexp)
	if err != nil {
		return err
	} else if !tv {
		return ErrInternalServer("Receipt Total was invalid, but validator did not return an error")
	}

	iv, err := validateItemsField(items)
	if err != nil {
		return err
	} else if !iv {
		return ErrInternalServer("Receipt Items were invalid, but validator did not return an error")
	}

	// TODO: @ashyrae - check the example receipts
	// date & time validation
	if date == "" {
		return ErrBadRequest("Receipt Purchase Date was invalid")
	}

	if time == "" {
		return ErrBadRequest("Receipt Purchase Time was invalid")
	}

	if rv && tv && iv {
		err = nil
	} else {
		err = ErrInternalServer("Receipt was invalid, but validator encountered no errors")
	}

	return err
}

func validateStringField(field string, regex constraint) (valid bool, err error) {
	if v, err := regexp.MatchString(field, regex); err != nil {
		return v, ErrBadRequest(err.Error())
	} else if !v {
		return false, ErrBadRequest("Receipt is invalid")
	} else {
		return true, nil
	}
}

func validateItemsField(items []Item) (valid bool, err error) {
	ctrValidated := 0
	for _, item := range items {
		sdValid, err := regexp.MatchString(item.ShortDescription, shortDesc_regexp)
		if err != nil {
			return false, ErrBadRequest(err.Error())
		} else if !sdValid {
			return false, ErrBadRequest("Item short description is invalid for item %s")
		}

		priceValid, err := regexp.MatchString(item.Price, price_regexp)
		if err != nil {
			return false, ErrBadRequest(err.Error())
		} else if !priceValid {
			return false, ErrBadRequest("Item Price is invalid")
		}

		if sdValid && priceValid {
			ctrValidated++
		}
	}
	if ctrValidated == len(items) {
		return true, nil
	} else {
		return false, ErrInternalServer("A Receipt Item was invalid, but validator did not return an error")
	}
}
