# listen

listen 是一个以 AI 指令式交互为入口的智能陪伴服务平台，目标形态包含 App、后台管理端与自适应 Web。

## 目录结构

```text
ai-listen/
  backend/             Go 后端
  web/                 Vue 前端
  app/                 App 端工程占位
  deploy/              Nginx 与 Docker 配置
  doc/                 产品、技术与协作文档
```

## 当前阶段

当前已完成：
- 项目方案与技术方案文档
- M1 基础工程目录初始化
- Go 后端基础服务骨架
- Vue 用户端和后台端基础骨架
- Nginx 和 Docker Compose 示例配置

## 启动说明

### 后端

```bash
cd backend
go run ./cmd/server
```

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

## 下一步建议

1. 完成数据库设计文档
2. 初始化 MySQL 表结构与迁移脚本
3. 实现认证与用户模块
4. 接入首页和声音页静态页面

