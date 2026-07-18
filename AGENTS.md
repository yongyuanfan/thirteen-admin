# AGENTS.md

本文件定义 OpenCode 在本仓库中开发 Go 项目时应遵循的通用规则。

## 基本原则

- 优先做最小且正确的改动，避免无关重构。
- 先理解现有代码结构和约定，再实施修改。
- 保持实现直接、可维护，避免引入不必要的抽象。
- 不得擅自修改与当前任务无关的代码、配置或目录结构。

## Go 代码约定

- 遵循标准 Go 风格，提交前保持 `gofmt` 格式化结果一致。
- 优先使用标准库；新增第三方依赖前，确认现有能力无法满足需求。
- 变量、函数、类型命名应清晰、语义化，避免无意义缩写。
- 错误处理要显式，不忽略错误；错误信息应包含足够上下文。
- 除非已有模式明确要求，否则不要过早抽象接口。
- 注释保持简洁，尽可能使用中文对变量、方法和代码块注释，仅解释不直观的设计意图、边界条件或复杂逻辑。

## 命名规则

- 尽量使用合理的名称对目录、文件、函数、方法和变量进行命名
- 尽可能使用简短的名称，但需要保证名称的合理性
- 命名尽量不要与同级或者同功能模块重复

## 项目结构

- `cmd/` 放置程序入口和启动相关代码。
- `app/` 放置应用 Bootstrap 与业务分层实现（扁平包 + 域前缀文件名，如 `auth.go`）：
  - `handlers/`：HTTP 处理器（入参校验、调用 service、写响应）；**不含路由注册**
  - `services/`：业务逻辑及支撑能力（如 JWT、验证码）
  - `models/`：GORM 数据模型
  - `seeds/`：AutoMigrate 与种子数据
  - `middlewares/`：Gin 中间件（如 JWT 鉴权）
  - `router/`：**唯一**放置路由注册；`router.go` 为注册表，业务路由如 `auth.go` 的 `RegisterAuthRoutes`
  - `app.go`：Bootstrap 装配依赖并 `InitRouter`
- `internal/base/` 分层通用代码（无具体业务域）：
  - `handler/`：统一响应信封 `{code,msg,data}`、错误码、`Handler` 基类
  - `service/`：`Service` 基类（持有 `*gorm.DB`）
  - `model/`：`BaseModel`、`StringSlice` 等
- `internal/` 其余为基础设施：`config`、`storage`、`logging`、`runtime`、`utils`；**不写业务域逻辑**。
- 禁止恢复 `app/<domain>/` 域包或把业务逻辑塞进 `internal/`（`base` 通用层除外）。
- 公共配置、静态资源等放在现有对应目录；新文件位置须与上述分层一致。
- 更细的 Cursor 会话规则见 `.cursor/rules/backend-architecture.mdc`。

## 分层依赖

允许的依赖方向：

```
router → handlers → services → models → internal/base/model
                ↘ middlewares → services
                ↘ internal/base/handler
seeds → models
Bootstrap(app) → handlers / services / seeds / middlewares / router
```

禁止：

- `handlers` import `router`
- `services` / `models` / `seeds` import `handlers`
- 业务包互相形成环；`internal/*`（除 `runtime` 调用 Bootstrap）依赖 `app/*`

## 新增业务模块步骤

1. 在 `app/models/` 增加模型（优先嵌入 `model.BaseModel`）。
2. 在 `app/services/` 实现业务（可嵌入 `service.Service`）。
3. 在 `app/handlers/` 实现 Handler（嵌入 `handler.Handler`，使用 `Success`/`Fail`）。
4. 在 `app/router/` 增加 `RegisterXxxRoutes`，通过 `RegisterFreeRouter` / `RegisterAuthRouter` 绑定。
5. 需要种子数据时在 `app/seeds/` 增加 AutoMigrate/Seed。
6. 在 `app/app.go` 的 `Bootstrap` 中完成构造、注册路由，最后 `InitRouter`。

## API 与响应约定

- 业务 API 使用根路径（如 `/login`、`/system/user/list`），与 `docs/api-docs` / thirteen-admin-ui 对齐；**无** `/api/v1` 前缀。
- 统一响应：`{ "code": number, "msg": string, "data"?: any }`；成功码 `200`；业务失败也返回 HTTP 200，由 `code` 区分（`401`/`500`/`601` 等）。
- 分页列表：`{ "code", "msg", "total", "rows" }`。
- 部分接口字段平铺根级（如登录 `token`、`getInfo` 的 `user/roles/permissions`）。
- JSON 字段使用 camelCase（如 `userName`、`userId`）；库表 `id`/`created_at` 经 VO 映射为 `userId`/`createdAt` 等。
- 需登录接口使用 `Authorization: Bearer <token>`；用户 ID 从 `middlewares.ContextUserIDKey` 读取。
- 接口契约真源：`docs/api-docs/`。

## 变更规则

- 修改前请先判断是否有变更的必要，避免无意义的变更
- 修改前先确认影响范围，避免破坏现有公开行为。
- 涉及跨文件改动时，保持接口、调用方和测试同步更新。
- 不为了“未来可能会用到”而提前加入兼容层、开关或复杂扩展点。
- 不擅自删除用户已有改动；发现冲突时优先保留并说明。

## 测试与验证

- 有对应测试位置时，优先补充或更新测试。
- 至少执行与改动直接相关的测试或验证命令。
- 若无法运行测试，应明确说明原因和未验证风险。
- 修复缺陷时，优先补充能覆盖该问题的回归测试。

## 配置与安全

- 不提交密钥、令牌、密码或其他敏感信息。
- 涉及配置变更时，优先复用现有配置方式，如 `.env`、配置文件或已有结构体。
- 对外部输入保持校验，对文件、网络、序列化操作做好错误处理。

## 性能与稳定性

- 在 I/O、并发、缓存、数据库或 RPC 路径上，优先保证正确性，再考虑性能优化。
- 新增 goroutine 时必须明确其生命周期、退出条件和错误传播方式。
- 避免引入潜在竞态、泄漏、阻塞或无界重试。

## 提交产出要求

- 输出应包含：修改了什么、为什么这样改、如何验证。
- 若存在未覆盖的边界情况、假设或后续建议，应明确列出。

## 禁止事项

- 不做与任务无关的大规模重命名或格式化。
- 不随意升级、替换或移除依赖。
- 不引入与当前仓库风格明显不一致的框架式写法。
- 不在未确认的情况下修改 CI、发布流程或生产配置。

## 必做事项

- 当变更完成后，必须完整的确认当前变更从可维护性，健壮性，简洁性，是否是最合理的，并且所实现的功能是完整的，如果不完整，需要给出后续执行计划


请确认当前变更从可维护性，健壮性，简洁性，是否是最合理的，并且所实现的功能是完整的，如果不完整，请给出后续执行计划