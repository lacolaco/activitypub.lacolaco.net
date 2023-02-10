package logger

import (
	"github.com/lacolaco/activitypub.lacolaco.net/config"
	"go.ajitem.com/zapdriver"
	"go.uber.org/zap"
)

func NewLogger(cfg *config.Config) *zap.Logger {
	var zc zap.Config
	zc = zapdriver.NewDevelopmentConfig()
	// if cfg.Env == env.EnvDevelopment {
	// } else if cfg.Env == env.EnvStaging {
	// 	zc = zapdriver.NewProductionConfig()
	// 	zc.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	// 	zc.EncoderConfig.TimeKey = zapcore.OmitKey
	// } else {
	// 	zc = zapdriver.NewProductionConfig()
	// 	zc.EncoderConfig.TimeKey = zapcore.OmitKey
	// }
	logger, _ := zc.Build()
	return logger
}
