# app workspace

App 端已拆分为两个独立工程：
- `apps/user-app`：用户 App（uni-app + Vue3）
- `apps/admin-app`：管理 App（uni-app + Vue3）

共享目录：
- `packages/shared`：跨 App 共享类型、常量、工具函数
  - `src/http`：统一请求封装（token 注入、未授权回调、标准响应解包）
  - `src/api`：接口类型定义（用户端、管理端、通用类型）
  - `src/design`：设计令牌（颜色、半径、间距、主题）

## 目录约定

用户 App：
- `apps/user-app/src/pages`：页面路由页面
- `apps/user-app/src/router`：路由常量与导航方法
- `apps/user-app/src/sdk`：用户端接口 SDK（依赖 shared api 类型）

管理 App：
- `apps/admin-app/src/pages`：页面路由页面
- `apps/admin-app/src/router`：路由常量与导航方法
- `apps/admin-app/src/sdk`：管理端接口 SDK（依赖 shared api 类型）

## 快速开始

```bash
cd app
pnpm install
```

启动用户 App（H5）：

```bash
pnpm --filter listen-user-app dev:h5
```

启动管理 App（H5）：

```bash
pnpm --filter listen-admin-app dev:h5
```
