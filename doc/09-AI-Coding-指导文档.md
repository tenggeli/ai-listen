# Listen AI Coding 指导文档

## 0. 文档定位

本文件用于指导 Codex 或其他 AI Coding 工具在 Listen 项目中稳定产出代码。

核心原则：

- 面向 AI 可执行，不追求传统文档叙述感
- 降低上下文歧义
- 固定开发顺序
- 固定目录和命名
- 固定任务颗粒度

## 1. 强制原则

## 1.1 总原则

- 先读 `doc/08-面向AI开发的技术方案设计.md`
- 所有实现必须服从端优先级与功能优先级
- 若需求与 PRD 冲突，以“AI 主入口优先、支付最后”作为最高裁决规则
- 若发现产品形态与服务方独立端冲突，优先实现“后台可管 + 用户端可消费”，服务方独立端延后

## 1.2 需求优先级

绝对优先：

1. 首页 AI 主入口
2. AI 对话
3. AI 推荐结果承接
4. 服务浏览
5. 声音内容
6. 管理后台审核与内容管理

低优先：

- 订单深流程
- 支付
- 结算
- 提现
- VIP

## 1.3 端优先级

1. 用户 Web
2. 管理后台 Web
3. 用户 App

## 2. AI 必须遵守的代码组织方式

## 2.1 Go 后端

强制分层：

- `interface/http`
- `application`
- `domain`
- `infrastructure`

禁止：

- Controller 直接访问数据库
- SQL 直接散落在业务流程代码里
- 把状态判断写在多个 controller 中重复实现

必须：

- 每个核心模块有 `entity.go`, `repository.go`, `service.go` 或 `usecase.go`
- 每个状态机有枚举定义
- 每个接口有请求 DTO 和响应 DTO

## 2.2 Vue 前端

强制分层：

- `views/pages`
- `components`
- `domain`
- `application`
- `api`

禁止：

- 直接在页面里硬编码业务枚举
- 模板中出现复杂权限判断
- API 返回原始结构直接在 3 个以上组件中传播

必须：

- 定义页面模型 `ViewModel`
- 定义领域实体 `Entity`
- 定义 API Adapter
- 页面状态统一为 `idle/loading/success/empty/error`

## 2.3 面向对象约束

AI 生成代码时应尽量优先以下结构：

- `Entity`：带行为的实体
- `ValueObject`：不可变值对象
- `Repository`：仓储接口
- `UseCase`：应用服务
- `Factory`：复杂对象构建
- `Assembler`：DTO 与 Entity 转换

不要把所有逻辑压成：

- 巨大的 `utils.ts`
- 巨大的 `helpers.go`
- 巨大的单文件页面

## 3. 推荐目录模板

## 3.1 用户 Web

```text
web/apps/user-web/src/
  main.ts
  router/
  views/
    home/
    chat/
    services/
    sound/
    me/
    orders/
  components/
    ai/
    service/
    audio/
    common/
  domain/
    ai/
    service/
    audio/
    user/
    order/
  application/
    ai/
    service/
    audio/
    user/
    order/
  api/
  styles/
```

## 3.2 管理后台 Web

```text
web/apps/admin-web/src/
  views/
    login/
    dashboard/
    providers/
    service-items/
    audio/
    orders/
    complaints/
    configs/
  components/
  domain/
  application/
  api/
```

## 3.3 Go 后端

```text
backend/internal/
  interface/http/
  application/
  domain/
  infrastructure/
```

## 4. 命名规范

## 4.1 Go 命名

- 实体：`User`, `AiSession`, `Provider`, `Order`
- 仓储：`UserRepository`, `OrderRepository`
- 用例：`CreateAiSessionUseCase`, `ListProvidersUseCase`
- 服务：`MatchService`, `ComplaintResolutionService`

## 4.2 TypeScript 命名

- 实体：`AiConversationSession`, `MatchCandidate`, `AudioTrack`
- 页面模型：`HomePageViewModel`, `OrderDetailViewModel`
- API 适配器：`AiApi`, `ProviderApi`
- 应用服务：`AiChatService`, `ProviderQueryService`

## 4.3 枚举命名

- `OrderStatus`
- `ProviderReviewStatus`
- `AiMatchState`
- `PageLoadState`

枚举值用明确字符串，不用数字。

## 5. AI 开发任务拆分标准

## 5.1 合格任务粒度

一条 AI 开发任务应只覆盖下面之一：

- 一个页面及其依赖接口
- 一个领域模块的实体 + 仓储 + 用例
- 一个后台列表页 + 详情页
- 一个状态机
- 一组迁移脚本

## 5.2 不合格任务粒度

不要给 AI 这种模糊任务：

- “把整个项目做完”
- “实现完整产品”
- “把订单支付都补齐”

应改成：

- “实现用户 Web 首页 AI 主入口页面与匹配次数接口”
- “实现后台服务方审核列表、详情、审核通过/拒绝接口”
- “实现声音分类、列表、播放记录表与接口”

## 5.3 推荐开发顺序

### 阶段 A：基础骨架

- 初始化前后端工程
- 建立目录
- 建立公共响应格式
- 建立配置与路由

### 阶段 B：AI 主入口

- 首页
- AI 会话
- 流式输出
- 匹配次数限制
- 匹配结果承接

### 阶段 C：内容与服务承接

- 服务页
- 服务方详情
- 声音页
- 我的页基础

### 阶段 D：后台运营能力

- 登录
- 服务方审核
- 服务项目管理
- 声音内容管理
- 平台配置

### 阶段 E：交易骨架

- 订单模型
- 订单状态机
- 订单列表与详情
- 评价与投诉

### 阶段 F：低优先商业化

- 支付
- 结算
- 提现
- VIP

## 6. AI 写代码时的页面规则

## 6.1 用户首页

必须有：

- 核心呼吸球
- 文本输入入口
- 语音输入入口占位
- 推荐 3 人承接区
- 高级筛选入口
- AI 次数限制提示

禁止：

- 做成传统信息流首页
- 堆砌过多图标入口
- 让“服务列表”压过 AI 主入口

## 6.2 AI 对话页

必须有：

- 会话流
- AI 流式输出
- 快捷回复
- 跳转卡片
- 安全兜底文案

## 6.3 服务页

必须有：

- 分类
- 筛选
- 服务方卡片
- 详情

P0 可不做：

- 真支付
- 完整订单联动

## 6.4 声音页

必须有：

- 分类 tabs
- 播放中卡片
- 推荐列表
- 播放记录接口占位

## 6.5 管理后台

必须先做：

- 服务方审核
- 服务项目管理
- 声音内容管理
- 订单/投诉只做结构骨架也可以

## 7. AI 写代码时的后端规则

## 7.1 Controller 规则

- 只负责参数接收、校验、调用 use case、返回 response
- 不能承载复杂状态迁移逻辑

## 7.2 UseCase 规则

- 一次只完成一个业务动作
- 必须有输入对象和输出对象
- 必须明确依赖哪些 repository / service

## 7.3 Repository 规则

- 接口先行
- MySQL 实现后补
- 查询对象要明确，不要无限参数列表

## 7.4 状态机规则

- 状态迁移必须集中实现
- 非法迁移直接报业务错误
- 每次迁移都写状态日志

## 8. AI 写代码时的数据库规则

- 先写 migration，再写 repository
- 不在 SQL 中编码复杂业务状态机
- 关键快照字段直接冗余
- 所有列表接口默认分页
- 所有外键字段都建立索引

## 9. AI 测试与验收规则

## 9.1 后端最低测试要求

- 状态机单测
- use case 单测
- controller 基本接口测试

## 9.2 前端最低验收要求

- 页面能渲染
- 空态、加载态、错误态齐全
- API mock 下可完成主路径交互

## 9.3 每次任务完成时，AI 必须自检

- 是否符合当前阶段优先级
- 是否引入了支付等低优先复杂度
- 是否把业务逻辑写进 UI
- 是否新增了重复枚举或重复 DTO
- 是否保留了后续扩展点

## 10. 推荐 Prompt 模板

## 10.1 实现页面时

```text
请基于 doc/08-面向AI开发的技术方案设计.md 与 doc/09-AI-Coding-指导文档.md，
在不实现支付的前提下，完成用户 Web 的首页 AI 主入口。
要求：
1. 使用 Vue
2. 页面包含核心呼吸球、文本输入、推荐结果承接区、高级筛选入口
3. 接 API 契约但允许先 mock
4. 页面状态包含 idle/loading/success/error/limited
5. 抽离 domain model、view model、api adapter
```

## 10.2 实现后端模块时

```text
请基于 doc/08-面向AI开发的技术方案设计.md 与 doc/09-AI-Coding-指导文档.md，
实现 Go 后端 ai_companion 上下文的最小闭环。
要求：
1. 先定义 entity、repository interface、use case、controller、dto
2. 提供创建会话、发送消息、查询剩余匹配次数接口
3. AI 网关先用 mock 实现
4. 匹配上限为单用户每日 5 次
5. 给出 migration 与基本单测
```

## 10.3 实现后台时

```text
请基于 doc/08-面向AI开发的技术方案设计.md 与 doc/09-AI-Coding-指导文档.md，
实现管理后台的服务方审核模块。
要求：
1. 包含列表、详情、审核通过、审核拒绝、补充资料
2. 前后端分层明确
3. 所有审核状态使用明确枚举
4. 生成最小可运行页面与 API
```

## 11. 常见错误清单

- 把首页做成普通社交信息流
- 先做支付，导致主入口迟迟不可用
- 没有把 AI 输出结构化
- 把服务方做成第四端，偏离当前产品形态
- 在前端模板里写复杂业务判断
- 在 controller 里直接写状态流转
- 没有状态日志
- 没有空态和错误态

## 12. 当前阶段的最优开发入口

如果 AI 只能从一个任务开始，必须从这里开始：

1. 用户 Web 首页 AI 主入口
2. Go 后端 AI 会话与匹配次数接口
3. 管理后台服务方审核列表

这是当前最符合产品价值、最容易建立整体骨架、也最适合 AI 连续开发的起点。
