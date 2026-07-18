# thirteen-admin HTTP 接口文档

> 本文档依据 **thirteen-admin-ui**（若依 Vue3 + Element Plus）实际调用契约整理，作为后端实现与联调的权威参考。  
> 源码对照：`thirteen-admin-ui/src/api/**`、`thirteen-admin-ui/src/types/api/**`。

## 文档索引

| 文档 | 说明 |
|------|------|
| [conventions.md](./conventions.md) | 鉴权、响应信封、分页、导出、通用约定 |
| [endpoints.md](./endpoints.md) | 全部接口 Method + Path 速查表 |
| [auth.md](./auth.md) | 登录 / 注册 / 验证码 / 用户信息 / 路由 |
| [common.md](./common.md) | 文件上传 / 下载 |
| [system/user.md](./system/user.md) | 用户管理、个人中心 |
| [system/role.md](./system/role.md) | 角色管理、授权、数据权限 |
| [system/menu.md](./system/menu.md) | 菜单管理 |
| [system/dept.md](./system/dept.md) | 部门管理 |
| [system/post.md](./system/post.md) | 岗位管理 |
| [system/dict.md](./system/dict.md) | 字典类型与字典数据 |
| [system/config.md](./system/config.md) | 参数配置 |
| [system/notice.md](./system/notice.md) | 通知公告、已读 |
| [monitor/job.md](./monitor/job.md) | 定时任务与任务日志 |
| [monitor/operlog.md](./monitor/operlog.md) | 操作日志 |
| [monitor/logininfor.md](./monitor/logininfor.md) | 登录日志 |
| [monitor/online.md](./monitor/online.md) | 在线用户 |
| [monitor/cache.md](./monitor/cache.md) | 缓存监控 |
| [monitor/server.md](./monitor/server.md) | 服务监控 |
| [tool/gen.md](./tool/gen.md) | 代码生成 |

## 环境与前缀

| 前端环境变量 | 值 | 说明 |
|--------------|-----|------|
| `VITE_APP_BASE_API`（开发） | `/dev-api` | Vite 代理到后端，**去掉此前缀**后转发 |
| 生产 | `/prod-api` | 由网关/Nginx 转发 |

文档中的路径均为**去掉环境前缀后**的后端真实路径，例如：

- 前端请求：`POST /dev-api/login`
- 后端实现：`POST /login`

## 快速对照：端点数量概览

| 域 | 约计 |
|----|------|
| 认证与路由 | 8 |
| 公共上传下载 | 3 |
| 系统管理 | ~70 |
| 监控 | ~30 |
| 代码生成 | ~11 |
| **合计** | **~120+**（含 export / import） |

## 嵌入页（非 JSON API）

前端部分页面以 iframe 打开，不走 axios 业务信封：

| 路径 | 用途 |
|------|------|
| `/swagger-ui/index.html` | API 文档页（可选） |
| `/druid/login.html` | 数据源监控（可选） |
| `/v3/api-docs/*` | OpenAPI（可选，开发代理已配置） |

## 数据库

表结构见 [../ry.sql](../ry.sql)。

## 实现状态说明

当前 `thirteen-admin` 后端仅实现部分鉴权接口，且路径/信封可能与本文档不完全一致。  
**本文档描述的是前端所需的完整目标契约**；实现时请逐步对齐本文档，而非现有临时路由。
