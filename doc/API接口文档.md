# ai-listen API 接口文档

## 1. 文档说明
本文件为 ai-listen 项目的接口目录级文档，适用于前后端联调、接口评审、测试用例设计。

建议规范：
- 协议：HTTPS
- 风格：RESTful + JSON
- 字符编码：UTF-8
- 鉴权：Token / JWT

统一返回示例：
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

统一错误码建议：
- 0：成功
- 400：参数错误
- 401：未登录
- 403：无权限
- 404：资源不存在
- 500：服务器错误

---

## 2. 鉴权接口

### 2.1 微信登录
- `POST /api/auth/wechat-login`

### 2.2 手机验证码登录
- `POST /api/auth/mobile-login`

### 2.3 发送验证码
- `POST /api/auth/send-code`

### 2.4 退出登录
- `POST /api/auth/logout`

---

## 3. 用户接口

### 3.1 获取用户资料
- `GET /api/user/profile`

### 3.2 更新用户资料
- `PUT /api/user/profile`

### 3.3 获取用户首页信息
- `GET /api/user/home-summary`

### 3.4 获取我的订单
- `GET /api/user/orders`

### 3.5 获取收藏列表
- `GET /api/user/favorites`

---

## 4. 首页 AI 接口

### 4.1 AI 匹配
- `POST /api/home/ai-match`

请求参数建议：
```json
{
  "inputType": "voice",
  "content": "帮我找个今晚7点一起吃饭的电影搭子",
  "cityId": 1,
  "lng": 120.1,
  "lat": 30.2
}
```

### 4.2 AI 追问
- `POST /api/home/ai-followup`

### 4.3 手动筛选推荐
- `POST /api/home/filter-match`

### 4.4 获取首页示例标签
- `GET /api/home/tags`

---

## 5. 服务提供方接口

### 5.1 服务方列表
- `GET /api/provider/list`

### 5.2 服务方详情
- `GET /api/provider/detail/{id}`

### 5.3 服务方动态列表
- `GET /api/provider/{id}/posts`

### 5.4 服务方评价列表
- `GET /api/provider/{id}/reviews`

### 5.5 收藏服务方
- `POST /api/provider/favorite`

---

## 6. 订单接口

### 6.1 创建订单
- `POST /api/order/create`

### 6.2 订单支付
- `POST /api/order/pay`

### 6.3 获取订单详情
- `GET /api/order/detail/{id}`

### 6.4 取消订单
- `POST /api/order/cancel`

### 6.5 用户确认开始
- `POST /api/order/confirm-start`

### 6.6 用户确认完成
- `POST /api/order/confirm-finish`

### 6.7 用户发起投诉
- `POST /api/order/complaint`

### 6.8 订单列表
- `GET /api/order/list`

---

## 7. 评价接口

### 7.1 提交评价
- `POST /api/review/create`

### 7.2 追评
- `POST /api/review/append`

### 7.3 评价列表
- `GET /api/review/list`

---

## 8. 广场接口

### 8.1 动态列表
- `GET /api/post/list`

### 8.2 发布动态
- `POST /api/post/create`

### 8.3 动态详情
- `GET /api/post/detail/{id}`

### 8.4 评论动态
- `POST /api/post/comment`

### 8.5 点赞动态
- `POST /api/post/like`

### 8.6 收藏动态
- `POST /api/post/favorite`

### 8.7 举报动态
- `POST /api/post/report`

---

## 9. 声音接口

### 9.1 音频列表
- `GET /api/audio/list`

### 9.2 音频详情
- `GET /api/audio/detail/{id}`

### 9.3 音频点赞
- `POST /api/audio/like`

### 9.4 音频收藏
- `POST /api/audio/favorite`

### 9.5 音频评论
- `POST /api/audio/comment`

---

## 10. 钱包与 VIP 接口

### 10.1 钱包信息
- `GET /api/wallet/info`

### 10.2 钱包流水
- `GET /api/wallet/transactions`

### 10.3 钱包充值
- `POST /api/wallet/recharge`

### 10.4 开通 VIP
- `POST /api/vip/open`

### 10.5 VIP 权益信息
- `GET /api/vip/info`

---

## 11. 消息接口

### 11.1 消息列表
- `GET /api/message/list`

### 11.2 标记已读
- `POST /api/message/read`

### 11.3 未读数
- `GET /api/message/unread-count`

---

## 12. 服务提供方端接口

### 12.1 提交入驻申请
- `POST /api/provider-center/apply`

### 12.2 获取审核状态
- `GET /api/provider-center/audit-status`

### 12.3 更新服务方资料
- `PUT /api/provider-center/profile`

### 12.4 新增服务项目
- `POST /api/provider-center/service-item/create`

### 12.5 编辑服务项目
- `PUT /api/provider-center/service-item/update`

### 12.6 服务项目列表
- `GET /api/provider-center/service-item/list`

### 12.7 切换工作状态
- `POST /api/provider-center/status/change`

### 12.8 服务方订单列表
- `GET /api/provider-center/order/list`

### 12.9 接单
- `POST /api/provider-center/order/accept`

### 12.10 出发
- `POST /api/provider-center/order/depart`

### 12.11 到达
- `POST /api/provider-center/order/arrive`

### 12.12 开始服务
- `POST /api/provider-center/order/start`

### 12.13 完成订单
- `POST /api/provider-center/order/finish`

### 12.14 异常取消
- `POST /api/provider-center/order/cancel`

### 12.15 收益中心
- `GET /api/provider-center/income/summary`

### 12.16 提现申请
- `POST /api/provider-center/withdraw/apply`

### 12.17 提现记录
- `GET /api/provider-center/withdraw/list`

---

## 13. 后台管理接口

### 13.1 服务方审核列表
- `GET /admin/provider/audit/list`

### 13.2 审核服务方
- `POST /admin/provider/audit/do`

### 13.3 订单列表
- `GET /admin/order/list`

### 13.4 订单详情
- `GET /admin/order/detail/{id}`

### 13.5 投诉处理
- `POST /admin/order/complaint/handle`

### 13.6 用户列表
- `GET /admin/user/list`

### 13.7 服务方列表
- `GET /admin/provider/list`

### 13.8 服务项目管理
- `GET /admin/service-item/list`
- `POST /admin/service-item/save`

### 13.9 财务报表
- `GET /admin/report/revenue`

### 13.10 提现审核
- `POST /admin/withdraw/audit`

### 13.11 内容审核
- `GET /admin/content/post/list`
- `POST /admin/content/post/audit`
- `GET /admin/content/audio/list`
- `POST /admin/content/audio/audit`

---

## 14. WebSocket 事件建议

### 用户侧
- `order.accepted`
- `order.departed`
- `order.arrived`
- `order.started`
- `order.finished`
- `message.new`

### 服务方侧
- `order.new`
- `order.canceled`
- `order.complaint`
- `message.new`

---

## 15. 接口开发优先级

### P0（MVP 必做）
- 鉴权
- 首页 AI 匹配
- 服务方列表/详情
- 订单创建/支付/状态流转
- 服务方接单/出发/到达/开始/完成
- 评价
- 审核后台

### P1
- 广场
- 声音
- 钱包/VIP
- 消息中心

### P2
- 高级运营能力
- 更复杂的推荐接口
- 红娘/管家服务

---

## 16. 说明
当前为目录级接口文档，下一步可继续拆为：
- 请求参数明细
- 响应参数明细
- 错误码细化
- 示例报文
- Swagger/OpenAPI 文档
