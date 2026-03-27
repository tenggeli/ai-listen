# Listen Monorepo (MVP Skeleton)

本次初始化遵循 `doc/03-AI开发` 约束，优先完成：
- 基础工程骨架（用户 Web / 管理后台 / 用户 App / Go 后端）
- 用户 Web 首页 AI 主入口最小可运行闭环
- 管理后台服务方审核模块（列表/详情/审核动作）
- MySQL 仓储真实接入（支持 `memory/mysql` 切换）

## 目录

- `web/apps/user-web` 用户 Web（Vue）
- `web/apps/admin-web` 管理后台（已实现服务方审核模块）
- `app/user-app` 用户 App（Vue 骨架）
- `backend` Go 后端（`interface/application/domain/infrastructure`）
- `deploy/nginx` Nginx 路由配置草案

## 快速启动

### 1) 启动后端

```bash
cd backend
go run ./cmd/server
```

默认端口 `8080`，健康检查：`GET /healthz`。

后端仓储驱动：
- 默认：`LISTEN_REPOSITORY_DRIVER=memory`
- MySQL：`LISTEN_REPOSITORY_DRIVER=mysql`
- DSN：`LISTEN_MYSQL_DSN`（示例默认值见 `internal/infrastructure/config/config.go`）

### 2) 启动用户 Web

```bash
cd web/apps/user-web
npm install
npm run dev
```

默认端口 `5173`。

- 默认使用 Mock：`VITE_USE_MOCK` 非 `false`
- 对接真实后端：设置 `VITE_USE_MOCK=false`

### 3) 启动管理后台

```bash
cd web/apps/admin-web
npm install
npm run dev
```

默认端口 `5174`。

- 默认使用 Mock：`VITE_USE_MOCK` 非 `false`
- 对接真实后端：设置 `VITE_USE_MOCK=false`
- 审核页面路由：`/admin/providers/review`

### 4) 启动用户 App 骨架

```bash
cd app/user-app
npm install
npm run dev
```

默认端口 `5175`。
