package easyzap

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestLogNormal(t *testing.T) {
	zlog := New(&Config{
		LogPath: "logs/out.log",
		Level:   zapcore.InfoLevel,
	})

	zlog.Info("TestLogNormal")
	zlog.Warn("TestLogNormal")
	zlog.Error("TestLogNormal")
}

func TestLogDisableStd(t *testing.T) {
	zlog := New(&Config{
		LogPath:    "logs/out.log",
		Level:      zapcore.InfoLevel,
		DisableStd: true,
	})

	zlog.Info("TestLogDisableStd")
	zlog.Warn("TestLogDisableStd")
	zlog.Error("TestLogDisableStd")
}

func TestLogJSONOut(t *testing.T) {
	zlog := New(&Config{
		LogPath:    "logs/out.log",
		Level:      zapcore.InfoLevel,
		JSONFormat: true,
		DisableStd: true,
	})

	zlog.Info("TestLogJSONOut")
	zlog.Warn("TestLogJSONOut")
	zlog.Error("TestLogJSONOut")
}

func TestLogErrOut(t *testing.T) {
	zlog := New(&Config{
		LogPath:    "logs/out.log",
		ErrPath:    "logs/err.log",
		Level:      zapcore.InfoLevel,
		JSONFormat: true,
	})

	zlog.Info("TestLogErrOut")
	zlog.Warn("TestLogErrOut")
	zlog.Error("TestLogErrOut")
}

func TestLogTrace(t *testing.T) {
	zlog := New(&Config{
		LogPath: "logs/out.log",
		ErrPath: "logs/err.log",
		Level:   zapcore.InfoLevel,
		Trace:   true,
	})

	zlog.Info("TestLogTrace")
	zlog.Warn("TestLogTrace")
	zlog.Error("TestLogTrace")
}

func TestLogTraceJSON(t *testing.T) {
	zlog := New(&Config{
		LogPath:    "logs/out.log",
		ErrPath:    "logs/err.log",
		Level:      zapcore.InfoLevel,
		Trace:      true,
		JSONFormat: true,
	})

	zlog.Info("TestLogTrace")
	zlog.Warn("TestLogTrace")
	zlog.Error("TestLogTrace")
}
