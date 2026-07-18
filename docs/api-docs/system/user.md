# 用户管理

前端：`src/api/system/user.ts`  
类型：`src/types/api/system/user.ts`  
页面：`views/system/user/**`、个人中心

所有接口需 Token（除特别说明）。

---

## 数据模型 SysUser

| 字段 | 类型 | 说明 |
|------|------|------|
| `userId` | number | 用户 ID |
| `deptId` | number | 部门 ID |
| `userName` | string | 账号 |
| `nickName` | string | 昵称 |
| `email` | string | 邮箱 |
| `phonenumber` | string | 手机号 |
| `sex` | `"0"\|"1"\|"2"` | 0男 1女 2未知 |
| `avatar` | string | 头像地址 |
| `password` | string | 密码（仅写入；响应勿回传明文） |
| `status` | `"0"\|"1"` | 0正常 1停用 |
| `dept` | SysDept | 部门对象（列表展示） |
| `roles` | SysRole[] | 角色对象列表 |
| `roleIds` | number[] | 角色 ID 列表（提交用） |
| `postIds` | number[] | 岗位 ID 列表（提交用） |
| `loginIp` | string | 最后登录 IP（详情展示） |
| `loginDate` | string | 最后登录时间 |
| + BaseEntity | | `createdBy`/`createdAt`/`remark` 等 |

---

## GET `/system/user/list`

分页查询用户。

### Query

| 参数 | 说明 |
|------|------|
| `pageNum` / `pageSize` | 分页 |
| `userName` | 账号模糊 |
| `phonenumber` | 手机号 |
| `status` | `0`/`1` |
| `deptId` | 部门（含子部门时由后端约定） |
| `params[beginTime]` / `params[endTime]` | 创建时间范围 |
| `orderByColumn` / `isAsc` | 排序 |

### 响应

`TableDataInfo<SysUser>`：`{ code, msg, total, rows }`

---

## GET `/system/user/` 或 `/system/user/{userId}`

获取用户表单数据。`userId` 为空时表示新增初始化（仍返回全部角色/岗位选项）。

### 响应（UserFormDataResult）

```json
{
  "code": 200,
  "msg": "操作成功",
  "data": { "userId": 1, "userName": "admin", "...": "..." },
  "posts": [{ "postId": 1, "postName": "董事长" }],
  "roles": [{ "roleId": 1, "roleName": "超级管理员" }],
  "postIds": [1],
  "roleIds": [1]
}
```

| 字段 | 说明 |
|------|------|
| `data` | 用户详情；新增时可为 `null`/`{}` |
| `posts` | 全部岗位选项 |
| `roles` | 全部角色选项（可含 `flag` 表示是否已选） |
| `postIds` / `roleIds` | 当前用户已选 ID |

---

## POST `/system/user`

新增用户。Body：`SysUser`（含 `roleIds`、`postIds`、`password`）。

### 响应

`AjaxResult`

---

## PUT `/system/user`

修改用户。Body：`SysUser`。

---

## DELETE `/system/user/{userId}`

删除用户。`userId` 可为 `1,2,3`。

---

## PUT `/system/user/resetPwd`

管理员重置密码。

### 请求体

```json
{ "userId": 2, "password": "123456" }
```

---

## PUT `/system/user/changeStatus`

启停用户。

### 请求体

```json
{ "userId": 2, "status": "1" }
```

---

## GET `/system/user/profile`

个人中心资料。

### 响应（UserProfileResult）

```json
{
  "code": 200,
  "msg": "操作成功",
  "data": { "userId": 1, "userName": "admin", "nickName": "若依" },
  "roleGroup": "超级管理员",
  "postGroup": "董事长"
}
```

---

## PUT `/system/user/profile`

修改个人资料。Body：`SysUser`（昵称、手机、邮箱、性别等）。

---

## PUT `/system/user/profile/updatePwd`

修改自己的密码。

### 请求体

```json
{ "oldPassword": "old", "newPassword": "new" }
```

---

## POST `/system/user/profile/avatar`

上传头像。

- Header：`Content-Type: application/x-www-form-urlencoded`（前端 axios 配置；实际多为 FormData）
- Body：文件字段（常见 `avatarfile`）

### 响应

```json
{
  "code": 200,
  "msg": "操作成功",
  "imgUrl": "/profile/avatar/xxx.png"
}
```

---

## GET `/system/user/authRole/{userId}`

查询用户可分配角色。

### 响应

```json
{
  "code": 200,
  "user": { "userId": 2, "userName": "ry" },
  "roles": [
    { "roleId": 2, "roleName": "普通角色", "flag": true }
  ]
}
```

`flag=true` 表示该角色已分配给用户。

---

## PUT `/system/user/authRole`

保存用户角色。**参数在 Query（非 Body）。**

### Query

| 参数 | 说明 |
|------|------|
| `userId` | 用户 ID |
| `roleIds` | 角色 ID，逗号分隔字符串，如 `1,2` |

---

## GET `/system/user/deptTree`

用户管理页部门树下拉。

### 响应

```json
{
  "code": 200,
  "msg": "操作成功",
  "data": [
    { "id": 100, "label": "若依科技", "disabled": false, "children": [] }
  ]
}
```

---

## POST `/system/user/export`

导出 Excel。筛选条件同 list（form）。响应 blob。

---

## POST `/system/user/importData`

导入用户。`multipart`；Query：`updateSupport=0|1`（是否更新已存在用户）。

前端组件：`action="/system/user/importData"`。

---

## POST `/system/user/importTemplate`

下载导入模板。响应 xlsx blob。

---

## 关联配置

初始密码展示：`GET /system/config/configKey/sys.user.initPassword`，**值在响应 `msg` 字段**。
