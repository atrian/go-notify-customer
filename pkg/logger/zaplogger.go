package logger

import (
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"strings"

	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() *ZapLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	return &ZapLogger{
		logger: logger,
	}
}

func NewFatalZapLogger() *ZapLogger {
	errorFatalLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.FatalLevel
	})

	// write syncers
	stderrSyncer := zapcore.Lock(os.Stderr)

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			stderrSyncer,
			errorFatalLevel,
		),
	)

	logger := zap.New(core)

	return &ZapLogger{
		logger: logger,
	}
}

func (z ZapLogger) Fatal(message string, err error) {
	z.logger.Fatal(message, zap.Error(err))
}

func (z ZapLogger) Panic(message string, err error) {
	z.logger.Panic(message, zap.Error(err))
}

func (z ZapLogger) Error(message string, err error) {
	z.logger.Error(message, zap.Error(err))
}

func (z ZapLogger) Warning(message string) {
	z.logger.Warn(message)
}

func (z ZapLogger) Info(message string) {
	z.logger.Info(message)
}

func (z ZapLogger) Debug(message ...string) {
	z.logger.Debug(strings.Join(message, " "))
}

func (z ZapLogger) Sync() {
	err := z.logger.Sync()
	if err != nil {
		log.Print(err)
	}
}
