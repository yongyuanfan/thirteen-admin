# 操作日志

前端：`src/api/monitor/operlog.ts`  
类型：`src/types/api/monitor/operlog.ts`

---

## 数据模型 SysOperLog

| 字段 | 类型 | 说明 |
|------|------|------|
| `operId` | number | 日志主键 |
| `title` | string | 模块标题 |
| `businessType` | number | 0其它 1新增 2修改 3删除 4授权 5导出 6导入 7强退 8生成代码 9清空 |
| `method` | string | 方法名 |
| `requestMethod` | string | HTTP 方法 |
| `operatorType` | string/number | 操作类别 |
| `operName` | string | 操作人员 |
| `deptName` | string | 部门名称 |
| `operUrl` | string | 请求 URL |
| `operIp` | string | 主机 IP |
| `operLocation` | string | 操作地点 |
| `operParam` | string | 请求参数 |
| `jsonResult` | string | 返回参数 |
| `status` | `"0"\|"1"` | 0正常 1异常 |
| `errorMsg` | string | 错误消息 |
| `operTime` | string/Date | 操作时间 |
| `costTime` | number | 耗时（毫秒） |

---

## GET `/monitor/operlog/list`

### Query

`operIp`、`title`、`operName`、`businessType`、`status`、`params[beginTime/endTime]`、分页。

### 响应

`TableDataInfo<SysOperLog>`

---

## DELETE `/monitor/operlog/{operId}`

支持多 ID。

---

## DELETE `/monitor/operlog/clean`

清空操作日志。

---

## POST `/monitor/operlog/export`

导出 Excel。
