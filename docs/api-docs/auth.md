# 认证与路由

前端定义：`src/api/login.ts`、`src/api/menu.ts`  
类型：`src/types/api/login.ts`、`src/types/api/menu.ts`

---

## POST `/login`

用户名密码登录。**无需 Token。**

### 请求体

```json
{
  "username": "admin",
  "password": "admin123",
  "code": "12",
  "uuid": "验证码缓存 key"
}
```

| 字段 | 必填 | 说明 |
|------|------|------|
| `username` | 是 | 用户账号 |
| `password` | 是 | 密码（明文传输；HTTPS 部署） |
| `code` | 视验证码开关 | 图形验证码 |
| `uuid` | 视验证码开关 | 与 `/captchaImage` 返回的 uuid 对应 |

当 `captchaEnabled=false` 时，`code`/`uuid` 可放宽校验。

### 响应

```json
{
  "code": 200,
  "msg": "操作成功",
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

`token` 在根级（非 `data` 内）。

### curl

```bash
curl -s -X POST 'http://localhost:8080/login' \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123","code":"1","uuid":"xxx"}'
```

---

## POST `/register`

用户自助注册。**无需 Token。**  
前端还有 `confirmPassword` 仅做本地校验，不提交后端。

### 请求体

```json
{
  "username": "newuser",
  "password": "pass123",
  "code": "12",
  "uuid": "..."
}
```

### 响应

标准 `AjaxResult`。

---

## GET `/captchaImage`

获取图形验证码。**无需 Token。** 前端超时 20s。

### 响应

```json
{
  "code": 200,
  "msg": "操作成功",
  "uuid": "a1b2c3",
  "img": "data:image/gif;base64,... 或纯 base64",
  "captchaEnabled": true
}
```

| 字段 | 说明 |
|------|------|
| `uuid` | 验证码缓存 key，登录时回传 |
| `img` | 验证码图片 Base64 |
| `captchaEnabled` | 是否启用验证码；`false` 时前端可不展示验证码 |

---

## GET `/getInfo`

获取当前登录用户信息、角色与权限。**需要 Token。**

### 响应

```json
{
  "code": 200,
  "msg": "操作成功",
  "user": {
    "userId": 1,
    "userName": "admin",
    "nickName": "管理员",
    "avatar": "",
    "deptId": 103,
    "email": "ry@163.com",
    "phonenumber": "15888888888",
    "sex": "1",
    "status": "0"
  },
  "roles": ["admin"],
  "permissions": ["*:*:*"],
  "isDefaultModifyPwd": false,
  "isPasswordExpired": false,
  "pwdChrtype": "0"
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `user` | object | `SysUser`，见 [system/user.md](./system/user.md) |
| `roles` | string[] | 角色权限字符（如 `admin`、`common`），非数字 ID |
| `permissions` | string[] | 权限标识；超管可为 `["*:*:*"]` |
| `isDefaultModifyPwd` | boolean | 是否需提醒修改初始密码 |
| `isPasswordExpired` | boolean | 密码是否过期 |
| `pwdChrtype` | string | 密码字符策略 `0`–`4`（前端校验用，可选） |

---

## POST `/logout`

退出登录。**需要 Token。**

### 响应

标准 `AjaxResult`。后端应使当前 Token / 会话失效。

---

## POST `/unlockscreen`

锁屏后输入密码解锁。**需要 Token。**

### 请求体

```json
{ "password": "admin123" }
```

### 响应

标准 `AjaxResult`；密码错误返回非 200。

---

## GET `/getRouters`

按当前用户权限返回前端动态路由树。**需要 Token。**

前端：`store/modules/permission.ts`

### 响应

```json
{
  "code": 200,
  "msg": "操作成功",
  "data": [
    {
      "name": "System",
      "path": "/system",
      "hidden": false,
      "redirect": "noRedirect",
      "component": "Layout",
      "alwaysShow": true,
      "meta": {
        "title": "系统管理",
        "icon": "system",
        "noCache": false,
        "link": null
      },
      "children": [
        {
          "name": "User",
          "path": "user",
          "hidden": false,
          "component": "system/user/index",
          "meta": {
            "title": "用户管理",
            "icon": "user",
            "noCache": false
          }
        }
      ]
    }
  ]
}
```

### RouterVo 字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `name` | string | 路由名称 |
| `path` | string | 路由路径 |
| `hidden` | boolean | 是否在侧栏隐藏 |
| `redirect` | string | 重定向；`noRedirect` 时面包屑不可点 |
| `component` | string | 组件：`Layout` / `ParentView` / `InnerLink` / 业务路径如 `system/user/index` |
| `query` | string | 路由参数 JSON 字符串 |
| `alwaysShow` | boolean | 始终显示根菜单 |
| `meta.title` | string | 侧栏/面包屑标题 |
| `meta.icon` | string | 图标名 |
| `meta.noCache` | boolean | 是否不缓存 |
| `meta.link` | string | 内链地址 |
| `children` | RouterVo[] | 子路由 |

仅返回当前用户有权限的菜单（类型 M/C）；按钮（F）通过 `permissions` 控制，不进入路由树。
