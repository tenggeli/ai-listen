# Listen Codex 启动入口

本文件是 Codex 进入 Listen 项目时的首读文档，用来快速统一目标、边界、代码现状与后续执行顺序。

## 1. 当前项目状态

当前仓库已经不是“纯骨架”，而是处于 `P0 早期可运行阶段`。

已落地模块：

- 用户 Web：AI 主入口页面 `/home`
- 用户 Web：AI 对话页 `/chat`
- 管理后台 Web：服务方审核页 `/admin/providers/review`
- Go 后端：AI 会话、消息追加、每日匹配次数、AI 推荐结果（Mock Match Service）
- Go 后端：服务方审核列表、详情、通过 / 拒绝 / 补充资料
- MySQL 迁移：AI 会话、AI 消息、匹配结果表、服务方审核相关表、匹配次数表
- 仓储切换：支持 `memory/mysql` 双仓储驱动

仍为骨架或未开始的部分：

- 用户 Web：服务页、声音页、我的页
- 管理后台：登录鉴权、服务项目管理、声音内容管理、平台配置、订单管理
- 用户 App：仅工程骨架
- 真正 AI 网关、真实支付、结算、提现、VIP

## 2. 首次读取顺序

Codex 进入仓库后，建议固定按以下顺序建立上下文：

1. `doc/03-AI开发/00-Codex-启动入口.md`
2. `doc/03-AI开发/01-技术方案设计.md`
3. `doc/03-AI开发/02-AI-Coding-指导.md`
4. `README.md`
5. `doc/01-产品设计/01-产品设计总览.md`
6. 相关 PRD
7. 相关原型 HTML

## 3. 当前裁决规则

若文档之间出现冲突，统一按以下规则处理：

1. AI 主入口优先于传统功能堆叠
2. 用户 Web 优先于管理后台，管理后台优先于用户 App
3. 已有目录与接口现状优先，不按旧文档重建一套新结构
4. 支付、结算、提现、VIP 等低优先能力不得阻塞 P0
5. 服务方当前优先按“平台内业务角色”实现，不急于独立成第四端
6. 开发环境未必能直连 MySQL；如涉及真实库验证，可保留人工验证步骤

## 4. 当前实际目录

```text
web/apps/user-web      用户 Web（Vue + Vite）
web/apps/admin-web     管理后台 Web（Vue + Vite）
app/user-app           用户 App 骨架（Vue + Vite）
backend                Go API 服务
backend/migrations     MySQL 迁移脚本
```

后端当前分层：

- `backend/internal/interface/http/user`
- `backend/internal/interface/http/admin`
- `backend/internal/application/ai`
- `backend/internal/application/provider`
- `backend/internal/domain/ai`
- `backend/internal/domain/provider`
- `backend/internal/infrastructure/...`

## 5. 当前推荐实现顺序

在现有代码基础上，建议按下列顺序继续推进：

1. 用户 Web 首页与 AI 对话页联调真实后端
2. 用户 Web 服务浏览页
3. 用户 Web 声音页
4. 用户 Web 我的页
5. 管理后台登录鉴权骨架
6. 管理后台服务项目管理
7. 管理后台声音内容管理
8. 订单骨架、评价、投诉

## 6. 当前最重要的业务边界

P0 继续聚焦：

- 首页 AI 主入口
- AI 对话页
- AI 会话与消息接口
- 每日匹配次数限制
- AI 推荐结果承接
- 后台服务方审核模块

P0 可保留占位但不做深实现：

- 真支付
- 托管结算
- 提现
- VIP
- 广场复杂互动
- 真正流式 AI 对话
- 真实语音输入

## 7. 当前可直接运行的入口

后端：

```bash
cd backend
go run ./cmd/server
```

默认端口 `8080`，健康检查：`GET /healthz`

用户 Web：

```bash
cd web/apps/user-web
npm install
npm run dev
```

默认端口 `5173`

管理后台：

```bash
cd web/apps/admin-web
npm install
npm run dev
```

默认端口 `5174`

用户 App：

```bash
cd app/user-app
npm install
npm run dev
```

默认端口 `5175`

## 8. 给 Codex 的编码约束

- 后端采用 Go + MySQL，按 `interface / application / domain / infrastructure` 分层。
- 前端采用 Vue，按 `views / components / domain / application / api` 分层。
- 业务逻辑不要直接写进 controller 或页面模板。
- 所有状态流转必须收敛到实体或领域服务。
- 优先在现有模块上续写，不平移出第二套目录结构。
- 每次实现只做一个清晰任务颗粒，不要一次铺开整个项目。

## 9. 常见误区

- 把仓库误判成“还没开始”，重复初始化工程
- 忽略已经存在的 `/api/v1/ai/*` 与 `/api/v1/admin/providers/*` 接口
- 把首页做成普通信息流，而不是 AI 驱动入口
- 先实现支付，导致主链路迟迟不可运行
- 在页面里散落业务枚举、状态判断和接口转换逻辑

## 10. 推荐任务写法

给 Codex 下任务时，优先使用这种格式：

```text
请基于 doc/03-AI开发/01-技术方案设计.md 与 doc/03-AI开发/02-AI-Coding-指导.md，
在现有代码基础上实现 [某个单一模块] 的最小可运行版本。

要求：
1. 明确端：用户 Web / 管理后台 / 用户 App / 后端
2. 明确范围：只做一个页面、一个模块或一个状态机
3. 明确是否允许 mock
4. 明确是否需要 migration、测试、API 契约
5. 默认不实现支付、结算、提现
6. 若涉及已有模块，优先复用现有目录与接口
```

## 11. 原始资料位置

如果需要补充理解业务背景，可再查看：

- `doc/01-产品设计/00-产品构想原稿.md`
- `doc/01-产品设计/01-产品设计总览.md`
- `doc/02-PRD/`
- `doc/04-原型素材/`

这些文档主要用于补充业务背景，不应替代本目录中的 AI 开发文档作为最终执行依据。
