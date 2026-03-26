# listen

listen 是一个以 AI 交互为入口的陪伴服务平台，包含用户端与管理端的 Web/App，以及统一 Go 后端。

## 项目结构

```text
ai-listen/
  backend/             Go API 服务
  web/                 Web 前端（user-web、admin-web）
  app/                 App 前端（user-app、admin-app）
  deploy/              Nginx 与 Docker Compose 示例
  doc/                 项目文档
```

## 当前状态（简要）
- 后端已完成统一路由、认证、订单与支付基础链路、管理端 RBAC 基础能力。
- Web 与 App 四端均已建立工程骨架；App 已拆分为 `user-app` 与 `admin-app`。
- 用户 Web 已补齐首页、短信登录、我的订单、个人中心基础业务页并完成接口联调。
- 音频与 AI 匹配模块已从占位升级为基础可用实现（匹配会话、音频列表/详情/播放记录/收藏）。
- 存储支持内存/MySQL 双实现；未配置 `MYSQL_DSN` 时默认走内存模式。
- 部分模块仍为轻实现（内容互动深水区、投诉治理深度能力、资金闭环等）。

## 快速启动

### 1) 启动后端（内存模式）

```bash
cd backend
go run ./cmd/server
```

### 2) 启动后端（MySQL 模式）

```bash
cd backend
MYSQL_DSN='user:password@tcp(127.0.0.1:3306)/listen?charset=utf8mb4&parseTime=True&loc=Local' go run ./cmd/server
```

支付回调签名密钥可通过 `PAYMENT_CALLBACK_SECRET` 配置（默认开发值：`listen-dev-callback-secret`）。

### 3) 启动 Web

用户端：
```bash
cd web/apps/user-web
npm install
npm run dev
```

管理端：
```bash
cd web/apps/admin-web
npm install
npm run dev
```

### 4) 启动 App（uni-app H5）

```bash
cd app
pnpm install
pnpm --filter listen-user-app dev:h5
pnpm --filter listen-admin-app dev:h5
```

## 常用信息
- 健康检查：`GET /api/v1/health`
- 管理员登录：`POST /api/v1/admin/auth/login`

默认开发账号（内存/MySQL 种子一致）：
- `admin / admin123456`
- `content_admin / admin123456`

## 测试

```bash
cd backend
go test ./...
```

## 文档索引（已精简对齐）
- [01-项目方案设计：定位、现状与优先级](./doc/01-项目方案设计.md)
- [02-技术方案设计：架构、目录、技术缺口](./doc/02-技术方案设计.md)
- [03-开发规范与面向对象设计：分层与协作约束](./doc/03-开发规范与面向对象设计.md)
- [04-开发任务拆解与里程碑：进度与后续计划](./doc/04-开发任务拆解与里程碑.md)
- [05-数据库设计：现有表与关键状态](./doc/05-数据库设计.md)
- [06-接口设计规范：完整路由清单与成熟度](./doc/06-接口设计规范.md)
- [07-四端页面矩阵与功能边界：四端职责与边界规则](./doc/07-四端页面矩阵与功能边界.md)
