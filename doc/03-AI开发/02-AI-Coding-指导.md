# Listen AI Coding 指导文档

## 0. 文档定位

本文件用于指导 Codex 或其他 AI Coding 工具在 Listen 项目中稳定产出代码。

核心原则：

- 面向 AI 可执行，不追求传统文档叙述感
- 降低上下文歧义
- 固定开发顺序
- 固定目录和命名
- 先复用现有实现，再扩展新模块
- 以最新原型 HTML 为执行主线

## 1. 强制原则

### 1.1 总原则

- 先读 `doc/03-AI开发/00-Codex-启动入口.md`
- 再读 `doc/03-AI开发/01-技术方案设计.md`
- 所有实现必须服从端优先级与功能优先级
- 若需求与 PRD 冲突，以“AI 主入口优先、第三方统一 mock、数据库优先”作为最高裁决规则
- 若仓库已有实现，优先沿用现有目录、接口、命名与状态模型
- 若最新原型与旧文档冲突，以 `doc/04-原型素材` 最新页面流为准

### 1.2 当前优先级

绝对优先：

1. 登录 / 资料补全 / 性格设置链路
2. 服务页与服务方详情页
3. 声音页
4. 我的页
5. 订单创建、订单详情、评价投诉骨架
6. 管理后台登录鉴权骨架
7. 管理后台服务项目 / 声音内容 / 订单投诉管理

低优先：

- 真实支付
- 结算
- 提现
- VIP

### 1.3 端优先级

1. 用户 Web
2. 管理后台 Web
3. 用户 App

## 2. AI 必须遵守的代码组织方式

### 2.1 Go 后端

当前实际目录：

```text
backend/internal/
  application/
    ai/
    provider/
  domain/
    ai/
    provider/
  infrastructure/
    ai/
    config/
    persistence/
      memory/
      mysql/
  interface/http/
    user/
    admin/
```

后续新增模块建议：

```text
backend/internal/
  application/
    identity/
    service/
    audio/
    order/
    feedback/
    payment/
  domain/
    identity/
    service/
    audio/
    order/
    feedback/
    payment/
```

强制要求：

- 新增模块继续放入 `interface/http`、`application`、`domain`、`infrastructure`
- Controller 不直接访问数据库
- SQL 不直接散落在 UseCase 中
- 状态流转写入实体方法或领域服务
- 每个接口有请求 DTO 和响应 DTO
- 运行配置优先从 `~/conf/listenbase.cof` 读取

### 2.2 Vue 前端

当前用户 Web 实际目录：

```text
web/apps/user-web/src/
  api/
  application/ai/
  components/ai/
  domain/ai/
  router/
  styles/
  views/
    home/
    chat/
```

当前管理后台实际目录：

```text
web/apps/admin-web/src/
  api/
  application/provider/
  components/providers/
  domain/provider/
  router/
  styles/
  views/
    dashboard/
    providers/
```

强制要求：

- 新页面继续按 `views / components / domain / application / api` 分层
- API 返回原始结构必须先经过 Adapter 转换
- 页面状态统一为 `idle/loading/success/empty/error`
- Mock 与 HTTP 适配器优先保持同构接口
- 原型页中的状态按钮、操作文案、流程跳转要尽量映射成明确状态或事件

禁止：

- 直接在页面里硬编码业务枚举
- 模板中出现复杂权限判断
- API 返回原始结构跨多个组件直接传播

### 2.3 面向对象约束

AI 生成代码时优先采用以下结构：

- `Entity`：带行为的实体
- `Repository`：仓储接口
- `UseCase`：应用服务
- `ViewModel`：页面状态编排
- `Api Adapter`：DTO 转换与请求封装

不要把所有逻辑压成：

- 巨大的 `utils.ts`
- 巨大的 `helpers.go`
- 巨大的单文件页面

## 3. 配置、数据库、Mock 的执行规则

### 3.1 配置规则

所有后端服务配置相关修改，统一遵守：

- 主配置来源：`~/conf/listenbase.cof`
- 仓库内配置仅允许作为兜底或示例
- 读取逻辑集中在 `infrastructure/config`
- 不要在业务代码中散落 `os.Getenv` 读取

### 3.2 数据库规则

涉及数据相关开发，尽量设计为从数据库读取。

允许使用：

```text
MYSQL_DSN=root:YUIO980-klio@tcp(192.168.2.24:13306)/listen?parseTime=true&charset=utf8mb4
```

规则：

- 列表、详情、字典、配置、订单、评价、投诉优先设计 MySQL 仓储
- 前端 mock 数据只为页面开发提速，不替代后端数据结构设计
- 若当前无法直连数据库，允许保留 mock / memory 验证，但 migration、repository interface、mysql repository 仍应预留

### 3.3 第三方交互 Mock 规则

所有涉及第三方交互的逻辑，当前统一先 mock，暂不进行真实对接。

包括：

- 首页推荐结果
- AI 对话推荐文案
- AI 对话自动回复
- 服务推荐理由
- 声音推荐理由
- 短信验证码发送
- 微信登录 / 微信授权
- 支付下单结果

规则：

- mock 统一收口到 `AiApi` 或后端 `MockMatchService` 风格的适配层
- 登录、支付等第三方交互也要有各自 adapter / gateway
- 业务页面与 UseCase 不直接拼装第三方返回结构

### 3.4 登录页默认值规则

登录页面默认填充：

- 手机号：`13800113800`
- 验证码：`123`

规则：

- 页面初始值直接展示这组默认值
- 发送验证码按钮仅做 mock 倒计时或 mock 成功提示
- 不真实发送短信

### 3.5 支付规则

支付开发统一排最后。

当前可以做：

- 支付确认页
- 订单创建
- 订单详情
- mock 支付成功
- 评价投诉链路

当前不要做：

- 真实支付下单
- 回调通知
- 对账
- 分账

实现方式：

- 后端提供 `mock-success` 风格支付结果接口即可
- 前端支付页点击确认后可直接进入订单详情页
- 不引入真实支付 SDK

## 4. 当前已实现模块的续写规则

### 4.1 用户 Web AI 模块

续写时优先复用：

- `AiApi`
- `HomePageViewModel`
- `ChatPageViewModel`
- `AiSession`
- `AiMessage`
- `MatchCandidate`
- `PageLoadState`

约束：

- 不要在页面里重新发明一套 AI 状态管理
- 若新增 AI 回复或推荐承接，应从 `AiApi` 与 `ChatPageViewModel` 扩展

### 4.2 管理后台服务方审核模块

续写时优先复用：

- `ProviderAdminApi`
- `ProviderReviewViewModel`
- `ProviderReviewStatus`
- `ProviderSummary`
- `ProviderDetail`

### 4.3 后端 AI 模块

续写时优先复用：

- `domain/ai`
- `application/ai`
- `interface/http/user`
- `infrastructure/persistence/{memory,mysql}`

### 4.4 新模块续写建议

建议新增独立模块而不是塞进现有 AI 模块：

- 登录 / 用户资料 / 性格设置：`identity`
- 服务页 / 服务方详情：`service`
- 声音页：`audio`
- 订单详情 / 状态机：`order`
- 评价 / 投诉：`feedback`
- 支付 mock：`payment`

## 5. 命名规范

### 5.1 Go 命名

- 实体：`Session`, `DailyQuota`, `Provider`, `User`, `Order`, `Review`, `Complaint`
- 仓储：`SessionRepository`, `UserRepository`, `OrderRepository`
- 用例：`CreateAiSessionUseCase`, `SaveUserProfileUseCase`, `CreateOrderUseCase`
- 服务：`MatchService`, `PaymentMockService`

### 5.2 TypeScript 命名

- 实体：`AiSession`, `AiMessage`, `MatchCandidate`, `UserProfile`, `OrderDetail`
- 页面模型：`HomePageViewModel`, `ChatPageViewModel`, `OrderDetailPageViewModel`
- API 适配器：`AiApi`, `IdentityApi`, `OrderApi`

### 5.3 枚举命名

- `PageLoadState`
- `ProviderReviewStatus`
- `OrderStatus`
- `ComplaintType`

枚举值用明确字符串，不用数字。

## 6. AI 开发任务拆分标准

### 6.1 合格任务粒度

一条 AI 开发任务应只覆盖下面之一：

- 一个页面及其依赖接口
- 一个领域模块的实体 + 仓储 + 用例
- 一个后台列表页 + 详情页
- 一个状态机
- 一组 migration 脚本
- 一个现有模块的小范围联调或补测

### 6.2 不合格任务粒度

不要给 AI 这种模糊任务：

- “把整个项目做完”
- “实现完整产品”
- “把登录、订单、支付一起写完”

## 7. 提交前检查

每次修改后，至少检查：

1. 是否沿用了现有目录结构
2. 是否把业务逻辑留在了 ViewModel / UseCase / Entity
3. 是否补了需要的 migration / mock / dto
4. 是否覆盖了 `idle/loading/success/empty/error`
5. 是否说明了启动与验证方式
6. 是否遵守了 `~/conf/listenbase.cof` 配置规则
7. 是否把可持久化数据按数据库读取方式设计
8. 是否把 AI / 支付的 mock 边界隔离清楚
