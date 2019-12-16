package server

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func serverUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("invoke server method=%s duration=%s error=%v", info.FullMethod, time.Since(start), err)
	return resp, err
}

func serverStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	start := time.Now()
	err := handler(srv, ss)
	log.Printf("invoke server method=%s duration=%s error=%v", info.FullMethod, time.Since(start), err)
	return err
}

func createServerUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(serverUnaryInterceptor)
}

func createServerStreamInterceptor() grpc.ServerOption {
	return grpc.StreamInterceptor(serverStreamInterceptor)
}
