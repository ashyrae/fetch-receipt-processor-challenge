package model

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// regex
var id_regexp = regexp.MustCompile(`^\S+$`)
var retailer_regexp = regexp.MustCompile(`[\w\s\-&]+$`)
var total_regexp = regexp.MustCompile(`^\d+\.\d{2}$`)
var date_regexp = regexp.MustCompile(`^\d{4}\-(0?[1-9]|1[012])\-(0?[1-9]|[12][0-9]|3[01])$`)
var shortDesc_regexp = regexp.MustCompile(`^[\w\s\-]+$`)
var price_regexp = regexp.MustCompile(`^\d+\.\d{2}$`)

func validateReceipt(
	retailer string,
	date string,
	time string,
	total string,
	items []*Item,
) (err error) {
	invalid := make([]string, 0)

	if _, e := validateField(retailer, retailer_regexp); e != nil {
		invalid = append(invalid, "retailer")
	}

	if _, e := validateField(date, date_regexp); e != nil {
		invalid = append(invalid, "date")
	}

	if _, e := validateTime(time); e != nil {
		invalid = append(invalid, "time")
	}

	if _, e := validateItems(items); e != nil {
		invalid = append(invalid, "items")

	} else if _, e := validateTotal(total, items); e != nil {
		invalid = append(invalid, "total")
	}

	if len(invalid) > 0 {
		return ErrBadRequest(fmt.Sprintf("Receipt is invalid: fields are invalid %s", invalid))
	} else {
		return nil
	}
}

func validateField(field string, regex *regexp.Regexp) (valid bool, err error) {
	if v := regex.MatchString(field); !v {
		return false, ErrBadRequest(fmt.Sprintf("Receipt field %s is invalid:", field))
	} else {
		return true, nil
	}
}

func validateTotal(total string, items []*Item) (valid bool, err error) {
	if v := total_regexp.MatchString(total); !v {
		return false, ErrBadRequest(fmt.Sprintf("Receipt total is invalid: %d", err))
	} else {
		var reconcile float64
		for _, item := range items {
			priceFloat, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return false, ErrInternalServer(fmt.Sprintf("error reconciling receipt total and item prices: %d", err))
			}
			reconcile = reconcile + priceFloat
		}
		if total != strconv.FormatFloat(reconcile, 'f', -1, 64) {
			ErrBadRequest("Receipt total is invalid: total price and item prices do not match")
		}
	}

	return true, nil
}

func validateItems(items []*Item) (valid bool, err error) {
	ctrValidated := 0
	for _, item := range items {
		sv := shortDesc_regexp.MatchString(item.ShortDescription)
		if !sv {
			return false, ErrInternalServer("Item short description is invalid for item")
		}

		pv := price_regexp.MatchString(item.Price)
		if !pv {
			return false, ErrBadRequest("Item Price is invalid")
		}

		if sv && pv {
			ctrValidated++
		}
	}
	if ctrValidated == len(items) {
		return true, nil
	} else {
		return false, ErrInternalServer("A Receipt Item was invalid, but validator did not return an error")
	}
}

func validateTime(t string) (valid bool, err error) {
	// sanitize to HH:MM:SS
	t = t + ":00"
	_, err = time.Parse(time.TimeOnly, t)
	if err != nil {
		return false, ErrInternalServer("Receipt Purchase Time could not be processed: " + err.Error())
	} else {
		return true, nil
	}
}

func idFactory() (id string, err error) {
	uid := uuid.New()
	v := id_regexp.MatchString(uid.String())
	if !v {
		return "", ErrInternalServer("Generated ID did not satisfy ID constraint")
	} else {
		id = uid.String()
		return id, nil
	}
}
