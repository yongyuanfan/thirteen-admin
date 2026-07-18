# 登录日志

前端：`src/api/monitor/logininfor.ts`  
类型：`src/types/api/monitor/logininfor.ts`

---

## 数据模型 SysLogininfor

| 字段 | 类型 | 说明 |
|------|------|------|
| `infoId` | number | 访问 ID |
| `userName` | string | 用户账号 |
| `ipaddr` | string | 登录 IP |
| `loginLocation` | string | 登录地点 |
| `browser` | string | 浏览器 |
| `os` | string | 操作系统 |
| `status` | `"0"\|"1"` | 0成功 1失败 |
| `msg` | string | 提示消息 |
| `loginTime` | string/Date | 访问时间 |

---

## GET `/monitor/logininfor/list`

### Query

`ipaddr`、`userName`、`status`、`params[beginTime/endTime]`、分页。

### 响应

`TableDataInfo<SysLogininfor>`

---

## DELETE `/monitor/logininfor/{infoId}`

支持多 ID。

---

## GET `/monitor/logininfor/unlock/{userName}`

解锁因登录失败被锁定的账号。

---

## DELETE `/monitor/logininfor/clean`

清空登录日志。

---

## POST `/monitor/logininfor/export`

导出 Excel。
