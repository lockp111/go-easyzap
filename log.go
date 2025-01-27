package easyzap

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config easyzap config
type Config struct {
	LogPath       string        `json:"log_path"`
	ErrPath       string        `json:"err_path"`
	Level         zapcore.Level `json:"level"`
	MaxSize       int           `json:"max_size"`
	MaxBackups    int           `json:"max_backups"`
	MaxDay        int           `json:"max_day"`
	CallerSkip    int           `json:"caller_skip"`
	DisableStd    bool          `json:"disable_std"`
	JSONFormat    bool          `json:"json_format"`
	Trace         bool          `json:"trace"`
	Compress      bool          `json:"compress"`
	DisableCaller bool          `json:"disable_caller"`
}

// New
func New(cfg Config) *zap.Logger {
	var (
		cores []zapcore.Core
		ws    []zapcore.WriteSyncer
		opts  = []zap.Option{
			zap.WithPanicHook(zapcore.WriteThenPanic),
			zap.WithFatalHook(zapcore.WriteThenFatal),
		}
	)

	if !cfg.DisableStd {
		ws = append(ws, zapcore.AddSync(os.Stdout))
	}
	if !cfg.DisableCaller {
		opts = append(opts, zap.WithCaller(true), zap.AddCallerSkip(cfg.CallerSkip))
	}
	if cfg.Trace {
		opts = append(opts, zap.AddStacktrace(cfg.Level))
	}

	zapCfg := zap.NewProductionEncoderConfig()
	zapCfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	encoder := zapcore.NewConsoleEncoder(zapCfg)
	if cfg.JSONFormat {
		encoder = zapcore.NewJSONEncoder(zapCfg)
	}

	if cfg.LogPath != "" {
		ws = append(ws, zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.LogPath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxDay,
			Compress:   cfg.Compress,
			LocalTime:  true,
		}))
	}

	if cfg.ErrPath != "" {
		cores = append(cores, zapcore.NewCore(
			encoder,
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   cfg.ErrPath,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxDay,
				Compress:   cfg.Compress,
				LocalTime:  true,
			}),
			zap.NewAtomicLevelAt(zapcore.ErrorLevel),
		))
	}

	cores = append(cores, zapcore.NewCore(
		encoder,
		&zapcore.BufferedWriteSyncer{
			WS:            zapcore.NewMultiWriteSyncer(ws...),
			FlushInterval: time.Second,
		},
		zap.NewAtomicLevelAt(cfg.Level),
	))
	return zap.New(zapcore.NewTee(cores...), opts...)
}
