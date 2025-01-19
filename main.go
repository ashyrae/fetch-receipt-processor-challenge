package main

import (
	ctx "context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/ashyrae/fetch-receipt-processor-challenge/receipt-processor/api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type receiptService struct {
	pb.UnimplementedReceiptServiceServer
}

func (s *receiptService) ProcessReceipt(ctx ctx.Context, req *pb.ProcessReceiptRequest) (res *pb.ProcessReceiptResponse, err error) {
	// placeholder

	return res, nil
}

func AwardPoints(ctx ctx.Context, req *pb.AwardPointsRequest) (res *pb.AwardPointsResponse, err error) {
	// placeholder

	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterReceiptServiceServer(s, &receiptService{}) // Register the service
	// Enable server reflection
	reflection.Register(s)

	// Start serving in a goroutine to allow shutdown to proceed in parallel
	go startServer(s, lis)

	conn, err := grpc.NewClient("0.0.0.0:50051")
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterReceiptServiceHandler(ctx.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8081",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway for REST on http://0.0.0.0:8081")
	if err := gwServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed to serve gRPC-Gateway server: %v", err)
	}

	// Wait for the server to shut down gracefully when an OS signal is received
	waitForShutdown(s)
}

// Start the gRPC server and listen for incoming connections
func startServer(s *grpc.Server, lis net.Listener) {
	log.Println("Server started on port 50051")
	if err := s.Serve(lis); err != nil && err != grpc.ErrServerStopped {
		log.Fatalf("failed to serve: %v", err)
	}

}

// Wait for interrupt signal to gracefully shutdown the server
func waitForShutdown(s *grpc.Server) {
	// Create a channel to receive OS signals
	sigs := make(chan os.Signal, 1)
	// Create a channel to receive a signal when server shutdown is complete
	done := make(chan bool, 1)

	// Register to receive specific signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine that will listen for signals
	go func() {
		sig := <-sigs
		log.Printf("Received signal: %s", sig)

		// Perform server shutdown
		log.Println("Shutting down server...")
		s.GracefulStop() // Gracefully stop the server
		log.Println("Server has been shut down.")

		// Notify the main goroutine that we're done
		done <- true
	}()

	// Wait for shutdown to complete
	<-done
	log.Println("Graceful shutdown completed")
}
