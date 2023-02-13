package logging

import (
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"go.ajitem.com/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(conf *config.Config) *zap.Logger {
	zc := zapdriver.NewDevelopmentConfig()
	if conf.IsRunningOnCloud() {
		zc.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		zc.EncoderConfig.TimeKey = zapcore.OmitKey
	}
	logger, _ := zc.Build()
	return logger
}
