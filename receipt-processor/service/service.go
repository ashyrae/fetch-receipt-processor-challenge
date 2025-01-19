package receipt_service

import (
	ctx "context"

	pb "github.com/ashyrae/fetch-receipt-processor-challenge/receipt-processor/api/proto"
	model "github.com/ashyrae/fetch-receipt-processor-challenge/receipt-processor/service/model"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ReceiptService struct {
	pb.UnimplementedReceiptServiceServer
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
		// generate an ID for our receipt
		if receiptId, err := uuid.NewRandom(); err != nil {
			return &pb.ProcessReceiptResponse{}, model.ErrInternalServer(err.Error())
		} else {
			// validate against regex
			// then store in db & include id in response
			res = &pb.ProcessReceiptResponse{Id: receiptId.String()}
		}
	}

	return res, err
}

func (s *ReceiptService) AwardPoints(ctx ctx.Context, req *pb.AwardPointsRequest) (res *pb.AwardPointsResponse, err error) {
	// validate that the request actually contains an id of a processed receipt

	// otherwise, toss back NotFound

	// proceed with award

	return res, nil
}

func NewService() (s *grpc.Server) {
	serv := grpc.NewServer()
	pb.RegisterReceiptServiceServer(s, &ReceiptService{}) // Register the service
	// Enable server reflection
	reflection.Register(serv)
	return serv
}
