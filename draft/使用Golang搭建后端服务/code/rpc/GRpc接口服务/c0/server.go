package main

import (
	"context"
	"log"
	"net"

	pb "c0/squarerpc_service"

	grpc "google.golang.org/grpc"
)

const (
	// server port 5000
	ADDRESS = "0.0.0.0:5000"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) Square(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	log.Printf("Received: %v", in.Message)
	return &pb.Message{Message: in.Message * in.Message}, nil
}

func main() {
	lis, err := net.Listen("tcp", ADDRESS)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("server started @", ADDRESS)
	s := grpc.NewServer()
	pb.RegisterSquareServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
