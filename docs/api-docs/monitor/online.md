# 在线用户

前端：`src/api/monitor/online.ts`  
类型：`src/types/api/monitor/online.ts`

---

## 数据模型 SysUserOnline

| 字段 | 类型 | 说明 |
|------|------|------|
| `tokenId` | string | 会话编号（强退用） |
| `deptName` | string | 部门名称 |
| `userName` | string | 用户名称 |
| `ipaddr` | string | 登录 IP |
| `loginLocation` | string | 登录地址 |
| `browser` | string | 浏览器 |
| `os` | string | 操作系统 |
| `loginTime` | number | 登录时间（时间戳） |

---

## GET `/monitor/online/list`

### Query

`ipaddr`、`userName`（可无分页字段，若依实现常返回全量再前端分页，或带 `total`/`rows`）

### 响应

`TableDataInfo<SysUserOnline>`：`{ code, msg, total, rows }`

---

## DELETE `/monitor/online/{tokenId}`

强退指定会话。`tokenId` 为会话标识字符串。
