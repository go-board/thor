package client

import (
	"context"

	"google.golang.org/grpc"
)

func NewClient(ctx context.Context, target string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		ctx,
		target,
		grpc.WithChainUnaryInterceptor(),
		grpc.WithChainStreamInterceptor(),
	)
}
