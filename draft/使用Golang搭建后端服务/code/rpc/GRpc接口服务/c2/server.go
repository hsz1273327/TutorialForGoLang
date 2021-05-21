package main

import (
	"io"
	"log"
	"net"

	pb "c2/squarerpc_service"

	grpc "google.golang.org/grpc"
)

const (
	// server port 5000
	ADDRESS = "0.0.0.0:5000"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SumSquare(instream pb.SquareService_SumSquareServer) error {
	var sum float64 = 0.0
	for {
		data, err := instream.Recv()
		if err != nil {
			if err == io.EOF {
				return instream.SendAndClose(&pb.Message{Message: sum})
			}
			if err != nil {
				return err
			}
		} else {
			sum += data.Message * data.Message
		}
	}
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
