package cache

import (
	"context"
	"time"

	"github.com/go-board/x-go/xctx"
	"github.com/go-redis/redis/v7"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

type loggerMiddleware struct {
	logger *zap.Logger
}

func NewLoggerMiddleware(l *zap.Logger) redis.Hook {
	return &loggerMiddleware{logger: l}
}

func (l *loggerMiddleware) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (l *loggerMiddleware) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	err := cmd.Err()
	if err != nil {
		l.logger.With().Error("")
	} else {
		l.logger.With().Info("")
	}
	return nil
}

func (l *loggerMiddleware) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (l *loggerMiddleware) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}

type metricMiddleware struct {
	histogram *prometheus.HistogramVec
	timeKey   interface{}
}

func NewMetricMiddleware() redis.Hook {
	return &metricMiddleware{
		histogram: promauto.NewHistogramVec(prometheus.HistogramOpts{}, []string{"cmd", "status"}),
		timeKey:   "redis-timer-key",
	}
}

func (m *metricMiddleware) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	ctx = context.WithValue(ctx, m.timeKey, time.Now())
	return ctx, nil
}

func (m *metricMiddleware) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if start, ok := xctx.ReadTime(ctx, m.timeKey); ok {
		m.histogram.WithLabelValues(cmd.Name(), status(cmd.Err())).Observe(time.Since(start).Seconds())
	}
	return nil
}

func (m *metricMiddleware) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	ctx = context.WithValue(ctx, m.timeKey, time.Now())
	return ctx, nil
}

func (m *metricMiddleware) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if start, ok := xctx.ReadTime(ctx, m.timeKey); ok {
		m.histogram.WithLabelValues("pipeline", "success").Observe(time.Since(start).Seconds())
	}
	return nil
}

func status(err error) string {
	if err != nil {
		return "failure"
	}
	return "success"
}
