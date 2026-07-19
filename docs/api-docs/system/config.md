# 参数配置

前端：`src/api/system/config.ts`  
类型：`src/types/api/system/config.ts`

---

## 数据模型 SysConfig

| 字段 | 类型 | 说明 |
|------|------|------|
| `configId` | number | 参数主键 |
| `configName` | string | 参数名称 |
| `configKey` | string | 参数键名 |
| `configValue` | string | 参数键值 |
| `configType` | `"Y"\|"N"` | 是否系统内置 |
| + BaseEntity | | |

---

## GET `/system/config/list`

Query：`configName`、`configKey`、`configType`、`params[beginTime/endTime]`、分页。  
响应：`TableDataInfo<SysConfig>`

---

## GET `/system/config/{configId}`

`{ "code": 200, "data": SysConfig }`

---

## GET `/system/config/configKey/{configKey}`

按键名获取配置值。

### 重要约定

**配置值放在响应字段 `msg` 中**（不是 `data`）。前端直接使用 `res.msg`。

```json
{
  "code": 200,
  "msg": "123456"
}
```

前端已知使用：

| configKey | 用途 |
|-----------|------|
| `sys.user.initPassword` | 用户管理页展示初始密码 |

其它内置键见 `docs/thirteen-admin.sql` 种子数据（皮肤、验证码开关、注册开关、密码策略等）。

---

## POST `/system/config` / PUT `/system/config`

新增 / 修改。Body：`SysConfig`。

---

## DELETE `/system/config/{configId}`

支持多 ID。

---

## DELETE `/system/config/refreshCache`

刷新参数缓存。

---

## POST `/system/config/export`

导出 Excel。
