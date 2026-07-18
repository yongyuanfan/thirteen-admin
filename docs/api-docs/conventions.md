# 通用约定

## 1. 传输与 Content-Type

| 场景 | Content-Type |
|------|----------------|
| 默认 JSON | `application/json;charset=utf-8` |
| 表单导出 / 部分上传 | `application/x-www-form-urlencoded` |
| 文件上传 | `multipart/form-data` |
| 导出 / 下载 | 请求可带 form；响应为二进制 blob |

## 2. 鉴权

- Header：`Authorization: Bearer <token>`
- Token 来源：登录接口返回的 `token`；前端 Cookie 名 `Admin-Token`
- 以下接口**不需要** Token（请求头可设 `isToken: false`）：
  - `POST /login`
  - `POST /register`
  - `GET /captchaImage`
- 其余业务接口均需登录；Token 无效或过期时业务码为 `401`

## 3. 响应信封

### 3.1 AjaxResult（通用）

```json
{
  "code": 200,
  "msg": "操作成功",
  "data": {}
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | `number` | 业务状态码，**成功为 `200`**（不是字符串） |
| `msg` | `string` | 提示信息 |
| `data` | `any` | 可选；业务数据 |

部分接口会把关键字段**平铺在根级**（与 `data` 同级），例如登录的 `token`、`getInfo` 的 `user`/`roles`/`permissions`、上传的 `url`/`fileName` 等。前端以整包 `res.data` 使用，不单取 `data`。

### 3.2 TableDataInfo（分页列表）

```json
{
  "code": 200,
  "msg": "查询成功",
  "total": 100,
  "rows": []
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `total` | `number` | 总记录数 |
| `rows` | `array` | 当前页数据 |

### 3.3 业务码

| code | 含义 | 前端行为 |
|------|------|----------|
| `200` | 成功 | 正常处理 |
| `401` | 未认证 / 会话过期 | 弹窗引导重新登录 |
| `500` | 服务错误 | Message 错误提示 |
| `601` | 业务警告 | Message 警告 |
| 其它非 200 | 失败 | Notification 错误 |

HTTP 状态码通常为 `200`；以响应体 `code` 为准。

## 4. 分页与排序（Query）

列表类 `GET` 通用查询参数（`PageDomain`）：

| 参数 | 类型 | 说明 |
|------|------|------|
| `pageNum` | number | 页码，从 1 开始 |
| `pageSize` | number | 每页条数 |
| `orderByColumn` | string | 排序列（属性名） |
| `isAsc` | string | `asc` / `desc` |
| `reasonable` | boolean | 分页合理化（可选） |

时间范围常用嵌套参数（前端序列化为 `params[beginTime]` / `params[endTime]`）：

```
params[beginTime]=2024-01-01
params[endTime]=2024-12-31
```

## 5. 路径中的多 ID

删除等接口支持逗号拼接多个主键：

```
DELETE /system/user/1,2,3
```

## 6. 导出约定

前端调用：`POST /{module}/export`，`Content-Type: application/x-www-form-urlencoded`，`responseType: blob`。

- 请求体：与列表查询相同的筛选条件（form 字段）
- 成功：Excel 二进制流；建议响应头带 `download-filename`（URL 编码文件名）
- 失败：仍可能返回 JSON `{ code, msg }`（前端会校验 blob 是否为 JSON 错误）

已知导出端点：

| Method | Path |
|--------|------|
| POST | `/system/user/export` |
| POST | `/system/role/export` |
| POST | `/system/post/export` |
| POST | `/system/dict/type/export` |
| POST | `/system/dict/data/export` |
| POST | `/system/config/export` |
| POST | `/monitor/job/export` |
| POST | `/monitor/jobLog/export` |
| POST | `/monitor/operlog/export` |
| POST | `/monitor/logininfor/export` |

## 7. 通用实体字段（BaseEntity）

多数业务对象可包含：

| 字段 | 类型 | 说明 |
|------|------|------|
| `createdBy` | string | 创建者 |
| `createdAt` | string | 创建时间（前端展示用字符串/日期） |
| `updatedBy` | string | 更新者 |
| `updatedAt` | string | 更新时间 |
| `remark` | string | 备注 |
| `params` | object | 请求扩展参数（如时间范围） |
| `searchValue` | string | 搜索关键字（部分接口） |

> 说明：HTTP JSON 与 Go 模型一致使用 `createdAt` / `createdBy`；数据库列为 `created_at` / `created_by`。

## 8. TreeSelect

```json
{
  "id": 100,
  "label": "研发部门",
  "disabled": false,
  "children": []
}
```

## 9. 状态与标志枚举（通用）

| 场景 | 值 |
|------|-----|
| 正常 / 停用 | `"0"` 正常，`"1"` 停用 |
| 显示 / 隐藏 | `"0"` 显示，`"1"` 隐藏 |
| 是否（Y/N） | `"Y"` / `"N"` |
| 删除标志（库） | `"0"` 存在，`"2"` 删除（列表接口一般不返回已删数据） |

## 10. curl 示例模板

```bash
# 登录
curl -s -X POST 'http://localhost:8080/login' \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123","code":"1","uuid":"..."}'

# 带 Token 的业务请求
curl -s 'http://localhost:8080/system/user/list?pageNum=1&pageSize=10' \
  -H 'Authorization: Bearer <token>'
```
