# 公共接口（上传 / 下载）

前端：`src/plugins/download.ts`、`ImageUpload` / `FileUpload` / `Editor` 组件

除非特别说明，均需 Bearer Token。

---

## POST `/common/upload`

通用文件上传（multipart）。

### 请求

- `Content-Type: multipart/form-data`
- 字段名通常为 `file`（与 Element Plus Upload 默认一致）

### 响应

```json
{
  "code": 200,
  "msg": "操作成功",
  "fileName": "/profile/upload/2024/01/01/xxx.png",
  "newFileName": "xxx.png",
  "url": "http://host/profile/upload/2024/01/01/xxx.png",
  "originalFilename": "photo.png"
}
```

| 字段 | 说明 |
|------|------|
| `fileName` | 含路径的存储名（业务常存此值） |
| `newFileName` | 不含路径的新文件名 |
| `url` | 可访问完整 URL |
| `originalFilename` | 原始文件名 |

字段与 `AjaxResult` 平铺在同一层。

---

## GET `/common/download`

按文件名下载（常用于临时导出文件）。

### Query

| 参数 | 类型 | 说明 |
|------|------|------|
| `fileName` | string | 文件名/路径（需 URL 编码） |
| `delete` | boolean | 下载后是否删除，`true`/`false` |

### 响应

- 成功：文件二进制流；响应头建议含 `download-filename`（URL 编码）
- 失败：JSON `{ code, msg }`

### 示例

```
GET /common/download?fileName=%2Fprofile%2Fxxx.xlsx&delete=true
```

---

## GET `/common/download/resource`

按资源路径下载。

### Query

| 参数 | 类型 | 说明 |
|------|------|------|
| `resource` | string | 资源标识/路径（URL 编码） |

### 响应

同文件二进制流。
