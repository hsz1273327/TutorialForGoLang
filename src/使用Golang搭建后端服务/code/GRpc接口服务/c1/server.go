package main

import (
	"log"
	"net"

	pb "c1/squarerpc_service"

	grpc "google.golang.org/grpc"
)

const (
	// server port 5000
	ADDRESS = "0.0.0.0:5000"
)

type server struct{}

// 第一位为请求,第二位为流对象,用于发送回消息
func (s *server) RangeSquare(in *pb.Message, stream pb.SquareService_RangeSquareServer) error {
	limit := int(in.Message)
	for i := 0; i <= limit; i++ {
		err := stream.Send(&pb.Message{Message: float64(i * i)})
		if err != nil {
			return err
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
