package logging_test

import (
	"clash-admin/internal/config"
	"clash-admin/internal/logging"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitDisabled(t *testing.T) {
	t.Cleanup(func() { _ = logging.Close() })

	dir := t.TempDir()
	logPath := filepath.Join(dir, "app.log")

	if err := logging.Init(config.Logging{Enabled: false, FilePath: logPath}); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if logging.Enabled() {
		t.Fatal("Enabled() = true, want false")
	}

	logging.Info("should not write")
	if _, err := os.Stat(logPath); !os.IsNotExist(err) {
		t.Fatalf("log file should not exist when disabled, stat err = %v", err)
	}
}

func TestInitEnabledWritesJSON(t *testing.T) {
	t.Cleanup(func() { _ = logging.Close() })

	dir := t.TempDir()
	logPath := filepath.Join(dir, "nested", "clash-admin.log")

	if err := logging.Init(config.Logging{
		Enabled:  true,
		Level:    "info",
		FilePath: logPath,
	}); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if !logging.Enabled() {
		t.Fatal("Enabled() = false, want true")
	}

	logging.Info("starting server", "component", "runtime", "service", "http")
	logging.Warn("tool call failed", "component", "mcp-tool", "tool", "doc_generation")

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d log lines, want 2", len(lines))
	}

	var infoEntry map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &infoEntry); err != nil {
		t.Fatalf("unmarshal info line: %v", err)
	}
	if infoEntry["level"] != "INFO" {
		t.Fatalf("info level = %v, want INFO", infoEntry["level"])
	}
	if infoEntry["msg"] != "starting server" {
		t.Fatalf("info msg = %v, want starting server", infoEntry["msg"])
	}
	if infoEntry["component"] != "runtime" {
		t.Fatalf("info component = %v, want runtime", infoEntry["component"])
	}
	if infoEntry["service"] != "http" {
		t.Fatalf("info service = %v, want http", infoEntry["service"])
	}

	var warnEntry map[string]any
	if err := json.Unmarshal([]byte(lines[1]), &warnEntry); err != nil {
		t.Fatalf("unmarshal warn line: %v", err)
	}
	if warnEntry["level"] != "WARN" {
		t.Fatalf("warn level = %v, want WARN", warnEntry["level"])
	}
	if warnEntry["tool"] != "doc_generation" {
		t.Fatalf("warn tool = %v, want doc_generation", warnEntry["tool"])
	}
}

func TestInitInvalidLevelFallsBackToInfo(t *testing.T) {
	t.Cleanup(func() { _ = logging.Close() })

	dir := t.TempDir()
	logPath := filepath.Join(dir, "clash-admin.log")

	if err := logging.Init(config.Logging{
		Enabled:  true,
		Level:    "unknown",
		FilePath: logPath,
	}); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	logging.Info("info message")
	logging.Warn("warn message")

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "info message") || !strings.Contains(content, "warn message") {
		t.Fatalf("expected both info and warn messages in log, got: %s", content)
	}
}

func TestInitDefaultFilePathWhenEmpty(t *testing.T) {
	t.Cleanup(func() { _ = logging.Close() })

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	if err := logging.Init(config.Logging{Enabled: true}); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	logging.Info("default path message")

	logPath := filepath.Join(dir, "runtime", "logs", "clash-admin.log")
	if _, err := os.Stat(logPath); err != nil {
		t.Fatalf("default log file stat error = %v", err)
	}
}

func TestInitLevelFiltersDebug(t *testing.T) {
	t.Cleanup(func() { _ = logging.Close() })

	dir := t.TempDir()
	logPath := filepath.Join(dir, "clash-admin.log")

	if err := logging.Init(config.Logging{
		Enabled:  true,
		Level:    "warn",
		FilePath: logPath,
	}); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	logging.Info("filtered info")
	logging.Warn("kept warn")

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	content := string(data)
	if strings.Contains(content, "filtered info") {
		t.Fatalf("info message should be filtered at warn level, got: %s", content)
	}
	if !strings.Contains(content, "kept warn") {
		t.Fatalf("warn message should be kept, got: %s", content)
	}
}

func TestDebugAndErrorLevels(t *testing.T) {
	t.Cleanup(func() { _ = logging.Close() })

	dir := t.TempDir()
	logPath := filepath.Join(dir, "clash-admin.log")

	if err := logging.Init(config.Logging{
		Enabled:  true,
		Level:    "debug",
		FilePath: logPath,
	}); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	logging.Debug("debug message", "component", "test")
	logging.Error("error message", "component", "test")

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "debug message") || !strings.Contains(content, "error message") {
		t.Fatalf("expected debug and error messages in log, got: %s", content)
	}
}

func TestReInitClosesPreviousFile(t *testing.T) {
	t.Cleanup(func() { _ = logging.Close() })

	dir := t.TempDir()
	firstPath := filepath.Join(dir, "first.log")
	secondPath := filepath.Join(dir, "second.log")

	if err := logging.Init(config.Logging{Enabled: true, FilePath: firstPath}); err != nil {
		t.Fatalf("first Init() error = %v", err)
	}
	logging.Info("first")

	if err := logging.Init(config.Logging{Enabled: true, FilePath: secondPath}); err != nil {
		t.Fatalf("second Init() error = %v", err)
	}
	logging.Info("second")

	firstData, err := os.ReadFile(firstPath)
	if err != nil {
		t.Fatalf("ReadFile(first) error = %v", err)
	}
	if !strings.Contains(string(firstData), "first") || strings.Contains(string(firstData), "second") {
		t.Fatalf("first log file content unexpected: %s", firstData)
	}

	secondData, err := os.ReadFile(secondPath)
	if err != nil {
		t.Fatalf("ReadFile(second) error = %v", err)
	}
	if !strings.Contains(string(secondData), "second") {
		t.Fatalf("second log file content unexpected: %s", secondData)
	}
}
