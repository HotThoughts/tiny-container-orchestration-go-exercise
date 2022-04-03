package logger

import "go.uber.org/zap"

func New() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	sugarLogger := logger.Sugar()
	return sugarLogger
}
