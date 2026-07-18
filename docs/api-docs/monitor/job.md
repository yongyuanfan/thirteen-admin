# 定时任务

前端：`src/api/monitor/job.ts`、`jobLog.ts`  
类型：`src/types/api/monitor/job.ts`、`jobLog.ts`

---

## 数据模型 SysJob

| 字段 | 类型 | 说明 |
|------|------|------|
| `jobId` | number | 任务 ID |
| `jobName` | string | 任务名称 |
| `jobGroup` | string | 任务组名 |
| `invokeTarget` | string | 调用目标字符串 |
| `cronExpression` | string | cron 表达式 |
| `nextValidTime` | string/Date | 下次执行时间（只读） |
| `misfirePolicy` | `"1"\|"2"\|"3"` | 1立即 2执行一次 3放弃 |
| `concurrent` | `"0"\|"1"` | 0允许 1禁止 |
| `status` | `"0"\|"1"` | 0正常 1暂停 |
| + BaseEntity | | |

## 数据模型 SysJobLog

| 字段 | 类型 | 说明 |
|------|------|------|
| `jobLogId` | number | 日志 ID |
| `jobName` | string | 任务名称 |
| `jobGroup` | string | 组名 |
| `invokeTarget` | string | 调用目标 |
| `jobMessage` | string | 日志信息 |
| `exceptionInfo` | string | 异常信息 |
| `status` | `"0"\|"1"` | 0正常 1失败 |
| + BaseEntity | | |

---

## 任务接口

### GET `/monitor/job/list`

Query：`jobName`、`jobGroup`、`status`、分页。  
响应：`TableDataInfo<SysJob>`

### GET `/monitor/job/{jobId}`

`{ "code": 200, "data": SysJob }`

### POST `/monitor/job` / PUT `/monitor/job`

新增 / 修改。Body：`SysJob`。

### DELETE `/monitor/job/{jobId}`

支持多 ID。

### PUT `/monitor/job/changeStatus`

```json
{ "jobId": 1, "status": "1" }
```

### PUT `/monitor/job/run`

立即执行一次。

```json
{ "jobId": 1, "jobGroup": "DEFAULT" }
```

### POST `/monitor/job/export`

导出 Excel。

---

## 任务日志接口

### GET `/monitor/jobLog/list`

Query：`jobName`、`jobGroup`、`status`、`params[beginTime/endTime]`、分页。  
响应：`TableDataInfo<SysJobLog>`

### DELETE `/monitor/jobLog/{jobLogId}`

支持多 ID。

### DELETE `/monitor/jobLog/clean`

清空全部任务日志。

### POST `/monitor/jobLog/export`

导出 Excel。
