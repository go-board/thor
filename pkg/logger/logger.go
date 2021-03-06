package logger

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	once        sync.Once
)

func ServeLevel(w http.ResponseWriter, r *http.Request) { globalLevel.ServeHTTP(w, r) }

// Initialize do some initial work for log.
// Info logger will write to <dir>/<namespace>-<service_name>-<service_id>-info.log
// Error logger will write to <dir>/<namespace>-<service_name>-<service_id>-error.log
// Each log file size is limited to 1GB,
// and maximum days of old log to retained in day, the default is 14 days(2 week),
// currently both file size and max ages not support change by user.
func Initialize(dir string, namespace, serviceName, serviceId string) {
	once.Do(func() {
		cfg := zapcore.EncoderConfig{
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		}

		infoOutput := &lumberjack.Logger{
			Filename:  fmt.Sprintf("%s/%s-%s-%s-info.log", dir, namespace, serviceName, serviceId),
			LocalTime: true,
			MaxSize:   1024,
			MaxAge:    14,
		}
		errorOutput := &lumberjack.Logger{
			Filename:  fmt.Sprintf("%s/%s-%s-%s-error.log", dir, namespace, serviceName, serviceId),
			LocalTime: true,
			MaxSize:   1024,
			MaxAge:    14,
		}

		logger := zap.New(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(cfg),
				zapcore.NewMultiWriteSyncer(zapcore.AddSync(infoOutput)),
				globalLevel,
			),
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.PanicLevel),
			zap.ErrorOutput(zapcore.AddSync(errorOutput)),
			zap.Fields(
				zap.String("namespace", namespace),
				zap.String("service_name", serviceName),
				zap.String("service_id", serviceId),
			),
		)
		zap.ReplaceGlobals(logger)
	})
}
