package main

import (
	"context"
	"log"
	"net"

	pb "github.com/AditSingh10/ratelimiter/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedRateLimiterServer
}

func (c *server) Allow(context.Context, *pb.AllowRequest) (*pb.AllowResponse, error) {
	out := new(pb.AllowResponse)
	out.Allowed = true
	out.Remaining = 99
	return out, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRateLimiterServer(grpcServer, &server{})

	log.Println("server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
