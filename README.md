# Listen Monorepo (MVP Skeleton)

本次初始化遵循 `doc/03-AI开发` 约束，优先完成：
- 基础工程骨架（用户 Web / 平台管理后台 / 用户 App / Go 后端）
- 用户 Web 首页 AI 主入口最小可运行闭环
- 平台管理后台服务方审核模块（列表/详情/审核动作）
- MySQL 仓储真实接入（支持 `memory/mysql` 切换）

## 目录

- `web/apps/user-web` 用户 Web（Vue）
- `web/apps/admin-web` 平台管理后台（已实现服务方审核模块）
- `doc/05-管理后台原型` 服务方管理后台 + 平台管理后台原型输入
- `app/user-app` 用户 App（Vue 骨架）
- `backend` Go 统一接口服务（`interface/application/domain/infrastructure`）
- `deploy/nginx` Nginx 路由配置草案

## 快速启动

### 1) 启动后端

```bash
cd backend
go run ./cmd/server
```

默认端口 `8080`，健康检查：`GET /healthz`。

后端配置读取优先级：
- 1) `~/conf/listenbase.cof`
- 2) 环境变量（如 `LISTEN_REPOSITORY_DRIVER`）
- 3) 代码默认值（仅开发兜底）

推荐先准备配置文件：

```bash
mkdir -p ~/conf
cp backend/config/listenbase.example.cof ~/conf/listenbase.cof
```

关键配置项：
- `server.port`
- `repository.driver`（`memory/mysql`）
- `mysql.dsn`
- `ai.mode`（`mock/real`）
- `mock.enable_payment_success`

接口分组规范：
- 用户侧：`/api/v1/...`
- 服务方侧：`/api/v1/provider/...`
- 平台管理侧：`/api/v1/admin/...`

新增接口命名约束：
- 路径使用小写英文
- 资源名优先复数，如 `orders`、`providers`
- 动作型接口放在资源详情后，如 `/orders/{id}/accept`
- 单例资源允许使用单数，如 `profile`

历史接口说明：
- 现有已实现接口不因该规范回溯重命名
- `GET /healthz` 属于基础设施健康检查接口

用户登录骨架接口（P0-02）：
- `POST /api/v1/auth/login/sms`（验证码 mock：`123`）
- `POST /api/v1/auth/login/wechat/mock`

用户基础资料接口（P0-03）：
- `GET /api/v1/users/me`
- `PUT /api/v1/users/me/profile`

用户性格设置接口（P0-04）：
- `PUT /api/v1/users/me/personality`
- `POST /api/v1/users/me/personality/skip`

### 2) 启动用户 Web

```bash
cd web/apps/user-web
npm install
npm run dev
```

默认端口 `5173`。

- 首页默认直连 Go 后端 `/api/v1/ai/home`、`/api/v1/ai/match/remaining`、`/api/v1/ai/match`
- 如需自定义用户 Web 首页接口前缀：设置 `VITE_AI_API_BASE_URL`
- `Chat` / `声音` 页面仍可通过 `VITE_USE_MOCK=false` 切到真实后端

### 3) 启动平台管理后台

```bash
cd web/apps/admin-web
npm install
npm run dev
```

默认端口 `5174`。

- 默认使用 Mock：`VITE_USE_MOCK` 非 `false`
- 对接真实后端：设置 `VITE_USE_MOCK=false`
- 审核页面路由：`/admin/providers/review`
- 当前仓库中的 `admin-web` 表示“平台管理后台”，不是服务方管理后台

### 4) 启动用户 App 骨架

```bash
cd app/user-app
npm install
npm run dev
```

默认端口 `5175`。
