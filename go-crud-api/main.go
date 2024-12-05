package main

import (
	"context"
	"example.com/go-crud-api/db"
	pb "example.com/go-crud-api/go-crud-api"
	"example.com/go-crud-api/router"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	return &pb.UserResponse{
		Name: "Test User",
		Age:  30,
	}, nil
}

func main() {
	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		// Create a new gRPC server
		grpcServer := grpc.NewServer()
		pb.RegisterUserServiceServer(grpcServer, &server{})
		log.Println("gRPC server is running on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()
	// Start HTTP server
	go func() {
		db.InitPostgresDB()
		router.InitRouter().Run()
	}()
	// Block main thread to keep servers running
	select {}

}
