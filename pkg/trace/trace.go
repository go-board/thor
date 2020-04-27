package trace

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

func Initialize(srv string, typ string, param float64) {
	_, err := config.Configuration{
		ServiceName: srv,
		Reporter: &config.ReporterConfig{
			BufferFlushInterval: time.Second * 10,
			QueueSize:           4096,
		},
		Sampler: &config.SamplerConfig{
			Type:  typ,
			Param: param,
		},
	}.InitGlobalTracer(srv)
	if err != nil {
		fmt.Printf("Init tracer failed, %s\n", err)
	}
}

func StartSpan(ctx context.Context, name string, baggageItems map[string]string, tags map[string]interface{}) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, name)
	for k, v := range baggageItems {
		span.SetBaggageItem(k, v)
	}
	for k, v := range tags {
		span.SetTag(k, v)
	}
	return span, ctx
}
