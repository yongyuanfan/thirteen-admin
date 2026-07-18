# 岗位管理

前端：`src/api/system/post.ts`  
类型：`src/types/api/system/post.ts`

---

## 数据模型 SysPost

| 字段 | 类型 | 说明 |
|------|------|------|
| `postId` | number | 岗位 ID |
| `postCode` | string | 岗位编码 |
| `postName` | string | 岗位名称 |
| `postSort` | number | 显示顺序 |
| `status` | `"0"\|"1"` | 0正常 1停用 |
| + BaseEntity | | |

---

## GET `/system/post/list`

### Query

`postCode`、`postName`、`status`、分页排序。

### 响应

`TableDataInfo<SysPost>`

---

## GET `/system/post/{postId}`

`{ "code": 200, "data": SysPost }`

---

## POST `/system/post` / PUT `/system/post`

新增 / 修改。Body：`SysPost`。

---

## DELETE `/system/post/{postId}`

支持多 ID。

---

## POST `/system/post/export`

导出 Excel。
