# Listen AI Coding 指导文档

## 0. 文档定位

本文件用于指导 Codex 或其他 AI Coding 工具在 Listen 项目中稳定产出代码。

核心原则：

- 面向 AI 可执行，不追求传统文档叙述感
- 降低上下文歧义
- 固定开发顺序
- 固定目录和命名
- 先复用现有实现，再扩展新模块

## 1. 强制原则

## 1.1 总原则

- 先读 `doc/03-AI开发/00-Codex-启动入口.md`
- 再读 `doc/03-AI开发/01-技术方案设计.md`
- 所有实现必须服从端优先级与功能优先级
- 若需求与 PRD 冲突，以“AI 主入口优先、支付最后”作为最高裁决规则
- 若仓库已有实现，优先沿用现有目录、接口、命名与状态模型

## 1.2 当前优先级

绝对优先：

1. 用户 Web 首页与 AI 对话真实联调
2. 服务浏览页
3. 声音页
4. 管理后台登录鉴权骨架
5. 管理后台服务项目 / 声音内容管理

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

强制要求：

- 新增模块继续放入 `interface/http`、`application`、`domain`、`infrastructure`
- Controller 不直接访问数据库
- SQL 不直接散落在 UseCase 中
- 状态流转写入实体方法或领域服务
- 每个接口有请求 DTO 和响应 DTO

## 2.2 Vue 前端

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

禁止：

- 直接在页面里硬编码业务枚举
- 模板中出现复杂权限判断
- API 返回原始结构跨多个组件直接传播

## 2.3 面向对象约束

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

## 3. 当前已实现模块的续写规则

## 3.1 用户 Web AI 模块

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
- 若新增 AI 回复或流式能力，应从 `AiApi` 与 `ChatPageViewModel` 扩展

## 3.2 管理后台服务方审核模块

续写时优先复用：

- `ProviderAdminApi`
- `ProviderReviewViewModel`
- `ProviderReviewStatus`
- `ProviderSummary`
- `ProviderDetail`

约束：

- 审核状态必须继续收敛为明确枚举字符串
- 列表与详情变更时，优先走 ViewModel 状态同步

## 3.3 后端 AI 模块

续写时优先复用：

- `domain/ai`
- `application/ai`
- `interface/http/user`
- `infrastructure/persistence/{memory,mysql}`

约束：

- 新能力先补齐仓储接口，再补 UseCase，再补 Controller
- 若新增表结构，需要同步补 migration
- 若修改 AI 用例，优先补充 `usecase_test.go`

## 4. 命名规范

## 4.1 Go 命名

- 实体：`Session`, `DailyQuota`, `Provider`
- 仓储：`SessionRepository`, `MatchQuotaRepository`
- 用例：`CreateAiSessionUseCase`, `ReviewProviderUseCase`
- 服务：`MatchService`

## 4.2 TypeScript 命名

- 实体：`AiSession`, `AiMessage`, `MatchCandidate`, `ProviderSummary`
- 页面模型：`HomePageViewModel`, `ChatPageViewModel`, `ProviderReviewViewModel`
- API 适配器：`AiApi`, `ProviderAdminApi`

## 4.3 枚举命名

- `PageLoadState`
- `ProviderReviewStatus`

枚举值用明确字符串，不用数字。

## 5. AI 开发任务拆分标准

## 5.1 合格任务粒度

一条 AI 开发任务应只覆盖下面之一：

- 一个页面及其依赖接口
- 一个领域模块的实体 + 仓储 + 用例
- 一个后台列表页 + 详情页
- 一个状态机
- 一组迁移脚本
- 一个现有模块的小范围联调或补测

## 5.2 不合格任务粒度

不要给 AI 这种模糊任务：

- “把整个项目做完”
- “实现完整产品”
- “把所有 P0 一次写完”

## 6. 提交前检查

每次修改后，至少检查：

1. 是否沿用了现有目录结构
2. 是否把业务逻辑留在了 ViewModel / UseCase / Entity
3. 是否补了需要的 migration / mock / dto
4. 是否覆盖了 `idle/loading/success/empty/error`
5. 是否说明了启动与验证方式
