package main

import (
	"io"
	"log"
	"net"

	pb "c4/squarerpc_service"

	grpc "google.golang.org/grpc"
)

const (
	// server port 5000
	ADDRESS = "0.0.0.0:5000"
)

type server struct{}

func (s *server) StreamrangeSquare(stream pb.SquareService_StreamrangeSquareServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		} else {
			err := stream.Send(&pb.Message{Message: in.Message * in.Message})
			if err != nil {
				return err
			}
		}
	}
	return nil
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
