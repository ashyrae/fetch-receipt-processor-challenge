package model_test

import (
	"testing"

	"github.com/ashyrae/fetch-receipt-processor-challenge/receipt-processor/service/model"
)

const (
	testId  = "test-key"
	noExist = "im-not-real"
)

func TestReceiptDB_Create(t *testing.T) {
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

	if id, err := testDB.Create(r); err != nil {
		t.Errorf("Error encountered setting test receipt into DB: %d", err)
	} else if receipt := testDB.Store[id]; receipt != r {
		t.Errorf("Receipt set in DB is not identical to provided receipt: expected %v, received %v", r, receipt)
	}
}

func TestReceiptDB_Set(t *testing.T) {
	var testDB = model.ReceiptDB{Store: make(map[string]*model.Receipt)}
	testDB.Store["test-key"] = &model.Receipt{
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
		Awarded: false,
	}

	set := &model.Receipt{
		Retailer: "TestTarget",
		Date:     "2025-01-27",
		Time:     "01:23",
		Total:    "100.00",
		Items: []*model.Item{
			{
				ShortDescription: "An item at Target",
				Price:            "100.00",
			},
		},
		Awarded: true,
	}

	if id, err := testDB.Set(testId, set); err != nil {
		t.Errorf("Error encountered setting test receipt into DB: %d", err)
	} else if receipt := testDB.Store[id]; receipt != set {
		t.Errorf("Receipt set in DB is not identical to provided receipt: expected %v, received %v", set, receipt)
	}

	if _, err := testDB.Set("", set); err == nil {
		t.Error("Expected BadRequest error was not encountered when not providing an id to Set")
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

	testDB.Store[testId] = r

	if receipt, err := testDB.Get(testId); err != nil {
		t.Errorf("Error encountered getting test receipt from DB: %d", err)
	} else if receipt != r {
		t.Errorf("Receipt set in DB is not identical to provided receipt: expected %v, received %v", r, receipt)
	}

	if _, err := testDB.Get(noExist); err == nil {
		t.Error("Expected BadRequest error not encountered fetching nonexistent receipt from DB")
	}
}
