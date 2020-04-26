package server

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/tap"
)

type statsHandler struct {
}

func (s *statsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return ctx
}

func (s *statsHandler) HandleRPC(ctx context.Context, stats stats.RPCStats) {

}

func (s *statsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return ctx
}

func (s *statsHandler) HandleConn(ctx context.Context, stats stats.ConnStats) {
}

func unknownServiceHandler(srv interface{}, stream grpc.ServerStream) error {
	return nil
}

func inTapHandle(ctx context.Context, info *tap.Info) (context.Context, error) {
	return ctx, nil
}
