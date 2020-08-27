package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//ZSLogger Zap Sugared Logger
var ZSLogger *zap.SugaredLogger
var ZLogger *zap.Logger

func init()  {
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, _ := cfg.Build()

	ZSLogger = logger.Sugar().With()
	ZLogger = logger
}
