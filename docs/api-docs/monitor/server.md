# 服务监控

前端：`src/api/monitor/server.ts`  
页面：`views/monitor/server/index.vue`

---

## GET `/monitor/server`

获取服务器运行信息。**需要 Token。**

### 响应

```json
{
  "code": 200,
  "msg": "操作成功",
  "data": {
    "cpu": {
      "cpuNum": 8,
      "used": 12.5,
      "sys": 5.0,
      "free": 82.5
    },
    "mem": {
      "total": 0.5,
      "used": 0.2,
      "free": 0.3,
      "usage": 40.0
    },
    "go": {
      "total": 512,
      "used": 256,
      "free": 256,
      "usage": 50.0,
      "name": "go",
      "version": "go1.25.1",
      "startTime": "2026-07-18 15:00:00",
      "runTime": "1小时2分钟",
      "goroot": "/usr/local/go",
      "goroutines": 42,
      "compiler": "gc",
      "goos": "darwin",
      "goarch": "arm64",
      "cmdLine": "thirteen-admin server -c config/settings.local.yml"
    },
    "sys": {
      "computerName": "host",
      "osName": "darwin",
      "computerIp": "10.0.0.1",
      "osArch": "arm64",
      "userDir": "/app"
    },
    "sysFiles": [
      {
        "dirName": "/",
        "sysTypeName": "darwin",
        "typeName": "/",
        "total": "100GB",
        "free": "40GB",
        "used": "60GB",
        "usage": 60.0
      }
    ]
  }
}
```

### 字段说明

| 块 | 字段 | 说明 |
|----|------|------|
| `cpu` | `cpuNum/used/sys/free` | CPU 核数与占用百分比 |
| `mem` | `total/used/free/usage` | 进程内存（GB） |
| `go` | 见下 | Go 运行时信息 |
| `sys` | 主机与 OS | |
| `sysFiles` | 磁盘分区列表 | |

#### `go` 字段

| 字段 | 说明 |
|------|------|
| `total/used/free/usage` | 运行时内存（MB）：`Sys` / `Alloc` / 剩余 / 使用率 |
| `name` | 运行时名称，固定为 `go` |
| `version` | `runtime.Version()`，如 `go1.25.1` |
| `startTime` / `runTime` | 进程启动时间与运行时长 |
| `goroot` | `runtime.GOROOT()` |
| `goroutines` | 当前 goroutine 数量 |
| `compiler` | 编译器（通常 `gc`） |
| `goos` / `goarch` | 编译目标 OS / 架构 |
| `cmdLine` | 启动命令行（`os.Args`） |
