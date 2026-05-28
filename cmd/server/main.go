package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/AditSingh10/ratelimiter/internal/limiter"
	pb "github.com/AditSingh10/ratelimiter/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedRateLimiterServer
	limiters map[string]limiter.Limiter
}

func (c *server) Allow(ctx context.Context, req *pb.AllowRequest) (*pb.AllowResponse, error) {
	l, ok := c.limiters[req.Algorithm]
	if !ok {
		return nil, fmt.Errorf("unknown algorithm: %s", req.Algorithm)

	}
	allowed, remaining, resetAt, err := l.Allow(ctx, req.ClientId)
	if err != nil {
		return nil, err
	}
	return &pb.AllowResponse{
		Allowed:       allowed,
		Remaining:     remaining,
		ResetAtUnix:   resetAt.Unix(),
		AlgorithmUsed: req.Algorithm,
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

	pb.RegisterRateLimiterServer(grpcServer, &server{limiters: map[string]limiter.Limiter{
		"token_bucket": limiter.NewTokenBucket(100, 10),
		"leaky_bucket": limiter.NewLeakyBucket(10, 10*time.Second)}})

	log.Println("server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
