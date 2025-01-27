package receipt_service

import (
	ctx "context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/ashyrae/fetch-receipt-processor-challenge/receipt-processor/api/proto"
	model "github.com/ashyrae/fetch-receipt-processor-challenge/receipt-processor/service/model"
)

type ReceiptService struct {
	pb.UnimplementedReceiptServiceServer
	db *model.ReceiptDB
}

func (s *ReceiptService) ProcessReceipt(ctx ctx.Context, req *pb.ProcessReceiptRequest) (res *pb.ProcessReceiptResponse, err error) {
	r := &pb.Receipt{
		Retailer:     req.Retailer,
		PurchaseDate: req.PurchaseDate,
		PurchaseTime: req.PurchaseTime,
		Items:        req.Items,
		Total:        req.Total,
	}

	if rec, err := model.ProcessReceipt(r); err != nil {
		return &pb.ProcessReceiptResponse{}, err
	} else if id, err := s.db.Create(&rec); err != nil {
		return &pb.ProcessReceiptResponse{}, err
	} else {
		return &pb.ProcessReceiptResponse{Id: id}, nil
	}
}

func (s *ReceiptService) AwardPoints(ctx ctx.Context, req *pb.AwardPointsRequest) (res *pb.AwardPointsResponse, err error) {
	// validate that the request actually contains an id of a processed receipt
	if receipt, err := s.db.Get(req.Id); err != nil {
		return &pb.AwardPointsResponse{}, err
	} else {
		// check if receipt was already awarded once
		if receipt.Awarded {
			// bail out, award nothing
			return &pb.AwardPointsResponse{Points: &pb.Points{Points: 0}}, nil
		} else {
			// proceed with award
			award := model.AwardPoints(receipt)
			// flag the receipt as awarded
			awardedReceipt := receipt
			awardedReceipt.Awarded = true
			if _, err := s.db.Set(req.Id, awardedReceipt); err != nil {
				return &pb.AwardPointsResponse{}, err
			} else {
				return &pb.AwardPointsResponse{Points: &pb.Points{Points: award}}, nil
			}
		}
	}
}

func NewService() (srv *grpc.Server) {
	// create the server
	srv = grpc.NewServer()
	// put it all together & register
	pb.RegisterReceiptServiceServer(srv, &ReceiptService{db: &model.ReceiptDB{Store: make(map[string]*model.Receipt)}})
	// enable server reflection
	reflection.Register(srv)
	return srv
}
