package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	pb "c0/squarerpc_service"

	grpc "google.golang.org/grpc"
)

const (
	//server address
	ADDRESS = "localhost:5000"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// 创建client对象
	c := pb.NewSquareServiceClient(conn)

	// Contact the server and print out its response.
	query := 2.0
	if len(os.Args) > 1 {
		query, err = strconv.ParseFloat(os.Args[1], 64)
		if err != nil {
			log.Fatalf("query can not parse from string: %v", err)
			panic(err)
		}
	}
	// 设置请求上下文的过期时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 调用rpc的函数
	r, err := c.Square(ctx, &pb.Message{Message: query})
	if err != nil {
		log.Fatalf("could not call: %v", err)
	}
	log.Printf("Result: %f", r.Message)
}
