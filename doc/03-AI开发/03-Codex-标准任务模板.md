# Listen Codex 标准任务模板（最新版）

本文件用于给 Codex 下达稳定、低歧义、可直接执行的任务。

通用原则：

- 一次只下达一个清晰任务
- 先复用已有代码，再新增模块
- 明确端、范围、约束、交付与验证
- 第三方交互默认 mock
- 数据结构默认面向 MySQL 设计

## 1. 通用模板

```text
请基于以下文档完成开发任务：

- doc/03-AI开发/00-Codex-启动入口.md
- doc/03-AI开发/01-技术方案设计.md
- doc/03-AI开发/02-AI-Coding-指导.md
- 对应 PRD
- 对应原型 HTML

任务目标：
在现有代码基础上实现【模块名称】的最小可运行版本。

任务范围：
- 端：用户 Web / 管理后台 / 用户 App / Go 后端
- 页面或模块：【填写具体对象】
- 只做：【填写本次范围】
- 不做：【填写明确排除项】

实现要求：
1. 先核对现有路由、接口与目录再编码
2. 遵循现有分层与命名规范
3. 优先复用现有 ViewModel / UseCase / API Adapter / Entity
4. 业务逻辑不要直接写进页面模板或 Controller
5. 第三方交互统一 mock，并保留 adapter / gateway 扩展点
6. 页面状态覆盖 idle / loading / success / empty / error
7. 后端配置优先读取 ~/conf/listenbase.cof
8. 涉及数据按 MySQL 可落地结构设计
9. 若修改接口，说明前后端契约变更点
10. 登录相关默认值保持 13800113800 / 123

交付要求：
1. 直接修改代码，不只给方案
2. 列出修改文件
3. 给出验证步骤与结果
4. 列出未完成项与后续建议
```

## 2. 页面开发模板

```text
请基于现有代码实现【页面名称】页面最小可运行版本。

上下文：
- doc/03-AI开发/00-Codex-启动入口.md
- doc/03-AI开发/01-技术方案设计.md
- doc/03-AI开发/02-AI-Coding-指导.md
- 对应 PRD 与原型 HTML

要求：
1. 页面按 views / components / application / domain / api 分层
2. 优先复用已有 PageLoadState、ViewModel 与 API adapter
3. 明确页面入口与路由
4. 覆盖 idle/loading/success/empty/error
5. 接口未完成时可先 mock，但契约需与后端规划一致
6. 不把复杂业务判断写在模板中
7. 输出时说明关键交互与验证方式

本次不做：
- 【如真实支付、复杂风控、真实流式 AI 等】
```

## 3. 后端模块模板

```text
请在现有后端代码基础上实现【模块名称】最小闭环。

上下文：
- doc/03-AI开发/00-Codex-启动入口.md
- doc/03-AI开发/01-技术方案设计.md
- doc/03-AI开发/02-AI-Coding-指导.md

要求：
1. 严格按 interface / application / domain / infrastructure 分层
2. 先定义 entity、repository interface、use case、dto、controller
3. 状态流转集中在实体方法或领域服务
4. 提供最小接口并说明请求/响应结构
5. 涉及持久化变更必须补 migration
6. 至少补核心用例单测
7. 配置统一从 ~/conf/listenbase.cof 读取
8. 第三方交互先 mock gateway/service

本次不做：
- 【如真实支付、真实短信、真实 AI 网关】
```

## 4. 修复问题模板

```text
请修复【问题名称】。

上下文：
- doc/03-AI开发/00-Codex-启动入口.md
- doc/03-AI开发/01-技术方案设计.md
- doc/03-AI开发/02-AI-Coding-指导.md

问题现象：
- 【填写现象】

期望结果：
- 【填写预期行为】

要求：
1. 先定位根因，再做最小修复
2. 优先在现有模块修复，不扩散重构
3. 必要时补测试
4. 若涉及 mock/mysql 差异，说明验证方式
5. 输出包含：根因、修复点、验证结果、剩余风险
```

## 5. 小范围重构模板

```text
请对【模块名称】进行小范围重构，不改变对外行为。

目标：
- 提升可维护性
- 收敛散落逻辑

范围：
- 仅限【指定目录或文件】

要求：
1. 不改变接口契约
2. 将散落逻辑收敛到 domain / application / api adapter
3. 不扩大为全项目重构
4. 输出重构前问题、重构后结构、验证结果
```

## 6. 推荐填写字段（最少 9 项）

1. 做什么模块
2. 属于哪个端
3. 本次只做什么
4. 本次明确不做什么
5. 是否涉及第三方 mock
6. 是否需要测试 / migration / 接口
7. 是否要求复用现有模块
8. 是否涉及数据库读取设计
9. 是否涉及 `~/conf/listenbase.cof` 配置

## 7. 可直接复制示例（当前建议方向）

```text
请基于以下文档完成开发任务：

- doc/03-AI开发/00-Codex-启动入口.md
- doc/03-AI开发/01-技术方案设计.md
- doc/03-AI开发/02-AI-Coding-指导.md
- doc/02-PRD/用户端PRD.md
- doc/04-原型素材/listen_page_services.html
- doc/04-原型素材/listen_page_provider_detail.html

任务目标：
在现有代码基础上实现用户 Web 服务页与服务方详情页，以及 Go 后端服务浏览接口最小闭环。

任务范围：
- 端：用户 Web + Go 后端
- 页面：服务页、服务方详情页
- 只做：分类筛选、服务方列表、详情基础信息、服务项目列表
- 不做：下单支付、复杂推荐、第三方真实对接

实现要求：
1. 用户 Web 按 views/components/application/domain/api 分层
2. Go 后端按 interface/application/domain/infrastructure 分层
3. 补齐最小接口：分类、列表、详情、服务项目
4. 数据结构面向 MySQL 设计，必要时补 migration
5. 页面状态覆盖 idle/loading/success/empty/error
6. 输出修改文件、验证步骤、后续待办
```
