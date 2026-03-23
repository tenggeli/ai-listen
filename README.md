# ai-listen

根据 `doc/` 下产品与技术文档实现的 MVP 后端工程。

## 技术栈
- 后端：Go 1.22 + Gin + GORM
- 数据库：MySQL 8
- 缓存：Redis 7
- 网关：Nginx

## 已实现 MVP 模块
- 登录/鉴权：
  - `POST /api/auth/send-code`
  - `POST /api/auth/mobile-login`
  - `POST /api/auth/logout`
- 用户：
  - `GET /api/user/profile`
- 服务方入驻：
  - `POST /api/provider-center/apply`
  - `GET /api/provider-center/audit-status`
  - `PUT /api/provider-center/profile`

统一返回结构：
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

## 启动步骤
### 1) 安装依赖
```bash
go mod tidy
```

### 2) 启动 MySQL / Redis
```bash
cd deploy
docker compose up -d mysql redis
cd ..
```

### 3) 初始化配置
```bash
cp config.example.yaml config.yaml
```

### 4) 启动 API
```bash
go run ./cmd/api -config config.yaml
```

健康检查：
- `GET http://localhost:8080/health`

## SQL 初始化
- `scripts/mysql/init/001_init.sql`：创建数据库
- `scripts/mysql/init/002_mvp_auth_provider.sql`：`users`、`provider_applications`、`providers` 三张 MVP 核心表

通过 `deploy/docker-compose.yml` 启动 MySQL 时会自动执行以上脚本。

## 接口调用示例
### 1) 发送验证码
```bash
curl -X POST http://localhost:8080/api/auth/send-code \
  -H "Content-Type: application/json" \
  -d '{"mobile":"13800000000"}'
```

### 2) 手机号验证码登录
```bash
curl -X POST http://localhost:8080/api/auth/mobile-login \
  -H "Content-Type: application/json" \
  -d '{"mobile":"13800000000","code":"<send-code返回的debugCode>"}'
```

登录成功后，从返回中取 `data.token` 并放到请求头：
- `Authorization: Bearer <token>`

### 3) 获取用户资料
```bash
curl http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer <token>"
```

### 4) 提交服务方入驻申请
```bash
curl -X POST http://localhost:8080/api/provider-center/apply \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "realName":"张三",
    "idCardNo":"110101199001011234",
    "idCardFront":"https://example.com/front.jpg",
    "idCardBack":"https://example.com/back.jpg",
    "faceVerifyStatus":1,
    "agreementSigned":true,
    "cityId":1,
    "intro":"5年陪伴服务经验",
    "photos":["https://example.com/p1.jpg"],
    "serviceDesc":"陪聊、陪看电影"
  }'
```

### 5) 获取审核状态
```bash
curl http://localhost:8080/api/provider-center/audit-status \
  -H "Authorization: Bearer <token>"
```

### 6) 更新服务方资料
```bash
curl -X PUT http://localhost:8080/api/provider-center/profile \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "intro":"新的服务介绍",
    "cityId":1,
    "zodiac":"龙",
    "constellation":"白羊座",
    "tags":["健谈","电影","同城"],
    "serviceStatus":1,
    "onlineStatus":1
  }'
```

### 7) 退出登录
```bash
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer <token>"
```

## 文档歧义下的 MVP 假设
- `send-code` 在 `dev` 环境会回传 `debugCode` 便于联调；`prod` 不回传。
- `mobile-login` 验证码通过 Redis 校验，成功后清除验证码。
- token 为随机字符串，Redis 存储 7 天过期（Bearer 方案）。
- `provider-center/profile` 仅允许审核通过（`audit_status=1`）后更新；若审核通过但 `providers` 尚未创建，会自动创建一条基础 `providers` 记录。
