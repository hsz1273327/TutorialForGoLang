package main

import (
	"container/list"
	"io"
	"log"
	"net"

	pb "c3/squarerpc_service"

	grpc "google.golang.org/grpc"
)

const (
	// server port 5000
	ADDRESS = "0.0.0.0:5000"
)

type server struct{}

func (s *server) StreamrangeSquare(stream pb.SquareService_StreamrangeSquareServer) error {
	l := list.New()
	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		} else {
			l.PushBack(in.Message * in.Message)
		}
	}
	for e := l.Front(); e != nil; e = e.Next() {
		stream.Send(&pb.Message{Message: e.Value.(float64)})
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
