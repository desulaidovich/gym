package logger

import (
	"fmt"
	"time"

	"github.com/desulaidovich/gym/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg *config.Application) (*zap.Logger, error) {
	zapCfg := zap.NewProductionConfig()
	if cfg.Development {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	level, err := zap.ParseAtomicLevel(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("log level %q: %w", cfg.LogLevel, err)
	}

	zapCfg.Level = level
	zapCfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	zapCfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	zapCfg.EncoderConfig.NameKey = "component"
	zapCfg.EncoderConfig.TimeKey = "time"

	log, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}

	return log, nil
}
