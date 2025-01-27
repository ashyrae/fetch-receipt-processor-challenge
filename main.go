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
	receipt_service "github.com/ashyrae/fetch-receipt-processor-challenge/receipt-processor/service"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO @ashyrae: Dedicated Error Types

func main() {
	// Initialize our loggers
	il := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	el := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Begin listening

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		el.Fatalf("Failed to begin listening: %v", err)
	}

	// Initialize the Receipt Service, DB, info logger, and error logger
	// We use a goroutine to allow shutdown to proceed in parallel
	s := receipt_service.NewService()
	go startServer(lis, s, il, el)

	// grpc-gateway to multiplex
	// http should go to 8081 to avoid protocol mishaps on 50051
	if conn, err := grpc.NewClient("0.0.0.0:50051", grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		// everything should explode - gracefully - if we can't reach the server internally
		el.Fatalln("Failed to dial gRPC server:", err)
	} else {
		// mux!
		gwmux := runtime.NewServeMux(
			runtime.WithMarshalerOption("application/json", &runtime.JSONPb{}),
		)
		if err = pb.RegisterReceiptServiceHandler(ctx.Background(), gwmux, conn); err != nil {
			el.Fatalln("Failed to register gateway:", err)
		} else {
			gwServer := &http.Server{
				Addr:    ":8081",
				Handler: gwmux,
			}

			log.Println("Serving Receipt Service gRPC-Gateway for REST on http://0.0.0.0:8081")
			if err := gwServer.ListenAndServe(); err != nil {
				el.Fatalf("Failed to serve gRPC-Gateway server for the Receipt Service: %v", err)
			}
		}
	}

	// Wait for the server to shut down gracefully when an OS signal is received
	waitForShutdown(s, il)
}

// Start the gRPC server and listen for incoming connections
func startServer(lis net.Listener, s *grpc.Server, il *log.Logger, el *log.Logger) {
	il.Println("gRPC Server starting on port 50051")
	if err := s.Serve(lis); err != nil && err != grpc.ErrServerStopped {
		el.Fatalf("Failed to serve: %v", err)
	}

}

// Wait for interrupt signal, then gracefully shut down server
func waitForShutdown(s *grpc.Server, il *log.Logger) {
	// Create a channel to receive OS signals
	sigs := make(chan os.Signal, 1)
	// Create a channel to receive a signal when server shutdown is complete
	done := make(chan bool, 1)

	// Register to receive specific signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine that will listen for signals
	go func() {
		sig := <-sigs
		il.Printf("Received signal: %s", sig)

		// Perform server shutdown
		il.Println("Shutting down server...")
		s.GracefulStop() // Gracefully stop the server
		il.Println("Server has been shut down.")

		// Notify the main goroutine that we're done
		done <- true
	}()

	// Wait for shutdown to complete
	<-done
	il.Println("Graceful shutdown completed")
}
