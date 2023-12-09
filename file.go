package easyzap

import (
	"compress/gzip"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

var patternConversionRegexps = []*regexp.Regexp{
	regexp.MustCompile(`%[%+A-Za-z]`),
	regexp.MustCompile(`\*+`),
}

type logFile struct {
	*slog.Logger
	pattern  string
	maxAge   time.Duration
	maxCount uint
}

func newLogFile(path string, rotationCount uint, rotationSeconds, maxDay int64) io.Writer {
	var (
		logPath     = strings.TrimSuffix(path, ".log") + ".%Y%m%d%H%M.log"
		filePattern = logPath + ".gz"
	)
	for _, re := range patternConversionRegexps {
		filePattern = re.ReplaceAllString(filePattern, "*")
	}

	// 生成rotatelogs的Logger 实际生成的文件名 demo.YYmmddHH.log
	// 保存7天内的日志，每小时分割一次日志
	lf := &logFile{
		pattern:  filePattern,
		maxAge:   time.Hour * 24 * time.Duration(maxDay),
		maxCount: rotationCount,
	}
	writer, err := rotatelogs.New(
		logPath,
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(maxDay)),
		rotatelogs.WithRotationCount(rotationCount),
		rotatelogs.WithRotationTime(time.Second*time.Duration(rotationSeconds)),
		rotatelogs.WithHandler(lf),
	)
	if err != nil {
		panic(err)
	}
	lf.Logger = slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{
		Level:     slog.LevelError,
		AddSource: true,
	}))
	return writer
}

func (l *logFile) Handle(e rotatelogs.Event) {
	if e.Type() != rotatelogs.FileRotatedEventType {
		return
	}
	// 切割完成，获取上一个文件名
	prevFilePath := e.(*rotatelogs.FileRotatedEvent).PreviousFile()
	if prevFilePath == "" {
		return
	}

	// 压缩
	paths, fileName := filepath.Split(prevFilePath)
	prevFile, err := os.Open(prevFilePath)
	if err != nil {
		slog.Error("logFile open file fail", err,
			"prevFilePath", prevFilePath,
		)
		return
	}
	defer prevFile.Close()

	gzipFile, err := os.Create(paths + fileName + ".gz")
	if err != nil {
		slog.Error("logFile create gz fail", err,
			"prevFilePath", prevFilePath,
		)
		return
	}
	defer gzipFile.Close()

	gw := gzip.NewWriter(gzipFile)
	defer gw.Close()

	prevInfo, err := prevFile.Stat()
	if err != nil {
		slog.Error("logFile get prevFile stat fail", err,
			"prevFilePath", prevFilePath,
		)
		return
	}
	gw.Header.Name = prevInfo.Name()

	_, err = io.Copy(gw, prevFile)
	if err != nil {
		slog.Error("logFile copy file fail", err,
			"prevFilePath", prevFilePath,
		)
		return
	}
	// 删除原文件
	os.RemoveAll(prevFilePath)
	// 删除过期压缩日志
	if err := l.removeGzip(l.pattern); err != nil {
		slog.Error("logFile remove gz fail", err,
			"prevFilePath", prevFilePath,
		)
	}
}

func (l *logFile) removeGzip(filePattern string) error {
	matches, err := filepath.Glob(filePattern)
	if err != nil {
		return err
	}

	cutoff := time.Now().Add(-1 * l.maxAge)
	var toUnlink []string
	for _, path := range matches {
		// Ignore lock files
		if strings.HasSuffix(path, "_lock") || strings.HasSuffix(path, "_symlink") {
			continue
		}

		fi, err := os.Stat(path)
		if err != nil {
			continue
		}

		fl, err := os.Lstat(path)
		if err != nil {
			continue
		}

		if l.maxAge > 0 && fi.ModTime().After(cutoff) {
			continue
		}

		if l.maxCount > 0 && fl.Mode()&os.ModeSymlink == os.ModeSymlink {
			continue
		}
		toUnlink = append(toUnlink, path)
	}

	if l.maxCount > 0 {
		// Only delete if we have more than rotationCount
		if l.maxCount >= uint(len(toUnlink)) {
			return nil
		}

		toUnlink = toUnlink[:len(toUnlink)-int(l.maxCount)]
	}

	if len(toUnlink) <= 0 {
		return nil
	}

	go func() {
		// unlink files on a separate goroutine
		for _, path := range toUnlink {
			os.Remove(path)
		}
	}()

	return nil
}
