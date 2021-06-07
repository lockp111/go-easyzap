package easyzap

import (
	"io"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config easyzap config
type Config struct {
	LogPath    string
	ErrPath    string
	Level      zapcore.Level
	JSONFormat bool
	Trace      bool
	DisableStd bool
}

// New
func New(cfg *Config, cores ...zapcore.Core) *zap.SugaredLogger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		TimeKey:     "ts",
		CallerKey:   "caller",
		LineEnding:  zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		// EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
	}

	if cfg.Trace {
		encoderCfg.StacktraceKey = "trace"
	}

	var (
		encoder  = zapcore.NewConsoleEncoder(encoderCfg)
		newCores []zapcore.Core
	)
	if !cfg.DisableStd {
		newCores = append(newCores, zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			cfg.Level,
		))
	}

	if cfg.JSONFormat {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	if len(cfg.LogPath) != 0 {
		fileOut := getWriter(cfg.LogPath)
		newCores = append(newCores, zapcore.NewCore(
			encoder,
			zapcore.AddSync(fileOut),
			cfg.Level,
		))
	}

	if len(cfg.ErrPath) != 0 {
		errOut := getWriter(cfg.ErrPath)
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
	defer logger.Sync()
	return logger.Sugar()
}

// getWriter 获取日志文件的io.Writer
func getWriter(path string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.YYmmdd.log
	// 保存7天内的日志，每天分割一次日志
	hook, err := rotatelogs.New(
		strings.TrimRight(path, ".log")+".%Y%m%d.log",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		panic(err)
	}
	return hook
}
