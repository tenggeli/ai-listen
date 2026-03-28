# Listen Codex 启动入口

本文件是 Codex 进入 Listen 项目时的首读文档，用于快速对齐“当前代码真实状态、边界约束、执行顺序”。

## 1. 当前项目状态（以仓库代码为准）

当前仓库处于 `P0 持续推进阶段`，已经从“仅 AI 首页 + 审核”扩展到“用户登录与引导链路可运行”。

已落地模块：

- 用户 Web：`/auth` 登录页（手机号 + 微信 mock 登录）
- 用户 Web：`/profile/setup` 基础资料页
- 用户 Web：`/personality/setup` 性格设置页（兴趣 + MBTI + 跳过）
- 用户 Web：`/home` AI 主入口页
- 用户 Web：`/chat` AI 对话页
- 用户 Web：`/sound` 声音页（当前默认 MockSoundApi，可切 HTTP）
- 管理后台 Web：`/admin/dashboard`、`/admin/providers/review`
- Go 后端：AI 会话/匹配接口 + 身份登录/资料/性格接口 + 管理端服务方审核接口
- 配置：后端优先读取 `~/conf/listenbase.cof`
- MySQL：AI 与服务方审核相关迁移及仓储；支持 `memory/mysql` 切换（身份模块当前仍使用 memory 仓储）

仍未落地或仅占位：

- 用户 Web：服务页、服务方详情页、我的页、支付确认页、订单详情页、评价投诉页
- Go 后端：服务浏览、声音内容、订单、评价投诉、支付 mock 接口
- 管理后台：登录鉴权、服务项目管理、声音内容管理、订单/投诉管理
- 用户 App：仍为骨架，仅 `/home`
- 第三方真实对接：AI 网关、短信、微信真实授权、真实支付

## 2. 首次读取顺序

1. `doc/03-AI开发/00-Codex-启动入口.md`
2. `doc/03-AI开发/01-技术方案设计.md`
3. `doc/03-AI开发/02-AI-Coding-指导.md`
4. `README.md`
5. `web/apps/user-web/src/router/index.ts`
6. `backend/internal/interface/http/user/router.go`
7. `backend/internal/interface/http/admin/router.go`
8. 相关原型：`doc/04-原型素材/*.html`

## 3. 当前裁决规则

若文档与代码冲突，按以下优先级裁决：

1. 已有代码与路由现状优先
2. AI 主入口优先于传统信息流
3. 用户 Web 优先于管理后台，管理后台优先于用户 App
4. 配置统一按 `~/conf/listenbase.cof` 规则
5. 数据结构优先面向 MySQL 设计，允许阶段性 mock/memory 验证
6. 所有第三方交互先 mock，暂不接真实 SDK
7. 支付能力放在后续阶段，当前可先 mock-success

## 4. 配置与环境约束

### 4.1 配置文件来源

后端服务配置优先读取：

- `~/conf/listenbase.cof`

仓库内示例：

- `backend/config/listenbase.example.cof`

### 4.2 数据库约束

- 已有 MySQL 表：AI 会话/消息/匹配、每日配额、服务方审核相关
- 新增业务（identity/service/audio/order/feedback/payment）应优先按 MySQL 可落地结构设计
- 当前身份模块仓储为 memory；若切 mysql，需补 identity 的 migration 与 mysql repository

### 4.3 第三方交互约束

当前全部 mock：

- AI 推荐与自动回复
- 短信验证码发送
- 微信登录授权
- 支付下单与回调

### 4.4 登录演示默认值

登录页默认值：

- 手机号：`13800113800`
- 验证码：`123`

## 5. 当前实际目录

```text
web/apps/user-web      用户 Web（Vue + Vite）
web/apps/admin-web     管理后台 Web（Vue + Vite）
app/user-app           用户 App 骨架（Vue + Vite）
backend                Go API 服务
backend/migrations     MySQL 迁移脚本
```

后端已落地模块：

- `application/ai`
- `application/identity`
- `application/provider`
- `domain/ai`
- `domain/identity`
- `domain/provider`

## 6. 当前推荐推进顺序

1. 服务页 + 服务方详情页（用户 Web）
2. 服务浏览后端模块（service_discovery）
3. 声音内容后端模块（audio_content）并对齐 `/sound` HTTP 接口
4. 我的页（用户 Web）
5. 订单模块（后端 + 支付确认页 + 订单详情页）
6. 评价/投诉模块（后端 + 页面）
7. 管理后台登录鉴权骨架
8. 管理后台服务项目/声音内容/订单投诉管理
9. 用户 App 页面化

## 7. 当前最重要的业务边界

P0/P1 聚焦：

- 用户引导链路（已可运行，继续稳态）
- 首页 AI 主入口
- AI 对话页
- 服务浏览与服务方详情（下一优先）
- 声音页（前端已落地，后端待补）
- 订单、评价投诉（待建设）
- 后台服务方审核（已落地，继续补运营能力）

当前可保留占位但不做深实现：

- 真实支付
- 清结算/提现
- 第三方真实网关接入
