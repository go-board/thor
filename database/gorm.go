package database

import (
	"fmt"
	"time"

	"github.com/go-board/x-go/xdatabase/xsql"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

type Callback interface {
	Name() string
	BeforeQuery(s *gorm.Scope)
	AfterQuery(s *gorm.Scope)
	BeforeRowQuery(s *gorm.Scope)
	AfterRowQuery(s *gorm.Scope)
	BeforeCreate(s *gorm.Scope)
	AfterCreate(s *gorm.Scope)
	BeforeUpdate(s *gorm.Scope)
	AfterUpdate(s *gorm.Scope)
	BeforeDelete(s *gorm.Scope)
	AfterDelete(s *gorm.Scope)
}

func ApplyCallback(db *gorm.DB, callback Callback) {
	db.Callback().Query().Before("gorm:query").Register(fmt.Sprintf("before_query_%s", callback.Name()), callback.Before)
	db.Callback().Query().After("gorm:query").Register(fmt.Sprintf("after_query_%s", callback.Name()), callback.After)

	db.Callback().RowQuery().Before("gorm:row_query").Register(fmt.Sprintf("before_rowquery_%s", callback.Name()), callback.Before)
	db.Callback().RowQuery().After("gorm:row_query").Register(fmt.Sprintf("after_rowquery_%s", callback.Name()), callback.After)

	db.Callback().Create().Before("gorm:create").Register(fmt.Sprintf("before_create_%s", callback.Name()), callback.Before)
	db.Callback().Create().After("gorm:create").Register(fmt.Sprintf("after_create_%s", callback.Name()), callback.After)

	db.Callback().Update().Before("gorm:update").Register(fmt.Sprintf("before_update_%s", callback.Name()), callback.Before)
	db.Callback().Update().After("gorm:update").Register(fmt.Sprintf("after_update_%s", callback.Name()), callback.After)

	db.Callback().Delete().Before("gorm:delete").Register(fmt.Sprintf("before_delete_%s", callback.Name()), callback.Before)
	db.Callback().Delete().After("gorm:delete").Register(fmt.Sprintf("after_delete_%s", callback.Name()), callback.After)
}

type loggerCallback struct {
	logger *zap.Logger
}

func NewLoggerCallback(l *zap.Logger, options xsql.ConnectionOptions) Callback {
	return &loggerCallback{
		logger: l.With(
			zap.String("dialect", options.DriverName),
			zap.String("dsn", options.Dsn),
		),
	}
}

func (l *loggerCallback) Name() string {
	return "logger"
}

func (l *loggerCallback) BeforeQuery(s *gorm.Scope) {
	l.before(s, "query")
}

func (l *loggerCallback) AfterQuery(s *gorm.Scope) {
	l.after(s, "query")
}

func (l *loggerCallback) BeforeRowQuery(s *gorm.Scope) {
	l.before(s, "row_query")
}

func (l *loggerCallback) AfterRowQuery(s *gorm.Scope) {
	l.after(s, "row_query")
}

func (l *loggerCallback) BeforeCreate(s *gorm.Scope) {
	l.before(s, "create")
}

func (l *loggerCallback) AfterCreate(s *gorm.Scope) {
	l.after(s, "create")
}

func (l *loggerCallback) BeforeUpdate(s *gorm.Scope) {
	l.before(s, "update")
}

func (l *loggerCallback) AfterUpdate(s *gorm.Scope) {
	l.after(s, "update")
}

func (l *loggerCallback) BeforeDelete(s *gorm.Scope) {
	l.before(s, "delete")
}

func (l *loggerCallback) AfterDelete(s *gorm.Scope) {
	l.after(s, "delete")
}

func (l *loggerCallback) before(s *gorm.Scope, sqlType string) {}

func (l *loggerCallback) after(s *gorm.Scope, sqlType string) {
	var err error
	if s.HasError() {
		err = s.DB().Error
	}
	if err != nil {
		l.logger.With(
			zap.String("sql_type", sqlType),
			zap.String("instance_id", s.InstanceID()),
			zap.String("sql", s.SQL),
			zap.String("table_name", s.TableName()),
			zap.Error(err),
		).Error("execute sql failed")
	} else {
		l.logger.With(
			zap.String("sql_type", sqlType),
			zap.String("instance_id", s.InstanceID()),
			zap.String("sql", s.SQL),
			zap.String("table_name", s.TableName()),
		).Info("")
	}
}

type metricCallback struct {
	histogram *prometheus.HistogramVec
	timerKey  string
}

func NewMetricCallback(options xsql.ConnectionOptions, timerKey string) Callback {
	if timerKey == "" {
		timerKey = "gorm_timer_key"
	}
	return &metricCallback{
		timerKey: timerKey,
		histogram: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "gorm_operation_seconds",
				Help:        "database histogram",
				ConstLabels: prometheus.Labels{"dialect": options.DriverName, "dsn": options.Dsn},
			},
			[]string{"sql_type", "instance_id"},
		),
	}
}

func (m *metricCallback) Name() string {
	return "metric"
}

func (m *metricCallback) BeforeQuery(s *gorm.Scope) {
	m.before(s, "query")
}

func (m *metricCallback) AfterQuery(s *gorm.Scope) {
	m.after(s, "query")
}

func (m *metricCallback) BeforeRowQuery(s *gorm.Scope) {
	m.before(s, "row_query")
}

func (m *metricCallback) AfterRowQuery(s *gorm.Scope) {
	m.after(s, "row_query")
}

func (m *metricCallback) BeforeCreate(s *gorm.Scope) {
	m.before(s, "create")
}

func (m *metricCallback) AfterCreate(s *gorm.Scope) {
	m.after(s, "create")
}

func (m *metricCallback) BeforeUpdate(s *gorm.Scope) {
	m.before(s, "update")
}

func (m *metricCallback) AfterUpdate(s *gorm.Scope) {
	m.after(s, "update")
}

func (m *metricCallback) BeforeDelete(s *gorm.Scope) {
	m.before(s, "delete")
}

func (m *metricCallback) AfterDelete(s *gorm.Scope) {
	m.after(s, "delete")
}

func (m *metricCallback) before(s *gorm.Scope, sqlType string) {
	s.Set(m.timerKey, time.Now())
}

func (m *metricCallback) after(s *gorm.Scope, sqlType string) {
	val, ok := s.Get(m.timerKey)
	if ok {
		startTime := val.(time.Time)
		m.histogram.WithLabelValues(sqlType, s.InstanceID()).Observe(time.Since(startTime).Seconds())
	}
}
