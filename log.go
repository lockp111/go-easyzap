package easyzap

import (
	"io"
	"os"
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
func New(cfg *Config) *zap.SugaredLogger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		TimeKey:     "ts",
		CallerKey:   "caller",
		LineEnding:  zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		//EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
	}

	if cfg.Trace {
		encoderCfg.StacktraceKey = "trace"
	}

	var (
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
		cores   []zapcore.Core
	)
	if !cfg.DisableStd {
		cores = append(cores, zapcore.NewCore(
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
		cores = append(cores, zapcore.NewCore(
			encoder,
			zapcore.AddSync(fileOut),
			cfg.Level,
		))
	}

	if len(cfg.ErrPath) != 0 {
		errOut := getWriter(cfg.ErrPath)
		cores = append(cores, zapcore.NewCore(
			encoder,
			zapcore.AddSync(errOut),
			zap.ErrorLevel,
		))
	}

	// 需要传入zap.AddCaller()才会显示打日志点的文件名和行数
	logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	defer logger.Sync()
	return logger.Sugar()
}

// getWriter 获取日志文件的io.Writer
func getWriter(path string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// 保存7天内的日志，每天分割一次日志
	hook, err := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		panic(err)
	}
	return hook
}
