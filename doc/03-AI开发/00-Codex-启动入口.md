# Listen Codex 启动入口

本文件是 Codex 进入 Listen 项目时的首读文档，用于快速对齐“当前代码真实状态、边界约束、执行顺序”。

## 1. 当前项目状态（以仓库代码为准）

当前仓库处于 `P1 扩展阶段`，已经从“仅 AI 首页 + 审核”扩展到“用户引导、AI 主入口、服务浏览、声音页、我的页、设置页后端持久化、支付与订单最小闭环、平台管理后台鉴权与运营管理可运行、服务方侧后端接口已落地”。

当前端形态统一按以下口径理解：

- 用户 Web：面向 C 端用户
- 服务方管理后台：面向服务提供方的经营与履约后台
- 平台管理后台：面向平台运营/审核/客服/财务的后台
- 用户 App：后续移动端承接
- Go 后端：统一接口服务，后端不拆成多个独立服务

已落地模块：

- 用户 Web：`/auth` 登录页（手机号 + 微信 mock 登录）
- 用户 Web：`/profile/setup` 基础资料页
- 用户 Web：`/personality/setup` 性格设置页（兴趣 + MBTI + 跳过）
- 用户 Web：`/home` AI 主入口页
- 用户 Web：`/chat` AI 对话页
- 用户 Web：`/sound` 声音页，支持 `MockSoundApi` 与 HTTP `GET /api/v1/sounds`
- 用户 Web：`/services` 服务页
- 用户 Web：`/providers/:id` 服务方详情页
- 用户 Web：`/me` 我的页，可读取 `/api/v1/users/me`
- 用户 Web：`/payment/confirm` 支付确认页，当前调用后端创建订单并执行 mock-success 支付
- 用户 Web：`/orders/:id` 订单详情页，读取真实后端订单详情
- 用户 Web：`/orders` 订单列表页，读取真实后端订单列表
- 用户 Web：`/orders/:id/feedback` 评价/投诉页，读取与提交真实后端反馈
- 用户 Web：`/settings` 设置页，支持账号资料查看、设置项后端持久化（MySQL）与 memory 模式本地兜底、退出登录
- 平台管理后台 Web：当前仓库为 `web/apps/admin-web`，已实现 `/admin/login`、`/admin/dashboard`、`/admin/providers/review`、`/admin/services/manage`、`/admin/orders/manage`、`/admin/complaints/manage`
- 服务方管理后台：当前已有 `doc/05-管理后台原型`，前端工程尚未正式落地；但后端 provider 侧鉴权与订单履约接口已可用
- Go 后端：统一单体接口服务，当前已包含 AI、identity、user_settings、service discovery、sound、order、feedback、admin_auth、admin_order、service_item_admin、provider_auth 与平台管理侧/服务方侧路由
- 配置：后端优先读取 `~/conf/listenbase.cof`
- MySQL：AI、服务浏览、服务方审核、identity、order、feedback、user_settings、服务项目管理、后台订单操作日志均具备 migration / mysql repository，可在 `memory/mysql` 间切换

仍未落地或仅占位：

- Go 后端：真实支付模块、服务方资料与服务项目管理接口
- 服务方管理后台：前端工程骨架、登录鉴权、履约/经营/结算页面
- 平台管理后台：声音内容管理
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
8. `doc/05-管理后台原型/*.html`
9. 相关原型：`doc/04-原型素材/*.html`

## 3. 当前裁决规则

若文档与代码冲突，按以下优先级裁决：

1. 已有代码与路由现状优先
2. AI 主入口优先于传统信息流
3. 用户 Web 优先于双后台，平台管理后台优先于服务方管理后台，双后台优先于用户 App
4. 配置统一按 `~/conf/listenbase.cof` 规则
5. 数据结构优先面向 MySQL 设计，允许阶段性 mock/memory 验证
6. 所有第三方交互先 mock，暂不接真实 SDK
7. 支付能力当前仍允许先 mock-success，但真实支付仍暂不落地

## 4. 配置与环境约束

### 4.0 接口命名规范

后端接口按端分组如下：

- 用户侧：`/api/v1/...`
- 服务方侧：`/api/v1/provider/...`
- 平台管理侧：`/api/v1/admin/...`

新增接口命名要求：

- 路径统一使用小写英文与短横线风格，避免驼峰、拼音、下划线混用
- 资源名优先使用复数名词，例如 `orders`、`providers`、`service-items`
- 列表查询使用 `GET /resources`
- 详情查询使用 `GET /resources/{id}`
- 创建使用 `POST /resources`
- 更新优先使用 `PUT /resources/{id}`
- 动作型接口使用 `POST /resources/{id}/action`

历史接口说明：

- 现有已落地接口不因本规范回溯重命名
- 当前历史接口例外包括：`/api/v1/auth/login/*`、`/api/v1/users/me/*`、`/api/v1/ai/match/remaining`
- `GET /healthz` 属于基础设施健康检查接口，不属于用户侧/服务方侧/平台管理侧业务分组

推荐示例：

- 用户侧：`GET /api/v1/orders`
- 服务方侧：`POST /api/v1/provider/orders/{id}/accept`
- 平台管理侧：`POST /api/v1/admin/providers/{id}/approve`

### 4.1 配置文件来源

后端服务配置优先读取：

- `~/conf/listenbase.cof`

仓库内示例：

- `backend/config/listenbase.example.cof`

### 4.2 数据库约束

- 已有 MySQL 表：AI 会话/消息/匹配、每日配额、服务方审核、服务浏览、identity 用户账户、订单、订单反馈、用户设置、后台订单操作日志
- 新增业务（payment/admin 配置等）应继续按 MySQL 可落地结构设计
- 当前 `repository.driver` 可切 `memory/mysql`

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
web/apps/admin-web     平台管理后台 Web（Vue + Vite，当前已落地）
doc/05-管理后台原型      服务方管理后台 / 平台管理后台 原型输入
app/user-app           用户 App 骨架（Vue + Vite）
backend                Go 统一接口服务
backend/migrations     MySQL 迁移脚本
```

后端已落地模块：

- `application/ai`
- `application/admin_auth`
- `application/admin_order`
- `application/audio`
- `application/feedback`
- `application/identity`
- `application/order`
- `application/provider`
- `application/provider_auth`
- `application/service_discovery`
- `application/service_item_admin`
- `application/user_settings`
- `domain/admin_auth`
- `domain/ai`
- `domain/audio`
- `domain/feedback`
- `domain/identity`
- `domain/order`
- `domain/provider`
- `domain/provider_auth`
- `domain/service_discovery`
- `domain/service_item_admin`
- `domain/user_settings`

## 6. 当前推荐推进顺序

1. 服务方管理后台前端工程骨架
2. 双后台各自运营模块扩展
3. AI 对话自动回复后端化
4. 声音内容数据化
5. 订单状态流转扩展
6. 用户 App 页面化

## 7. 当前最重要的业务边界

P0/P1 聚焦：

- 用户引导链路稳态
- 首页 AI 主入口
- AI 对话页
- 服务浏览与服务方详情
- 声音页与后端接口
- 我的页与用户资料读取
- 支付确认 Mock 链路
- 订单与评价投诉最小闭环稳态
- 平台管理后台服务方审核稳态 + 继续补运营能力
- 服务方管理后台后续承接入驻、接单、履约、收益结算

当前可保留占位但不做深实现：

- 真实支付
- 清结算/提现
- 第三方真实网关接入
