package logger

import (
	"fmt"
	"net/http"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)

func ServeLogLevel(w http.ResponseWriter, r *http.Request) { globalLevel.ServeHTTP(w, r) }

func Initialize(dir string, namespace, serviceName, serviceId string) {
	cfg := zapcore.EncoderConfig{
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
	}

	infoLogger = initLogger(fmt.Sprintf("%s/%-%s-%s-info.log", dir, namespace, serviceName, serviceId), cfg, globalLevel)
	zap.ReplaceGlobals(infoLogger)
	errorLogger = initLogger(fmt.Sprintf("%s/%-%s-%s-error.log", dir, namespace, serviceName, serviceId), cfg, zapcore.ErrorLevel)
}

func initLogger(file string, cfg zapcore.EncoderConfig, lvl zapcore.LevelEnabler) *zap.Logger {
	l := &lumberjack.Logger{
		Filename:  file,
		LocalTime: true,
		MaxSize:   1024,
		MaxAge:    14,
	}
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(l)),
		lvl,
	))
	return logger
}

var errorLogger = zap.L()
var infoLogger = zap.L()

func Error() *zap.Logger { return errorLogger }

func Info() *zap.Logger { return infoLogger }
