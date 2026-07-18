# 角色管理

前端：`src/api/system/role.ts`  
类型：`src/types/api/system/role.ts`

---

## 数据模型 SysRole

| 字段 | 类型 | 说明 |
|------|------|------|
| `roleId` | number | 角色 ID |
| `roleName` | string | 角色名称 |
| `roleKey` | string | 权限字符（如 `admin`） |
| `roleSort` | number | 显示顺序 |
| `dataScope` | `"1"`–`"5"` | 1全部 2自定义 3本部门 4本部门及以下 5仅本人 |
| `menuCheckStrictly` | boolean | 菜单树父子联动 |
| `deptCheckStrictly` | boolean | 部门树父子联动 |
| `menuIds` | number[] | 菜单 ID（授权提交） |
| `deptIds` | number[] | 部门 ID（数据权限提交） |
| `status` | `"0"\|"1"` | 0正常 1停用 |
| + BaseEntity | | |

---

## GET `/system/role/list`

### Query

`roleName`、`roleKey`、`status`、`params[beginTime/endTime]`、分页排序。

### 响应

`TableDataInfo<SysRole>`

---

## GET `/system/role/{roleId}`

### 响应

`{ "code": 200, "data": SysRole }`

---

## POST `/system/role` / PUT `/system/role`

新增 / 修改。Body：`SysRole`（可含 `menuIds`）。

---

## PUT `/system/role/dataScope`

保存数据权限。Body：`SysRole`（含 `roleId`、`dataScope`、`deptIds`、`deptCheckStrictly` 等）。

---

## PUT `/system/role/changeStatus`

```json
{ "roleId": 2, "status": "1" }
```

---

## DELETE `/system/role/{roleId}`

支持多 ID：`1,2,3`。

---

## GET `/system/role/authUser/allocatedList`

已分配该角色的用户（分页）。

### Query

`roleId` + 用户筛选（`userName`、`phonenumber` 等）+ 分页。

### 响应

`TableDataInfo<SysUser>`

---

## GET `/system/role/authUser/unallocatedList`

未分配该角色的用户。Query 同上。

---

## PUT `/system/role/authUser/cancel`

取消单个用户授权。

### 请求体

```json
{ "userId": 2, "roleId": 2 }
```

---

## PUT `/system/role/authUser/cancelAll`

批量取消。**Query：**

| 参数 | 说明 |
|------|------|
| `roleId` | 角色 ID |
| `userIds` | 用户 ID，逗号分隔 |

---

## PUT `/system/role/authUser/selectAll`

批量授权。**Query：** `roleId`、`userIds`（逗号分隔）。

---

## GET `/system/role/deptTree/{roleId}`

数据权限部门树。

### 响应

```json
{
  "code": 200,
  "checkedKeys": [100, 101],
  "depts": [
    { "id": 100, "label": "若依科技", "disabled": false, "children": [] }
  ]
}
```

---

## 关联菜单树

角色菜单授权使用菜单模块接口（见 [menu.md](./menu.md)）：

- `GET /system/menu/treeselect`
- `GET /system/menu/roleMenuTreeselect/{roleId}` → `{ menus, checkedKeys }`

---

## POST `/system/role/export`

导出 Excel。筛选同 list。
