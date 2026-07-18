package logging

import (
	"clash-admin/internal/config"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const defaultFilePath = "runtime/logs/clash-admin.log"

var (
	mu       sync.Mutex
	enabled  bool
	logger   *slog.Logger
	logFile  *os.File
	minLevel slog.Level
)

// Init 根据配置初始化 JSON 文件日志；未启用时直接返回。
func Init(loggingConfig config.Logging) error {
	mu.Lock()
	defer mu.Unlock()

	resetLocked()

	if !loggingConfig.Enabled {
		return nil
	}

	filePath := strings.TrimSpace(loggingConfig.FilePath)
	if filePath == "" {
		filePath = defaultFilePath
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return fmt.Errorf("create logging directory: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open logging file %s: %w", filePath, err)
	}

	minLevel = parseLevel(loggingConfig.Level)
	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: minLevel,
	})
	logger = slog.New(handler)
	logFile = file
	enabled = true
	return nil
}

// Enabled 返回文件日志是否已启用。
func Enabled() bool {
	mu.Lock()
	defer mu.Unlock()
	return enabled
}

// Info 写入 INFO 级别 JSON 日志；未启用时为 no-op。
func Info(message string, attrs ...any) {
	write(slog.LevelInfo, message, attrs...)
}

// Warn 写入 WARN 级别 JSON 日志；未启用时为 no-op。
func Warn(message string, attrs ...any) {
	write(slog.LevelWarn, message, attrs...)
}

// Debug 写入 DEBUG 级别 JSON 日志；未启用时为 no-op。
func Debug(message string, attrs ...any) {
	write(slog.LevelDebug, message, attrs...)
}

// Error 写入 ERROR 级别 JSON 日志；未启用时为 no-op。
func Error(message string, attrs ...any) {
	write(slog.LevelError, message, attrs...)
}

// Close 关闭日志文件句柄。
func Close() error {
	mu.Lock()
	defer mu.Unlock()
	return resetLocked()
}

func write(level slog.Level, message string, attrs ...any) {
	mu.Lock()
	defer mu.Unlock()
	if !enabled || logger == nil {
		return
	}
	if level < minLevel {
		return
	}
	switch level {
	case slog.LevelWarn:
		logger.Warn(message, attrs...)
	case slog.LevelError:
		logger.Error(message, attrs...)
	case slog.LevelDebug:
		logger.Debug(message, attrs...)
	default:
		logger.Info(message, attrs...)
	}
}

func resetLocked() error {
	enabled = false
	logger = nil
	minLevel = slog.LevelInfo

	if logFile == nil {
		return nil
	}

	var err error
	if closeErr := logFile.Close(); closeErr != nil {
		err = closeErr
	}
	logFile = nil
	return err
}

func parseLevel(raw string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
