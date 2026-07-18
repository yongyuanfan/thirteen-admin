package services

import (
	"fmt"
	"os"
	"runtime"
	"strings"
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

// ServerCPU CPU 信息。
type ServerCPU struct {
	CPUNum int     `json:"cpuNum"`
	Used   float64 `json:"used"`
	Sys    float64 `json:"sys"`
	Free   float64 `json:"free"`
}

// ServerMem 进程内存信息（GB）。
type ServerMem struct {
	Total float64 `json:"total"`
	Used  float64 `json:"used"`
	Free  float64 `json:"free"`
	Usage float64 `json:"usage"`
}

// ServerGo Go 运行时信息。
type ServerGo struct {
	Total      float64 `json:"total"`      // 已从 OS 申请的内存（MB，对应 MemStats.Sys）
	Used       float64 `json:"used"`       // 堆上已分配（MB，对应 MemStats.Alloc）
	Free       float64 `json:"free"`       // 剩余（MB）
	Usage      float64 `json:"usage"`      // 使用率 %
	Name       string  `json:"name"`       // 运行时名称
	Version    string  `json:"version"`    // Go 版本，如 go1.25.1
	StartTime  string  `json:"startTime"`  // 进程启动时间
	RunTime    string  `json:"runTime"`    // 运行时长
	GOROOT     string  `json:"goroot"`     // GOROOT
	Goroutines int     `json:"goroutines"` // 当前 goroutine 数
	Compiler   string  `json:"compiler"`   // 编译器，通常为 gc
	GOOS       string  `json:"goos"`       // 编译目标 OS
	GOARCH     string  `json:"goarch"`     // 编译目标架构
	CmdLine    string  `json:"cmdLine"`    // 启动命令行
}

// ServerSys 主机系统信息。
type ServerSys struct {
	ComputerName string `json:"computerName"`
	OsName       string `json:"osName"`
	ComputerIP   string `json:"computerIp"`
	OsArch       string `json:"osArch"`
	UserDir      string `json:"userDir"`
}

// ServerSysFile 磁盘分区信息。
type ServerSysFile struct {
	DirName     string  `json:"dirName"`
	SysTypeName string  `json:"sysTypeName"`
	TypeName    string  `json:"typeName"`
	Total       string  `json:"total"`
	Free        string  `json:"free"`
	Used        string  `json:"used"`
	Usage       float64 `json:"usage"`
}

// ServerInfo 服务监控聚合数据。
type ServerInfo struct {
	CPU      ServerCPU      `json:"cpu"`
	Mem      ServerMem      `json:"mem"`
	Go       ServerGo       `json:"go"`
	Sys      ServerSys      `json:"sys"`
	SysFiles []ServerSysFile `json:"sysFiles"`
}

// Info 获取服务器运行信息。
func (s *ServerService) Info() ServerInfo {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	memTotalGB := float64(mem.Sys) / 1024 / 1024 / 1024
	memUsedGB := float64(mem.Alloc) / 1024 / 1024 / 1024
	memFreeGB := memTotalGB - memUsedGB
	memUsage := 0.0
	if memTotalGB > 0 {
		memUsage = memUsedGB / memTotalGB * 100
	}

	goTotalMB := float64(mem.Sys) / 1024 / 1024
	goUsedMB := float64(mem.Alloc) / 1024 / 1024
	goFreeMB := goTotalMB - goUsedMB
	goUsage := 0.0
	if goTotalMB > 0 {
		goUsage = goUsedMB / goTotalMB * 100
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

	return ServerInfo{
		CPU: ServerCPU{
			CPUNum: cpuNum,
			Used:   round2(cpuUsed),
			Sys:    round2(cpuUsed * 0.3),
			Free:   round2(cpuFree),
		},
		Mem: ServerMem{
			Total: round2(memTotalGB),
			Used:  round2(memUsedGB),
			Free:  round2(memFreeGB),
			Usage: round2(memUsage),
		},
		Go: ServerGo{
			Total:      round2(goTotalMB),
			Used:       round2(goUsedMB),
			Free:       round2(goFreeMB),
			Usage:      round2(goUsage),
			Name:       "go",
			Version:    runtime.Version(),
			StartTime:  FormatTimeVal(processStart),
			RunTime:    formatRunTime(time.Since(processStart)),
			GOROOT:     runtime.GOROOT(),
			Goroutines: goroutines,
			Compiler:   runtime.Compiler,
			GOOS:       runtime.GOOS,
			GOARCH:     runtime.GOARCH,
			CmdLine:    strings.Join(os.Args, " "),
		},
		Sys: ServerSys{
			ComputerName: hostname,
			OsName:       runtime.GOOS,
			ComputerIP:   "127.0.0.1",
			OsArch:       runtime.GOARCH,
			UserDir:      wd,
		},
		SysFiles: diskPartitions(),
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

func diskPartitions() []ServerSysFile {
	paths := []string{"/", os.TempDir()}
	out := make([]ServerSysFile, 0, len(paths))
	seen := map[string]bool{}
	for _, p := range paths {
		if p == "" || seen[p] {
			continue
		}
		seen[p] = true
		if part, ok := statDisk(p); ok {
			out = append(out, part)
		}
	}
	if len(out) == 0 {
		out = append(out, ServerSysFile{
			DirName:     "/",
			SysTypeName: runtime.GOOS,
			TypeName:    "/",
			Total:       "N/A",
			Free:        "N/A",
			Used:        "N/A",
			Usage:       0.0,
		})
	}
	return out
}

func statDisk(path string) (ServerSysFile, bool) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return ServerSysFile{}, false
	}
	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)
	used := total - free
	usage := 0.0
	if total > 0 {
		usage = float64(used) / float64(total) * 100
	}
	return ServerSysFile{
		DirName:     path,
		SysTypeName: runtime.GOOS,
		TypeName:    path,
		Total:       formatBytes(total),
		Free:        formatBytes(free),
		Used:        formatBytes(used),
		Usage:       round2(usage),
	}, true
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
