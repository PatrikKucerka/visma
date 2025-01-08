package main

import (
	"log"
	"net"

	"github.com/vismaml/hiring/image-service-grpc/pkg/service"
	pb "github.com/vismaml/hiring/image-service-grpc/proto"
	"google.golang.org/grpc"
)

func main() {
	service.InitRedisClient()
	// Define the gRPC server port
	const port = ":50051"

	// Create a new gRPC server instance
	grpcServer := grpc.NewServer()

	// Create an instance of the ImageProcessingServer
	imageProcessingServer := &service.ImageProcessingServer{}

	// Register the ImageProcessingServer with the gRPC server
	pb.RegisterImageProcessingServiceServer(grpcServer, imageProcessingServer)

	// Start listening on the specified port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	log.Printf("Server is listening on port %s", port)

	// Start the gRPC server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
