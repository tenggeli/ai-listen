# Listen AI Coding 指导文档（按当前仓库状态）

## 0. 文档定位

本文件用于指导 Codex 或其他 AI Coding 工具在 Listen 项目中稳定产出代码。

核心原则：

- 先看代码现状，再写代码
- 在现有模块续写，不平地重构
- 固定分层、固定命名、固定状态模型
- 对第三方交互统一 mock

## 1. 强制执行顺序

1. 先读：`doc/03-AI开发/00-Codex-启动入口.md`
2. 再读：`doc/03-AI开发/01-技术方案设计.md`
3. 核对路由与接口：
   - `web/apps/user-web/src/router/index.ts`
   - `backend/internal/interface/http/user/router.go`
   - `backend/internal/interface/http/admin/router.go`
4. 再开始编码

若文档与代码冲突：以代码现状为准，并回写文档。

## 2. 当前优先级（2026-03-28 代码快照）

高优先：

1. 用户 Web 服务页 + 服务方详情页
2. Go 后端服务浏览模块（service_discovery）
3. Go 后端声音模块（audio_content），打通 `/api/v1/sounds`
4. 用户 Web 我的页
5. 订单/支付确认（mock-success）/订单详情/评价投诉
6. 管理后台登录鉴权骨架 + 运营模块扩展

低优先：

- 真实 AI 网关
- 真实短信与微信授权
- 真实支付、结算、提现、VIP

## 3. 必须遵守的代码组织

### 3.1 Go 后端

当前实际目录核心：

```text
backend/internal/
  application/
    ai/
    identity/
    provider/
  domain/
    ai/
    identity/
    provider/
  infrastructure/
    ai/
    config/
    identity/
    persistence/
      memory/
      mysql/
  interface/http/
    user/
    admin/
```

新增模块必须继续按 `interface/application/domain/infrastructure` 分层。

强制要求：

- Controller 不直连数据库
- SQL 不散落在 UseCase
- 状态流转放到 Entity 或 DomainService
- 每个接口都定义请求/响应 DTO
- 配置读取统一走 `infrastructure/config`

### 3.2 Vue 前端

当前用户 Web 已有：

- `views/auth`
- `views/profile`
- `views/home`
- `views/chat`
- `views/sound`

新增页面保持：`views / components / application / domain / api`。

强制要求：

- API 原始结构先经 Adapter 转换
- 页面状态统一 `idle/loading/success/empty/error`
- Mock 与 HTTP Adapter 保持同构接口
- 复杂业务逻辑不写在模板里

禁止：

- 在页面直接硬编码业务枚举
- 用全局 `utils` 堆业务流程
- API 原始对象跨组件透传

## 4. 当前模块续写规则

### 4.1 用户身份链路（已落地）

续写优先复用：

- `AuthApi`
- `LoginPageViewModel`
- `ProfileSetupViewModel`
- `PersonalitySetupViewModel`
- `AuthSession`

注意：

- 登录默认值保持 `13800113800 / 123`
- 路由守卫依赖 session 的 `profileCompleted`、`personalityCompleted`

### 4.2 用户 AI 模块（已落地）

续写优先复用：

- `AiApi`
- `HomePageViewModel`
- `ChatPageViewModel`
- `PageLoadState`

### 4.3 声音模块（前端已落地，后端未打通）

续写优先复用：

- `SoundApi`
- `SoundPageViewModel`
- `SoundPageAggregate`

注意：

- 现有 `HttpSoundApi` 请求 `GET /api/v1/sounds?page=home&user_id=...`
- 后端暂未提供该接口，补后端时不要改动前端已用契约

### 4.4 管理后台审核模块（已落地）

续写优先复用：

- `ProviderAdminApi`
- `ProviderReviewViewModel`
- `ProviderReviewStatus`
- `ProviderSummary` / `ProviderDetail`

## 5. 配置、数据库、Mock 规则

### 5.1 配置

- 主配置来源：`~/conf/listenbase.cof`
- 示例配置：`backend/config/listenbase.example.cof`
- 不在业务代码散落 `os.Getenv`

### 5.2 数据库

- `ai/provider` 已有 migration + mysql repository
- `identity` 目前仍是 memory repository，若做 mysql 落地需补 migration/repository
- 新模块（service/audio/order/feedback/payment）默认按 MySQL 读取设计

### 5.3 第三方交互

统一 mock：

- AI 网关
- 短信
- 微信
- 支付

要求：

- mock 收口 adapter/gateway
- 页面和 UseCase 不直接感知第三方细节

## 6. 任务粒度标准

合格任务示例：

- “实现服务页 + 服务列表接口（含 DTO 与 ViewModel）”
- “实现声音首页接口 `/api/v1/sounds` 并打通 `HttpSoundApi`”
- “实现订单实体状态机 + 创建订单接口 + 单测”

不合格任务示例：

- “把整个项目都做完”
- “登录、订单、支付一次性全部完成”

## 7. 提交前检查清单

1. 是否复用现有目录结构与命名
2. 是否把逻辑放在 ViewModel/UseCase/Entity
3. 是否补齐 DTO、migration、repository（如需要）
4. 页面状态是否覆盖 `idle/loading/success/empty/error`
5. 是否遵守 `~/conf/listenbase.cof` 规则
6. 是否把第三方交互隔离在 mock adapter/gateway
7. 是否说明了验证方式与剩余未完成项
