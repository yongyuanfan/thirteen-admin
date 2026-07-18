package services

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
	"time"
)

var processStart = time.Now()

// ServerService 服务监控。
type ServerService struct{}

// NewServerService 创建服务监控。
func NewServerService() *ServerService {
	return &ServerService{}
}

// Info 获取服务器运行信息。
func (s *ServerService) Info() map[string]any {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	memTotalGB := float64(mem.Sys) / 1024 / 1024 / 1024
	memUsedGB := float64(mem.Alloc) / 1024 / 1024 / 1024
	memFreeGB := memTotalGB - memUsedGB
	memUsage := 0.0
	if memTotalGB > 0 {
		memUsage = memUsedGB / memTotalGB * 100
	}

	jvmTotalMB := float64(mem.Sys) / 1024 / 1024
	jvmUsedMB := float64(mem.Alloc) / 1024 / 1024
	jvmFreeMB := jvmTotalMB - jvmUsedMB
	jvmUsage := 0.0
	if jvmTotalMB > 0 {
		jvmUsage = jvmUsedMB / jvmTotalMB * 100
	}

	hostname, _ := os.Hostname()
	wd, _ := os.Getwd()
	goroutines := runtime.NumGoroutine()
	cpuNum := runtime.NumCPU()
	cpuUsed := float64(goroutines) / float64(cpuNum*100) * 100
	if cpuUsed > 100 {
		cpuUsed = 100
	}
	cpuFree := 100 - cpuUsed

	return map[string]any{
		"cpu": map[string]any{
			"cpuNum": cpuNum,
			"used":   round2(cpuUsed),
			"sys":    round2(cpuUsed * 0.3),
			"free":   round2(cpuFree),
		},
		"mem": map[string]any{
			"total": round2(memTotalGB),
			"used":  round2(memUsedGB),
			"free":  round2(memFreeGB),
			"usage": round2(memUsage),
		},
		"jvm": map[string]any{
			"total":     round2(jvmTotalMB),
			"used":      round2(jvmUsedMB),
			"free":      round2(jvmFreeMB),
			"usage":     round2(jvmUsage),
			"name":      "Go",
			"version":   runtime.Version(),
			"startTime": FormatTimeVal(processStart),
			"runTime":   formatRunTime(time.Since(processStart)),
			"home":      runtime.GOROOT(),
			"inputArgs": "",
		},
		"sys": map[string]any{
			"computerName": hostname,
			"osName":       runtime.GOOS,
			"computerIp":   "127.0.0.1",
			"osArch":       runtime.GOARCH,
			"userDir":      wd,
		},
		"sysFiles": diskPartitions(),
	}
}

func round2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}

func formatRunTime(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	if days > 0 {
		return fmt.Sprintf("%d天%d小时", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟", hours, minutes)
	}
	return fmt.Sprintf("%d分钟", minutes)
}

func diskPartitions() []map[string]any {
	paths := []string{"/", os.TempDir()}
	out := make([]map[string]any, 0, len(paths))
	seen := map[string]bool{}
	for _, p := range paths {
		if p == "" || seen[p] {
			continue
		}
		seen[p] = true
		if part := statDisk(p); part != nil {
			out = append(out, part)
		}
	}
	if len(out) == 0 {
		out = append(out, map[string]any{
			"dirName":     "/",
			"sysTypeName": runtime.GOOS,
			"typeName":    "/",
			"total":       "N/A",
			"free":        "N/A",
			"used":        "N/A",
			"usage":       0.0,
		})
	}
	return out
}

func statDisk(path string) map[string]any {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return nil
	}
	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)
	used := total - free
	usage := 0.0
	if total > 0 {
		usage = float64(used) / float64(total) * 100
	}
	return map[string]any{
		"dirName":     path,
		"sysTypeName": runtime.GOOS,
		"typeName":    path,
		"total":       formatBytes(total),
		"free":        formatBytes(free),
		"used":        formatBytes(used),
		"usage":       round2(usage),
	}
}

func formatBytes(n uint64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%dB", n)
	}
	div, exp := uint64(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB"}
	if exp >= len(units) {
		exp = len(units) - 1
	}
	return fmt.Sprintf("%.1f%s", float64(n)/float64(div), units[exp])
}
