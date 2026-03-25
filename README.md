# listen

listen 是一个以 AI 指令式交互为入口的智能陪伴服务平台，目标形态包含 App、后台管理端与自适应 Web。

## 项目结构

```text
ai-listen/
  backend/             Go 后端
  web/                 Vue 前端
  app/                 App 端工程占位
  deploy/              Nginx 与 Docker 配置
  doc/                 产品、技术、数据库、接口文档
```

## 当前代码状态

当前仓库已经不是纯工程初始化状态，已完成以下核心能力：
- Go 后端基础服务与统一路由
- Vue 用户端与管理后台基础工程
- 短信登录与 token 刷新
- 服务提供方入驻申请、资料配置、服务项目配置
- 管理端服务提供方审核通过/拒绝
- 订单主链路：创建订单、创建支付单、接单、出发、到达、开始服务、完单、确认完成
- 后台管理员登录与后台接口鉴权
- 后台按路由粒度 RBAC 鉴权（基于角色权限点）
- 后台 RBAC 权限配置持久化（角色、权限点、角色分配）
- MySQL 存储实现，覆盖认证、服务提供方入驻审核、订单状态机主链路
- 内存存储实现，方便本地调试与测试

当前仍以占位为主的模块：
- AI 推荐
- 广场与声音内容深层逻辑
- 评价、投诉、退款、结算、提现
- 后台权限配置化（当前权限点在代码内置）

## 后端运行

### 方式一：默认内存存储

```bash
cd backend
go run ./cmd/server
```

说明：
- 未配置 `MYSQL_DSN` 时，后端默认使用内存存储
- 适合本地快速联调与测试

### 方式二：启用 MySQL

```bash
cd backend
MYSQL_DSN='user:password@tcp(127.0.0.1:3306)/listen?charset=utf8mb4&parseTime=True&loc=Local' go run ./cmd/server
```

说明：
- 启动时会自动执行初始化建表 SQL
- 启动时会自动初始化基础服务项目数据
- 需要本地准备 MySQL 8.0+

### 可用基础接口
- `/`
- `/api/v1/health`

### 后台默认开发账号

管理员登录接口：`POST /api/v1/admin/auth/login`
角色分配接口：`PUT /api/v1/admin/rbac/users/{adminUserId}/roles`

默认开发账号（内存模式与 MySQL 初始化种子一致）：
- `username`: `admin`
- `password`: `admin123456`
- `username`: `content_admin`
- `password`: `admin123456`

## 前端运行

### 用户端 Web

```bash
cd web/apps/user-web
npm install
npm run dev
```

### 管理后台

```bash
cd web/apps/admin-web
npm install
npm run dev
```

## Docker 示例

仓库已提供基础 `docker-compose` 示例：

```bash
docker compose -f deploy/docker/docker-compose.yml up
```

说明：
- 当前 Compose 主要用于快速拉起 Go 后端和 Nginx
- 默认示例未内置 MySQL 服务

## 测试

后端测试：

```bash
cd backend
go test ./...
```

当前已覆盖的核心测试链路：
- 短信登录
- 服务提供方申请与后台审核
- 下单与支付
- 履约完成

## 文档索引

- [项目方案设计](./doc/01-项目方案设计.md)
- [技术方案设计](./doc/02-技术方案设计.md)
- [开发规范与面向对象设计](./doc/03-开发规范与面向对象设计.md)
- [开发任务拆解与里程碑](./doc/04-开发任务拆解与里程碑.md)
- [数据库设计](./doc/05-数据库设计.md)
- [接口设计规范](./doc/06-接口设计规范.md)

## 下一步建议

1. 补齐真实支付回调、退款、结算、提现审核。
2. 为后台补管理员登录、鉴权与核心管理列表查询。
3. 为订单补超时取消、超时自动确认和异常单处理。
4. 推进评价、投诉、广场、声音与 AI 推荐真实实现。
