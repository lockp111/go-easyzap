package easyzap

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestLogNormal(t *testing.T) {
	zlog := New(Config{
		LogPath: "logs/out.log",
		Level:   zapcore.InfoLevel,
	})

	zlog.Info("TestLogNormal")
	zlog.Warn("TestLogNormal")
	zlog.Error("TestLogNormal")
	zlog.Sync()
}

func TestLogDisableStd(t *testing.T) {
	zlog := New(Config{
		LogPath:    "logs/out.log",
		Level:      zapcore.InfoLevel,
		DisableStd: true,
	})

	zlog.Info("TestLogDisableStd")
	zlog.Warn("TestLogDisableStd")
	zlog.Error("TestLogDisableStd")
	zlog.Sync()
}

func TestLogJSONOut(t *testing.T) {
	zlog := New(Config{
		LogPath:    "logs/out.log",
		Level:      zapcore.InfoLevel,
		JSONFormat: true,
		DisableStd: false,
	})

	zlog.Info("TestLogJSONOut")
	zlog.Warn("TestLogJSONOut")
	zlog.Error("TestLogJSONOut")
	zlog.Sync()
}

func TestLogErrOut(t *testing.T) {
	zlog := New(Config{
		LogPath:    "logs/out.log",
		ErrPath:    "logs/err.log",
		Level:      zapcore.InfoLevel,
		JSONFormat: true,
	})

	zlog.Info("TestLogErrOut")
	zlog.Warn("TestLogErrOut")
	zlog.Error("TestLogErrOut")
	zlog.Sync()
}

func TestLogTrace(t *testing.T) {
	zlog := New(Config{
		LogPath: "logs/out.log",
		ErrPath: "logs/err.log",
		Level:   zapcore.InfoLevel,
		Trace:   true,
	})

	zlog.Info("TestLogTrace")
	zlog.Warn("TestLogTrace")
	zlog.Error("TestLogTrace")
	zlog.Sync()
}

func TestLogTraceJSON(t *testing.T) {
	zlog := New(Config{
		LogPath:    "logs/out.log",
		ErrPath:    "logs/err.log",
		Level:      zapcore.InfoLevel,
		Trace:      true,
		JSONFormat: true,
	})

	zlog.Info("TestLogTrace")
	zlog.Warn("TestLogTrace")
	zlog.Error("TestLogTrace")
	zlog.Sync()
}

func TestAPI(t *testing.T) {
	Info("TestLogAPI")
	Warn("TestLogAPI")
	Error("TestLogAPI")
	Sync()
}
