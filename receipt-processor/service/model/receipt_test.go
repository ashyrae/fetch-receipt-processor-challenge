package model_test

import (
	"encoding/json"
	"testing"

	pb "github.com/ashyrae/fetch-receipt-processor-challenge/receipt-processor/api/proto"
	"github.com/ashyrae/fetch-receipt-processor-challenge/receipt-processor/service/model"
)

func Test_ProcessReceipt(t *testing.T) {
	type testCase struct {
		receipt     *pb.Receipt
		errExpected bool
	}

	var testCases = []testCase{
		{
			receipt: &pb.Receipt{
				Retailer:     "TestTarget",
				PurchaseDate: "2025-01-21",
				PurchaseTime: "13:43",
				Items: []*pb.Item{
					{
						ShortDescription: "An item at Target",
						Price:            "40.29",
					},
				},
				Total: "40.29",
			},
			errExpected: false,
		},
		{
			receipt: &pb.Receipt{
				Retailer:     "TestSuperstore",
				PurchaseDate: "2025-01-21",
				PurchaseTime: "13:43",
				Items: []*pb.Item{
					{
						ShortDescription: "An item at Superstore",
						Price:            "100000.00",
					},
				},
				Total: "100000.00",
			},
			errExpected: false,
		},
		{
			receipt: &pb.Receipt{
				Retailer:     "",
				PurchaseDate: "2025-01-21",
				PurchaseTime: "13:43",
				Total:        "40.29",
				Items: []*pb.Item{
					{
						ShortDescription: "An item at someplace",
						Price:            "400000.29",
					},
				},
			},
			errExpected: true,
		},
		{
			receipt: &pb.Receipt{
				Retailer:     "TestWalmart",
				PurchaseDate: "",
				PurchaseTime: "13:43",
				Total:        "40.29",
				Items: []*pb.Item{
					{
						ShortDescription: "An item at Walmart",
						Price:            "40.29",
					},
				},
			},
			errExpected: true,
		},
		{
			receipt: &pb.Receipt{
				Retailer:     "TestWalmart",
				PurchaseDate: "0000-00-00",
				PurchaseTime: "13:43",
				Total:        "40.29",
				Items: []*pb.Item{
					{
						ShortDescription: "An item at Walmart",
						Price:            "40.29",
					},
				},
			},
			errExpected: true,
		},
		{
			receipt: &pb.Receipt{
				Retailer:     "TestHEB",
				PurchaseDate: "2025-01-21",
				PurchaseTime: "26:99",
				Total:        "40.29",
				Items: []*pb.Item{
					{
						ShortDescription: "An item at HEB",
						Price:            "40.29",
					},
				},
			},
			errExpected: true,
		},
		{
			receipt: &pb.Receipt{
				Retailer:     "TestSafeway",
				PurchaseDate: "2025-01-21",
				PurchaseTime: "13:43",
				Total:        "",
				Items: []*pb.Item{
					{
						ShortDescription: "An item at Safeway",
						Price:            "40.29",
					},
				},
			},
			errExpected: true,
		},
		{
			receipt: &pb.Receipt{
				Retailer:     "TestSafeway",
				PurchaseDate: "2025-01-21",
				PurchaseTime: "13:43",
				Total:        "40.29",
				Items: []*pb.Item{
					{
						ShortDescription: "",
						Price:            "",
					},
				},
			},
			errExpected: true,
		},
	}
	for i, tc := range testCases {
		if _, err := model.ProcessReceipt(tc.receipt); err != nil && !tc.errExpected {
			t.Errorf("Unexpected error processing receipt in test case %d: %v", i+1, err)
		} else if err == nil && tc.errExpected {
			t.Errorf("Did not receive expected error in test case %d", i+1)
		}
	}
}

func Test_ProcessReceipt_JSON(t *testing.T) {
	type testCase = string

	var testCases = []testCase{
		// Test case JSON is provided as example material by Fetch
		`{"retailer":"Walgreens","purchaseDate":"2022-01-02","purchaseTime":"08:13","total":"2.65","items":[{"shortDescription":"Pepsi - 12-oz","price":"1.25"},{"shortDescription":"Dasani","price":"1.40"}]}`,
		`{"retailer":"Target","purchaseDate":"2022-01-02","purchaseTime":"13:13","total":"1.25","items":[{"shortDescription": "Pepsi - 12-oz", "price": "1.25"}]}`,
	}
	for i, tc := range testCases {
		// We're not testing code outside our libraries
		// No code that isn't ours should be called during a test,
		// outside of test helper methods
		if r, err := unmarshalHelper(tc); err != nil {
			t.Errorf("Error unmarshaling JSON in test case %d: %v", i+1, err)
		} else if _, err := model.ProcessReceipt(r); err != nil {
			t.Errorf("Error processing receipt in test case %d: %v", i+1, err)
		}
	}
}

func Test_AwardPoints(t *testing.T) {
	type testCase struct {
		receipt *model.Receipt
		points  int64
	}

	var testCases = []testCase{
		{
			receipt: &model.Receipt{
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
			},
			points: 400,
		},
		{
			receipt: &model.Receipt{
				Retailer: "TestSuperstore",
				Date:     "2025-01-21",
				Time:     "13:43",
				Total:    "1000000000.00",
				Items: []*model.Item{
					{
						ShortDescription: "An item at Superstore",
						Price:            "1000000000.00",
					},
				},
			},
			points: 10000000000,
		},
		{
			receipt: &model.Receipt{
				Retailer: "TestSuperstore",
				Date:     "2025-01-21",
				Time:     "13:43",
				Total:    "0.01",
				Items: []*model.Item{
					{
						ShortDescription: "An item at Superstore",
						Price:            "0.01",
					},
				},
			},
			points: 25,
		},
	}

	for i, tc := range testCases {
		if award := model.AwardPoints(tc.receipt); award != tc.points {
			t.Errorf("Expected points were not awarded in test case %d: expected %d, got %d", i+1, tc.points, award)
		}
	}

}

func unmarshalHelper(body string) (r *pb.Receipt, err error) {
	if err = json.Unmarshal([]byte(body), &r); err != nil {
		return &pb.Receipt{}, err
	} else {
		return r, nil
	}
}
