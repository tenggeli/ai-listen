# Listen Codex 启动入口

本文件是 Codex 进入 Listen 项目时的首读文档，用来快速统一目标、边界、代码现状、配置约束与后续执行顺序。

## 1. 当前项目状态

当前仓库处于 `P0 可持续推进阶段`，已经具备基础 AI 能力与后台审核能力，产品原型已经扩展到完整的新用户进入链路与交易闭环原型。

已落地代码模块：

- 用户 Web：AI 主入口页面 `/home`
- 用户 Web：AI 对话页 `/chat`
- 管理后台 Web：服务方审核页 `/admin/providers/review`
- Go 后端：AI 会话、消息追加、每日匹配次数、AI 推荐结果（当前为 Mock Match Service）
- Go 后端：服务方审核列表、详情、通过 / 拒绝 / 补充资料
- MySQL 迁移：AI 会话、AI 消息、匹配结果、服务方审核、匹配次数等表
- 仓储切换：支持 `memory/mysql` 双仓储驱动

已完成原型、待继续开发的业务链路：

- 登录 / 注册
- 基础资料补全
- 性格设置（兴趣、MBTI、可跳过）
- 首页 / AI 对话 / 服务页 / 服务方详情 / 声音页 / 广场页 / 我的页
- 支付确认 / 订单详情 / 评价投诉

仍未落地或仅占位的部分：

- 用户登录后端与资料后端
- 服务浏览后端与声音内容后端
- 订单后端状态机与评价投诉后端
- 管理后台登录鉴权、服务项目管理、声音内容管理、订单与投诉管理
- 用户 App 业务页面
- 第三方对接能力（AI 网关、短信、微信登录、支付）
- 结算、提现、VIP

## 2. 首次读取顺序

Codex 进入仓库后，建议固定按以下顺序建立上下文：

1. `doc/03-AI开发/00-Codex-启动入口.md`
2. `doc/03-AI开发/01-技术方案设计.md`
3. `doc/03-AI开发/02-AI-Coding-指导.md`
4. `doc/04-原型素材/listen_prototype_hub.html`
5. 相关原型 HTML
6. `README.md`
7. `doc/01-产品设计/01-产品设计总览.md`
8. 相关 PRD

## 3. 当前裁决规则

若文档之间出现冲突，统一按以下规则处理：

1. AI 主入口优先于传统信息流
2. 用户 Web 优先于管理后台，管理后台优先于用户 App
3. 已有目录与接口现状优先，不按旧文档重建第二套结构
4. 所有配置优先按“用户目录配置文件”方案处理
5. 数据尽量从数据库读取，不新增大段页面硬编码业务数据
6. 所有涉及第三方交互的逻辑统一返回 mock 结果，当前不做真实对接
7. 支付相关开发放在最后，订单能力可以先做，最终支付统一先 mock 成功

## 4. 配置与环境约束

### 4.1 配置文件来源

后端服务配置相关修改，统一改为优先读取用户目录下配置文件：

- `~/conf/listenbase.cof`

要求：

- 不要把核心运行配置硬编码在仓库内
- 仓库内如保留默认配置，只能作为开发兜底，不作为主入口
- 新增配置项时，优先扩展 `listenbase.cof` 对应结构

### 4.2 数据库约束

涉及数据相关开发时，尽量设计为从数据库读取。

如开发依赖数据库配置，可使用：

```text
MYSQL_DSN=root:YUIO980-klio@tcp(192.168.2.24:13306)/listen?parseTime=true&charset=utf8mb4
```

使用原则：

- 查询类、列表类、详情类接口优先走 MySQL 仓储
- mock 数据只作为前端临时演示或 AI 临时输出，不作为长期数据源
- 若本地环境无法连库，允许保留 memory/mock 验证步骤，但接口设计仍要面向数据库落地

### 4.3 AI 数据约束

AI 相关能力当前统一采用 `AI Mock Data + 明确 Adapter 边界` 策略：

- 首页 AI 推荐结果先 mock
- AI 对话自动回复先 mock
- AI 推荐的服务方 / 声音 / 情绪承接文案先 mock
- 当前不对接真实 AI 网关

### 4.4 第三方交互约束

所有涉及第三方交互的逻辑，统一先 mock，暂不进行真实对接。

包括但不限于：

- AI 网关
- 短信验证码发送
- 微信登录 / 微信授权
- 支付下单
- 支付回调
- 第三方内容平台或第三方媒体服务

要求：

- 前后端都要保留 adapter / gateway 边界
- controller / 页面只能拿到 mock 结果，不直接耦合第三方 SDK
- 若要演示登录、支付、AI 结果，统一走本地 mock 成功 / mock 返回数据

### 4.5 支付约束

支付相关开发统一放在最后。

当前阶段允许：

- 订单创建
- 订单状态流转
- 订单详情展示
- 评价 / 投诉
- 支付确认页交互

当前阶段不做真实接入：

- 微信支付真实下单
- 回调验签
- 清结算
- 提现

所有最终支付结果统一先按：`mock 成功` 处理。

### 4.6 登录演示默认值

登录页面默认展示：

- 手机号：`13800113800`
- 验证码：`123`

约束：

- 登录页默认填充上述值，便于开发联调与原型演示
- 短信发送动作也统一 mock，不真实发送验证码

## 5. 当前实际目录

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

后续建议扩展：

- `backend/internal/application/identity`
- `backend/internal/application/service`
- `backend/internal/application/audio`
- `backend/internal/application/order`
- `backend/internal/domain/identity`
- `backend/internal/domain/service`
- `backend/internal/domain/audio`
- `backend/internal/domain/order`

## 6. 当前推荐实现顺序

在现有代码基础上，建议按下列顺序继续推进：

1. 用户登录 / 资料补全 / 性格设置基础链路
2. 用户 Web 服务页
3. 用户 Web 声音页
4. 用户 Web 我的页
5. 服务浏览后端与声音内容后端
6. 订单后端骨架与订单前端详情页
7. 评价 / 投诉骨架
8. 管理后台登录鉴权骨架
9. 管理后台服务项目管理
10. 管理后台声音内容管理
11. 管理后台订单 / 投诉管理
12. 用户 App 页面化

## 7. 当前最重要的业务边界

P0 / P1 继续聚焦：

- 登录 / 注册占位链路
- 基础资料补全
- 性格设置
- 首页 AI 主入口
- AI 对话页
- 服务浏览与服务方详情
- 声音页
- 我的页
- 订单创建与订单详情
- 评价与投诉
- 后台服务方审核

当前可保留占位但不做深实现：

- 真支付
- 托管结算
- 提现
- VIP
- 广场复杂互动
- 真正流式 AI 对话
- 真实语音输入
- AI 多模型路由

## 8. 当前可直接运行的入口

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

## 9. 给 Codex 的编码约束

- 后端采用 Go + MySQL，按 `interface / application / domain / infrastructure` 分层。
- 前端采用 Vue，按 `views / components / domain / application / api` 分层。
- 业务逻辑不要直接写进 controller 或页面模板。
- 所有状态流转必须收敛到实体或领域服务。
- 优先在现有模块上续写，不平移出第二套目录结构。
- 页面 mock 数据与第三方接口边界必须通过 adapter 隔离。
- 数据列表、详情、字典、配置优先设计为数据库读取。
- 每次实现只做一个清晰任务颗粒，不要一次铺开整个项目。

## 10. 常见误区

- 把仓库误判成“还没开始”，重复初始化工程
- 忽略原型已经扩展到登录、订单、评价投诉闭环
- 把配置继续写死在仓库文件里，而不是读取 `~/conf/listenbase.cof`
- 在页面里散落假枚举、假状态，而不抽到 domain / application
- 先实现真实支付，导致主链路开发阻塞
- 直接把 AI 文案或第三方结果写死到页面，而不预留 mock adapter

## 11. 推荐任务写法

给 Codex 下任务时，优先使用这种格式：

```text
请基于 doc/03-AI开发/01-技术方案设计.md 与 doc/03-AI开发/02-AI-Coding-指导.md，
在现有代码基础上实现 [某个单一模块] 的最小可运行版本。

要求：
1. 明确端：用户 Web / 管理后台 / 用户 App / 后端
2. 明确范围：只做一个页面、一个模块或一个状态机
3. 明确是否允许 mock
4. 明确是否需要 migration、测试、API 契约
5. 默认不实现真实支付、结算、提现
6. 若涉及已有模块，优先复用现有目录与接口
7. 后端配置优先读取 ~/conf/listenbase.cof
8. 涉及数据的能力优先按 MySQL 读取设计
```

## 12. 原始资料位置

如果需要补充理解业务背景，可再查看：

- `doc/01-产品设计/00-产品构想原稿.md`
- `doc/01-产品设计/01-产品设计总览.md`
- `doc/02-PRD/`
- `doc/04-原型素材/`

这些文档主要用于补充业务背景，不应替代本目录中的 AI 开发文档作为最终执行依据。
