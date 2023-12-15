package easyzap

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config easyzap config
type Config struct {
	LogPath         string
	ErrPath         string
	RotationCount   uint
	RotationSeconds int64
	MaxDay          int64
	Level           zapcore.Level
	JSONFormat      bool
	Trace           bool
	DisableStd      bool
}

// New
func New(cfg *Config, cores ...zapcore.Core) *zap.SugaredLogger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		CallerKey:      "caller",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
	}

	if cfg.Trace {
		encoderCfg.StacktraceKey = "trace"
	}

	var newCores []zapcore.Core
	if !cfg.DisableStd {
		newCores = append(newCores, zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderCfg),
			zapcore.AddSync(os.Stdout),
			cfg.Level,
		))
	}

	// 去掉颜色
	encoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	var encoder = zapcore.NewConsoleEncoder(encoderCfg)
	if cfg.JSONFormat {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	if len(cfg.LogPath) != 0 {
		fileOut := newLogFile(cfg.LogPath, cfg.RotationCount, cfg.RotationSeconds, cfg.MaxDay)
		newCores = append(newCores, zapcore.NewCore(
			encoder,
			zapcore.AddSync(fileOut),
			cfg.Level,
		))
	}

	if len(cfg.ErrPath) != 0 {
		errOut := newLogFile(cfg.ErrPath, cfg.RotationCount, cfg.RotationSeconds, cfg.MaxDay)
		newCores = append(newCores, zapcore.NewCore(
			encoder,
			zapcore.AddSync(errOut),
			zap.ErrorLevel,
		))
	}

	newCores = append(newCores, cores...)
	// 需要传入zap.AddCaller()才会显示打日志点的文件名和行数
	logger := zap.New(zapcore.NewTee(newCores...),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	return logger.Sugar()
}
