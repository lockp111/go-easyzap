# go-easylog
easy package by zap

## Usage

### Quick start
```golang
easyzap.Info("info")
easyzap.Warn("warn")
easyzap.Error("error")
```

### Inject
```golang
easyzap.Inject(easyzap.New(&Config{
		LogDir:   "logs/",
		Filename: "out.log",
		Level:    zapcore.InfoLevel,
    }))

easyzap.Info("info")
easyzap.Warn("warn")
easyzap.Error("error")
```

### New
```golang
zlog := New(&Config{
		LogDir:   "logs/",
		Filename: "out.log",
		Level:    zapcore.InfoLevel,
	})

zlog.Info("info")
zlog.Warn("warn")
zlog.Error("error")
```
