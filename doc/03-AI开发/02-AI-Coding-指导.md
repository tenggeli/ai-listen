# Listen AI Coding 指导文档（按 2026-04-02 当前仓库状态）

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
   - `doc/05-管理后台原型/*.html`
4. 再开始编码

若文档与代码冲突：以代码现状为准，并回写文档。

## 2. 当前优先级（2026-04-02 代码快照）

高优先：

1. 服务方管理后台前端工程骨架与工作台
2. 平台管理后台运营模块扩展（声音内容、订单、投诉）
3. AI 对话自动回复后端化
4. 声音页数据从 mock service 逐步迁移为可配置内容源
5. 订单状态流转扩展（接单、履约、完成、售后）

中优先：

- 平台与服务方双后台数据联动与运营指标
- 用户 App 页面化（优先首页/对话/声音）

低优先：

- 真实 AI 网关
- 真实短信与微信授权
- 真实支付、结算、提现、VIP
- 用户 App 页面化

## 3. 必须遵守的代码组织

### 3.1 Go 后端

当前实际目录核心：

```text
backend/internal/
  application/
    admin_auth/
    ai/
    audio/
    feedback/
    identity/
    order/
    provider/
    service_discovery/
    service_item_admin/
    user_settings/
  domain/
    admin_auth/
    ai/
    audio/
    feedback/
    identity/
    order/
    provider/
    service_discovery/
    service_item_admin/
    user_settings/
  infrastructure/
    ai/
    audio/
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

后端架构约束：

- 后端保持一个统一接口服务，不按前端后台数量拆成多个后端项目
- 可按路由分组拆 controller，例如 `user`、`admin`、`provider`
- 管理侧和服务方侧共享领域能力时，优先共享 application/domain，而不是复制一份后端

新增接口命名强制规则：

- 用户侧前缀固定：`/api/v1/...`
- 服务方侧前缀固定：`/api/v1/provider/...`
- 平台管理侧前缀固定：`/api/v1/admin/...`
- 不要把服务方接口写进 `/api/v1/admin/...`
- 不要把平台管理接口写成普通用户侧 `/api/v1/...`
- 不要使用 `get/list/create/update` 这类动词作为一级路径
- 动作型接口统一放在资源详情后，例如 `/orders/{id}/accept`
- 资源路径统一用小写英文复数名词，segment 之间用短横线
- 单例资源允许使用单数名词，例如 `profile`

历史接口例外：

- 当前已落地接口不要求按本规范回溯重命名
- 重点例外包括：`/api/v1/auth/login/*`、`/api/v1/users/me/*`、`/api/v1/ai/match/remaining`
- `GET /healthz` 视为基础设施接口，不纳入业务路由分组

强制要求：

- Controller 不直连数据库
- SQL 不散落在 UseCase
- 状态流转放到 Entity 或 DomainService
- 每个接口都定义请求/响应 DTO
- 配置读取统一走 `infrastructure/config`
- 新增接口前先确认其归属端，再决定是否落在 `user/admin/provider` 路由组

### 3.2 Vue 前端

当前用户 Web 已有：

- `views/auth`
- `views/profile`
- `views/home`
- `views/chat`
- `views/sound`
- `views/services`
- `views/me`
- `views/payment`
- `views/order`
- `views/settings`

新增页面保持：`views / components / application / domain / api`。

双后台约束：

- 平台管理后台与服务方管理后台是两个独立前端，不共用一个页面壳
- 当前 `web/apps/admin-web` 视为“平台管理后台”
- 服务方管理后台落地时应新建独立前端工程，优先参考 `doc/05-管理后台原型`

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
- `MyPageViewModel`
- `AuthSession`

注意：

- 登录默认值保持 `13800113800 / 123`
- 路由守卫依赖 session 的 `profileCompleted`、`personalityCompleted`
- identity 已支持 MySQL 仓储，不要再按“仅 memory”假设开发

### 4.2 用户 AI 模块（已落地）

续写优先复用：

- `AiApi`
- `HomePageViewModel`
- `ChatPageViewModel`
- `PageLoadState`

### 4.3 服务浏览模块（已落地）

续写优先复用：

- `ServiceDiscoveryApi`
- `ServicesPageViewModel`
- `ProviderDetailPageViewModel`
- `ServiceCategory`
- `ProviderPublicProfile`
- `ServiceItem`

注意：

- 服务链路 HTTP 契约已存在，优先兼容现有接口，不要重命名路由
- 当前服务详情页已把下单入口指向 `/payment/confirm`

### 4.4 声音模块（前后端已打通）

续写优先复用：

- `SoundApi`
- `SoundPageViewModel`
- `SoundPageAggregate`

注意：

- 现有 `HttpSoundApi` 请求 `GET /api/v1/sounds?page=home&user_id=...`
- 后端目前只支持 `page=home`
- 若扩展声音内容来源，优先保持既有响应结构稳定

### 4.5 支付、订单与评价链路（前后端已打通最小闭环）

续写优先复用：

- `PaymentConfirmPage`
- `HttpOrderApi`
- `HttpFeedbackApi`

注意：

- 当前支付确认页会先创建真实订单，再调用 `/api/v1/orders/{id}/pay/mock-success`
- 订单列表、订单详情、评价/投诉页已接真实后端接口
- 支付仍是 mock-success，扩展真实支付时优先保持订单接口与页面流程稳定

### 4.6 平台管理后台审核模块（已落地）

续写优先复用：

- `ProviderAdminApi`
- `ProviderReviewViewModel`
- `ProviderReviewStatus`
- `ProviderSummary` / `ProviderDetail`

### 4.7 服务方管理后台（待落地）

实现时优先参考：

- `doc/02-PRD/服务方端PRD.md`
- `doc/05-管理后台原型/listen_provider_console.html`

注意：

- 服务方管理后台是独立前端，不要直接在 `web/apps/admin-web` 里混做
- 后端接口仍放在同一个 Go 服务内，通过独立路由分组承接

### 4.8 用户设置模块（已落地）

续写优先复用：

- `HttpSettingsApi`
- `SettingsPage`
- `user_settings` 相关 UseCase/Repository

注意：

- 用户设置接口已落地：`GET/PUT /api/v1/users/me/settings`
- mysql 模式走后端持久化，memory 模式由后端返回 `501`，前端本地兜底

### 4.9 平台管理后台鉴权与服务项目管理（已落地）

续写优先复用：

- `HttpAdminAuthApi`
- `AuthService` / `AuthSessionStore`
- `HttpServiceItemAdminApi`
- `ServiceItemManageViewModel`

注意：

- 鉴权接口：`POST /api/v1/admin/auth/login/mock`、`GET /api/v1/admin/auth/me`
- 服务项目接口：`GET /api/v1/admin/service-items`、`GET /api/v1/admin/service-items/{id}`、`POST /api/v1/admin/service-items/{id}/activate|deactivate`

## 5. 配置、数据库、Mock 规则

### 5.1 配置

- 主配置来源：`~/conf/listenbase.cof`
- 示例配置：`backend/config/listenbase.example.cof`
- 不在业务代码散落 `os.Getenv`

### 5.2 数据库

- `ai/provider/service_discovery/identity/order/feedback/user_settings/service_item_admin` 已有 migration + mysql repository
- 新模块（payment、后台配置等）默认按 MySQL 设计
- 开发 mysql 模式时，应优先复用现有 `mysql.NewDB(...)` 与仓储注册方式

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

- “实现订单实体状态机 + 创建订单接口 + 单测”
- “为平台管理后台补声音内容管理最小接口与页面”
- “实现评价/投诉最小接口与前端提交页”

不合格任务示例：

- “把整个项目都做完”
- “登录、订单、支付一次性全部完成”
- “把服务方后台和平台管理后台混在一个前端里一起做”

## 7. 提交前检查清单

1. 是否复用现有目录结构与命名
2. 是否把逻辑放在 ViewModel/UseCase/Entity
3. 是否补齐 DTO、migration、repository（如需要）
4. 页面状态是否覆盖 `idle/loading/success/empty/error`
5. 是否遵守 `~/conf/listenbase.cof` 规则
6. 是否把第三方交互隔离在 mock adapter/gateway
7. 是否说明了验证方式与剩余未完成项
8. 若改动服务/声音/身份接口，是否兼容当前已上线前端契约
