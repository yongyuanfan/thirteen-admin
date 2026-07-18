# 部门管理

前端：`src/api/system/dept.ts`  
类型：`src/types/api/system/dept.ts`

---

## 数据模型 SysDept

| 字段 | 类型 | 说明 |
|------|------|------|
| `deptId` | number | 部门 ID |
| `parentId` | number | 父部门 ID |
| `ancestors` | string | 祖级列表，如 `0,100,101` |
| `deptName` | string | 部门名称 |
| `orderNum` | number | 显示顺序 |
| `leader` | string | 负责人 |
| `phone` | string | 电话 |
| `email` | string | 邮箱 |
| `status` | `"0"\|"1"` | 0正常 1停用 |
| `children` | SysDept[] | 子部门 |
| + BaseEntity | | |

---

## GET `/system/dept/list`

### Query

`deptName`、`status`

### 响应

`{ "code": 200, "data": SysDept[] }`（扁平或树，与若依后端一致时多为扁平）

---

## GET `/system/dept/list/exclude/{deptId}`

查询部门列表并排除指定节点及其子孙（修改上级时用）。

### 响应

`{ "code": 200, "data": SysDept[] }`

---

## GET `/system/dept/{deptId}`

### 响应

`{ "code": 200, "data": SysDept }`

---

## POST `/system/dept` / PUT `/system/dept`

新增 / 修改。Body：`SysDept`。

---

## PUT `/system/dept/updateSort`

```json
{
  "deptIds": "100,101,102",
  "orderNums": "1,2,3"
}
```

---

## DELETE `/system/dept/{deptId}`

删除部门（通常要求无子部门、无用户）。
