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
      "total": 16,
      "used": 8,
      "free": 8,
      "usage": 50.0
    },
    "jvm": {
      "total": 512,
      "used": 256,
      "free": 256,
      "usage": 50.0,
      "name": "OpenJDK",
      "version": "17",
      "startTime": "2024-01-01 00:00:00",
      "runTime": "1天2小时",
      "home": "/usr/lib/jvm",
      "inputArgs": "-Xms256m"
    },
    "sys": {
      "computerName": "host",
      "osName": "Linux",
      "computerIp": "10.0.0.1",
      "osArch": "amd64",
      "userDir": "/app"
    },
    "sysFiles": [
      {
        "dirName": "/",
        "sysTypeName": "ext4",
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
| `mem` | `total/used/free/usage` | 内存（GB 或后端约定单位） |
| `jvm` | 见上 | 运行时信息（Go 实现可映射为进程/运行时等价信息，字段名保持兼容） |
| `sys` | 主机与 OS | |
| `sysFiles` | 磁盘分区列表 | |

> 注：原若依后端为 Java/JVM 指标。Go 实现时应尽量保持**字段名与结构兼容**，数值语义用进程/runtime/系统指标填充。
