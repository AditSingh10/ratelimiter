package main

import (
	"context"
	"log"
	"net"

	"github.com/AditSingh10/ratelimiter/internal/limiter"
	pb "github.com/AditSingh10/ratelimiter/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedRateLimiterServer
	l limiter.Limiter
}

func (c *server) Allow(ctx context.Context, req *pb.AllowRequest) (*pb.AllowResponse, error) {
	allowed, remaining, resetAt, err := c.l.Allow(ctx, req.ClientId)
	if err != nil {
		return nil, err
	}
	return &pb.AllowResponse{
		Allowed:       allowed,
		Remaining:     remaining,
		ResetAtUnix:   resetAt.Unix(),
		AlgorithmUsed: "token_bucket",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	// add reflection grpcurl works
	reflection.Register(grpcServer)

	pb.RegisterRateLimiterServer(grpcServer, &server{l: limiter.NewTokenBucket(100, 10)})

	log.Println("server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
