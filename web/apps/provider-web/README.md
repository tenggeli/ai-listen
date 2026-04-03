# provider-web

服务方后台前端工程（Vue 3 + Vite），当前阶段聚焦登录鉴权骨架与工作台占位。

## 已实现能力

- 独立工程：`web/apps/provider-web`
- 路由：
  - `/provider/login`
  - `/provider/dashboard`
  - `/provider/orders`
  - `/provider/orders/:id`
  - `/provider/profile`
  - `/provider/services`
- 登录接口：`POST /api/v1/provider/auth/login/mock`
- 会话校验接口：`GET /api/v1/provider/profile`
- 资料编辑接口：`PUT /api/v1/provider/profile`
- 服务项目接口：`GET /api/v1/provider/services`
- 会话存储：`localStorage`（key: `provider_web_auth_session`）
- 路由守卫：
  - 未登录访问受保护页面会跳转登录页
  - 已登录访问登录页会跳转工作台
- 刷新恢复：
  - 刷新后通过 `GET /api/v1/provider/profile` 校验 token，成功则保留会话
  - token 失效（401）则清空会话并回到登录页
- 页面状态：统一覆盖 `idle/loading/success/error`
- 资料管理能力：
  - 读取当前资料
  - 编辑昵称与城市编码并保存
- 服务项目管理最小版：
  - 服务项目列表读取
  - 空态与异常态反馈
- 订单查看能力：
  - 列表分页
  - 状态筛选
  - 详情查看（含状态标签、创建时间、支付时间）
  - 空列表与异常反馈

## 本地运行

```bash
cd web/apps/provider-web
npm install
npm run dev
```

默认端口：`5176`，并将 `/api` 代理到 `http://localhost:8080`。

## 构建

```bash
npm run build
```

## 测试

```bash
npm run test:run
```

覆盖场景：

- 登录成功
- 登录失败
- token 401 失效处理
- 路由守卫重定向
- 刷新恢复会话
- 订单分页边界
- 订单空态
- 弱网超时错误
- 订单详情 404 错误
