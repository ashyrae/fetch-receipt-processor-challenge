package model

import (
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// regex
var id_regexp = regexp.MustCompile(`^\\S+$`).String()
var retailer_regexp = regexp.MustCompile(`^[\w\s&-]+$`).String()
var total_regexp = regexp.MustCompile(`^\\d+\\.\\d{2}$`).String()
var date_regexp = regexp.MustCompile(`^\d{4}\-(0?[1-9]|1[012])\-(0?[1-9]|[12][0-9]|3[01])$`).String()
var shortDesc_regexp = regexp.MustCompile(`^[\\w\\s\\-]+$`).String()
var price_regexp = regexp.MustCompile(`^\\d+\\.\\d{2}$`).String()

// time
var time_format = `23:59:59`

func validateReceipt(
	retailer string,
	date string,
	time string,
	total string,
	items []Item,
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

	if _, e := validateField(total, total_regexp); e != nil {
		invalid = append(invalid, "total")
	}

	if _, e := validateItems(items); e != nil {
		invalid = append(invalid, "items")
	}

	if len(invalid) > 0 {
		return ErrBadRequest(fmt.Sprintf("Receipt fields invalid; %s", invalid))
	} else {
		return nil
	}
}

func validateField(field string, regex string) (valid bool, err error) {
	if v, err := regexp.MatchString(field, regex); err != nil {
		return v, ErrBadRequest(err.Error())
	} else if !v {
		return false, ErrBadRequest(fmt.Sprintf("Receipt field %s is invalid:", field))
	} else {
		return true, nil
	}
}

func validateItems(items []Item) (valid bool, err error) {
	ctrValidated := 0
	for _, item := range items {
		sv, err := regexp.MatchString(item.ShortDescription, shortDesc_regexp)
		if err != nil {
			return false, ErrBadRequest(err.Error())
		} else if !sv {
			return false, ErrInternalServer("Item short description is invalid for item")
		}

		pv, err := regexp.MatchString(item.Price, price_regexp)
		if err != nil {
			return false, ErrBadRequest(err.Error())
		} else if !pv {
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
	_, err = time.Parse(time_format, t)
	if err != nil {
		return false, ErrInternalServer("Receipt Purchase Time could not be processed: " + err.Error())
	} else {
		return true, nil
	}
}

func idFactory() (id string, err error) {
	uid := uuid.New()
	v, err := regexp.MatchString(id_regexp, uid.String())
	if err != nil {
		return "", ErrInternalServer("Unable to generate new receipt ID: " + err.Error())
	} else if !v {
		return "", ErrInternalServer("Generated ID did not satisfy ID constraint")
	} else {
		id = uid.String()
		return id, nil
	}
}
