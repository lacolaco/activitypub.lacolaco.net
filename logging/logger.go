package logging

import (
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"go.ajitem.com/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(conf *config.Config) *zap.Logger {
	var zc zap.Config
	if conf.IsRunningOnCloud() {
		zc = zapdriver.NewProductionConfig()
		zc.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		zc.EncoderConfig.TimeKey = zapcore.OmitKey
	} else {
		zc = zapdriver.NewDevelopmentConfig()
	}
	logger, _ := zc.Build()
	return logger
}
