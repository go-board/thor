package trace

import (
	"context"
	"log"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

var (
	httpTag = opentracing.Tag{Key: string(ext.Component), Value: "HTTP"}
)

func IncomingHTTPRequest(r *http.Request) (*http.Request, opentracing.Span) {
	ctx, span := newSpanFromHTTPRequest(opentracing.GlobalTracer(), r)
	r = r.WithContext(ctx)
	return r, span
}

func OutgoingHTTPRequest(span opentracing.Span, r *http.Request) error {
	return span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
}

func newSpanFromHTTPRequest(tracer opentracing.Tracer, r *http.Request) (context.Context, opentracing.Span) {
	parentSpanContext, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil && err != opentracing.ErrSpanContextNotFound {
		log.Printf("grpc_opentracing: failed parsing trace information: %v", err)
	}
	serverSpan := tracer.StartSpan(
		r.URL.Path,
		// this is magical, it attaches the new span to the parent parentSpanContext, and creates an unparented one if empty.
		ext.RPCServerOption(parentSpanContext),
		httpTag,
	)
	return opentracing.ContextWithSpan(r.Context(), serverSpan), serverSpan
}
