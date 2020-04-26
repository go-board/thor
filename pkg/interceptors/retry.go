package interceptors

import (
	"context"
	"errors"

	"github.com/go-board/x-go/xslice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type RetryOption func(r *GrpcRetry)

func RetryCodes(codes ...int) RetryOption {
	return func(r *GrpcRetry) {
		r.retryCodes = codes
	}
}

type GrpcRetry struct {
	retryCodes  []int
	retryErrors []error
	maxRetries  int
}

func NewGrpcRetry(options ...RetryOption) *GrpcRetry {
	retry := &GrpcRetry{}
	for _, option := range options {
		option(retry)
	}
	return retry
}

func (g *GrpcRetry) UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if !isIdempotent(ctx) {
		select {
		default:
		case <-ctx.Done():
			return ctx.Err()
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	var err error
	for i := 0; i < g.maxRetries; i++ {
		select {
		default:
		case <-ctx.Done():
			return ctx.Err()
		}
		err = invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			if sts, ok := status.FromError(err); ok {
				if xslice.ContainsInt(g.retryCodes, int(sts.Code())) {
					continue
				}
			}
			for _, checkErr := range g.retryErrors {
				if errors.Is(err, checkErr) {
					continue
				}
			}
		}
	}
	return err
}

func (g *GrpcRetry) StreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	if !isIdempotent(ctx) {
		return streamer(ctx, desc, cc, method, opts...)
	}
	// todo: handle stream client retry
	return streamer(ctx, desc, cc, method, opts...)
}

func isIdempotent(ctx context.Context) bool {
	isIdempotent := false
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, ok := md["is_idempotent"]; ok {
			if len(val) > 0 && val[0] == "1" {
				isIdempotent = true
			}
		}
	}
	return isIdempotent
}
