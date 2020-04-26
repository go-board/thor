package server

import (
	"context"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/go-board/thor/pkg/registry"
)

type Server struct {
	srv      *grpc.Server
	registry registry.Registry
}

func NewServer() *grpc.Server {
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(func(ctx context.Context, p interface{}) (err error) {
				return
			})),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(zap.L()),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
			grpc_opentracing.StreamServerInterceptor(),
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(zap.L()),
		)),
		grpc.StatsHandler(&statsHandler{}),
		grpc.UnknownServiceHandler(unknownServiceHandler),
		grpc.InTapHandle(inTapHandle),
	)
	return srv
}

func (s *Server) RegisterService(sd *grpc.ServiceDesc, srv interface{}) {
	s.srv.RegisterService(sd, srv)
}

func (s *Server) Start(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// todo: register service
	return s.srv.Serve(ln)
}

func (s *Server) Close() {
	s.srv.GracefulStop()
}
