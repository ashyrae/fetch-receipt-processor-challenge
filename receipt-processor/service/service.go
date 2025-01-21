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
	store *model.ReceiptDB
}

func (s *ReceiptService) ProcessReceipt(ctx ctx.Context, req *pb.ProcessReceiptRequest) (res *pb.ProcessReceiptResponse, err error) {
	// parse our receipt fields from JSON
	// no non-JSON should make it through due to header validation
	// so we can focus on ensuring it's *valid* JSON that conforms to schema
	// regex validation in `api.yml`
	// if not, toss back BadRequest & punt
	if _, err := model.ProcessReceipt(req.Receipt); err != nil {
		return &pb.ProcessReceiptResponse{}, err
	} else {
		res = &pb.ProcessReceiptResponse{Id: ""}

	}

	return res, err
}

func (s *ReceiptService) AwardPoints(ctx ctx.Context, req *pb.AwardPointsRequest) (res *pb.AwardPointsResponse, err error) {
	// validate that the request actually contains an id of a processed receipt
	if receipt, err := s.store.Get(req.Id); err != nil {
		return &pb.AwardPointsResponse{}, err
	} else {
		// proceed with award
		award := model.AwardPoints(receipt)
		return &pb.AwardPointsResponse{Points: &pb.Points{Points: award}}, nil
	}
}

func NewService() (s *grpc.Server) {
	// init our DB
	db := make(map[string]*model.Receipt)
	store := model.ReceiptDB{Store: db}
	// create the server
	srv := grpc.NewServer()
	// put it all together & register
	pb.RegisterReceiptServiceServer(s, &ReceiptService{store: &store})
	// enable server reflection
	reflection.Register(srv)
	return srv
}
