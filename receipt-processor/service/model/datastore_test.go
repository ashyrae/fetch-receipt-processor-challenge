package model_test

import (
	"testing"

	"github.com/ashyrae/fetch-receipt-processor-challenge/receipt-processor/service/model"
)

const (
	newId = ""
)

func TestReceiptDB_Set(t *testing.T) {
	var testDB = model.ReceiptDB{Store: make(map[string]*model.Receipt)}

	r := &model.Receipt{
		Retailer: "TestTarget",
		Date:     "2025-01-21",
		Time:     "13:43",
		Total:    "40.29",
		Items: []*model.Item{
			{
				ShortDescription: "An item at Target",
				Price:            "40.29",
			},
		},
	}

	if id, err := testDB.Set(newId, r); err != nil {
		t.Errorf("Error encountered setting test receipt into DB: %d", err)
	} else if receipt := testDB.Store[id]; receipt != r {
		t.Errorf("Receipt set in DB is not identical to provided receipt: expected %v, received %v", r, receipt)
	}
}

func TestReceiptDB_Get(t *testing.T) {
	var testDB = model.ReceiptDB{Store: make(map[string]*model.Receipt)}

	r := &model.Receipt{
		Retailer: "TestTarget",
		Date:     "2025-01-21",
		Time:     "13:43",
		Total:    "40.29",
		Items: []*model.Item{
			{
				ShortDescription: "An item at Target",
				Price:            "40.29",
			},
		},
	}

	testDB.Store["test-key"] = r

	if receipt, err := testDB.Get("test-key"); err != nil {
		t.Errorf("Error encountered getting test receipt from DB: %d", err)
	} else if receipt != r {
		t.Errorf("Receipt set in DB is not identical to provided receipt: expected %v, received %v", r, receipt)
	}
}
