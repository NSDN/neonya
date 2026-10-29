# 后端架构

> 选型理由见[决策记录](/decisions/)。本文只写**结构与规则**。

## 1. 四层与依赖方向

```
transport/http (Huma operation)   解析 + 校验 + 身份提取 + 错误映射
        │
        ▼
service                           权限判断 + 事务边界 + 跨模块编排 + 入队
        │
        ▼
domain/*       entity + 业务规则 + Repository 接口（零依赖）
        ▲
        │ 实现接口
persistence (sqlc)                SQL 落地 + 行→entity + 驱动错误翻译
        │
        ▼
platform (pgx / redis / storage / mail / search / jobs)
```

**依赖只能往下走**；唯一反向箭头是 `persistence → domain`，因为它在实现 domain 定义的接口。

## 2. 各层能做与不能做

| 层 | 可以做 | 不能做 | 依赖 |
| --- | --- | --- | --- |
| `transport/http` | HTTP 编解码、Huma 校验、从 cookie 解析身份、错误映射、entity→DTO | 判断权限、开事务、拼 SQL | service、domain（只读类型） |
| `service` | 编排多个 repo、判权限、开事务、调度其他模块、入队、写审计 | 碰 `pgx` / `sqlc` 类型、解析 HTTP | domain 的接口 |
| `domain/*` | entity、业务不变式、sentinel error、Repository 接口 | import 任何框架 / 驱动 / HTTP 包 | 仅标准库 |
| `persistence` | 调 sqlc 查询、行→entity、唯一冲突翻译 | 业务判断 | domain、platform |
| `platform` | 连接池、Redis、存储、邮件、River 客户端 | 任何业务逻辑 | 第三方库 |

## 3. 强制规则

1. **`internal/domain` 零依赖**：不 import 框架、驱动、HTTP 包；只有 struct、业务规则、sentinel error。
2. **接口定义在消费方**：`domain/topic` 声明 `Repository`，`persistence/topic` 实现；service 单测用 fake。
3. **事务由 service 开启，tx 通过 context 传递**：`pgx.Tx` 不进入 repository 签名（会污染 domain）。用 `domain.Transactor`：

   - 好处是**不可能忘传**：显式传参漏一个参数，那条 SQL 就静默跑到事务外；
   - 代价是隐式：看 `repo.Insert(ctx, t)` 判不出是否在事务内。补救方法：persistence 打一条 debug 日志（`tx present / not present`），`task api:dev | jq 'select(.tx)'` 即可核对。

   ```go
   // internal/domain/tx.go —— 纯接口，不依赖 pgx
   type Transactor interface {
       WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
   }
   ```

   ```go
   // internal/platform/postgres/tx.go
   func (m *txManager) WithinTx(ctx context.Context, fn func(context.Context) error) error {
       return pgx.BeginFunc(ctx, m.pool, func(tx pgx.Tx) error {
           return fn(withTx(ctx, tx))     // tx 放进 ctx
       })
   }

   func (m *txManager) conn(ctx context.Context) gen.DBTX {
       if tx, ok := txFrom(ctx); ok { return tx }   // 事务内
       return m.pool                                // 事务外
   }
   ```

   `persistence` 每个方法开头 `sqlc.New(p.conn(ctx))` 即可，上层不感知事务。
4. **错误分层**：domain 只返回 sentinel（`ErrNoPermission` / `ErrTopicLocked` / `ErrTitleTooShort`），`transport/http` 统一映射成 `problem+json`。
5. **跨模块只走 service**：`domain/a` 不 import `domain/b` 的实现细节。
6. **装配集中**：`cmd/api/main.go` 手写全部依赖，不用 DI 容器。
7. **api 与 worker 共用 `internal/service`**：异步任务不许复制业务逻辑。
8. **运维走 CLI**：`api admin <cmd>`（`create-user` / `grant-role` / `recalc-counters` / `reindex` / `rerender`）复用同一套 service，不进 HTTP 层。

## 4. 一个垂直切片

```go
// ① transport —— 只翻译
func (h *Handler) CreateTopic(ctx context.Context, in *CreateTopicInput) (*TopicOutput, error) {
    actor, ok := auth.ActorFrom(ctx)
    if !ok { return nil, huma.Error401Unauthorized("login required") }

    tp, err := h.topics.Create(ctx, service.CreateTopicCmd{
        ActorID: actor.ID, BoardSlug: in.Body.BoardSlug,
        Title: in.Body.Title, Content: in.Body.Content,
    })
    if err != nil { return nil, mapError(err) }
    return toTopicOutput(tp), nil
}

// ② service —— 编排 + 权限 + 事务
func (s *TopicService) Create(ctx context.Context, cmd CreateTopicCmd) (*domain.Topic, error) {
    board, err := s.boards.BySlug(ctx, cmd.BoardSlug)
    if err != nil { return nil, err }
    if board.IsLocked() { return nil, domain.ErrBoardLocked }

    actor, err := s.users.ByID(ctx, cmd.ActorID)
    if err != nil { return nil, err }
    if !domain.CanCreateTopic(actor, board) { return nil, domain.ErrNoPermission }

    tp, err := domain.NewTopic(board.ID, actor.ID, cmd.Title, time.Now())
    if err != nil { return nil, err }

    return tp, s.tx.WithinTx(ctx, func(ctx context.Context) error {
        if _, err := s.topics.Insert(ctx, tp); err != nil { return err }
        if err := s.boards.BumpCounters(ctx, board.ID, +1); err != nil { return err }
        return s.jobs.EnqueueTx(ctx, jobs.IndexTopic{TopicID: tp.ID})  // 同事务入队
    })
}

// ③ domain —— 规则与契约
func NewTopic(boardID, authorID, title string, now time.Time) (*Topic, error) {
    title = strings.TrimSpace(title)
    if utf8.RuneCountInString(title) < 2 { return nil, ErrTitleTooShort }
    return &Topic{ID: uuid.NewString(), BoardID: boardID, AuthorID: authorID,
        Title: title, CreatedAt: now, LastPostedAt: now}, nil
}

// ④ persistence —— 只认 SQL
func (r *topicRepo) Insert(ctx context.Context, t *domain.Topic) (*domain.Topic, error) {
    row, err := r.conn(ctx).InsertTopic(ctx, toInsertParams(t))
    if isUniqueViolation(err) { return nil, domain.ErrDuplicate }
    if err != nil { return nil, err }
    return toDomainTopic(row), nil
}
```

## 5. 错误与 DTO 边界

```go
func mapError(err error) error {
    switch {
    case errors.Is(err, domain.ErrNoPermission):  return huma.Error403Forbidden("no permission")
    case errors.Is(err, domain.ErrTopicLocked):   return huma.Error409Conflict("topic locked")
    case errors.Is(err, domain.ErrTitleTooShort): return huma.Error422UnprocessableEntity("title too short")
    default:                                      return huma.Error500InternalServerError("internal error")
    }
}
```

- **domain entity ≠ API DTO**：`toTopicOutput()` 显式挑字段，否则新增内部字段会自动泄漏到公开响应。
- `default` 分支永不返回内部错误细节，靠 `slog` + 请求 ID 关联。

## 6. 契约生成

```
Go struct + Huma tag ──► api/openapi.yaml（提交） ──► openapi-typescript ──► web/src/lib/api/schema.d.ts
```

- 请求约束以 Go struct tag 为唯一来源（`minLength` / `maxLength` / `format`）。
- 前端不手写 fetch、不手写 interface。
- CI 重新导出 spec + `git diff --exit-code`，漂移即失败。
- `schema.d.ts` 提交进仓库，前端构建不依赖 Go 工具链。

**成本**：1 个前端 devDependency（`openapi-typescript`）、1 条 `task gen`、1 个生成文件、1 行 CI；**没有任何人手写的 YAML**。

**退出方式（两种粒度）**：

1. 想改成 spec-first：把导出的 `openapi.yaml` 冻结为手工源，Go 侧换 `oapi-codegen`，前端零改动；
2. 想彻底关掉 codegen：删导出任务与 `openapi-typescript` 步骤，前端改手写 `lib/api/types.ts`（约 10 分钟）。代价是失去跨语言编译期检查——Go 改字段时前端只在运行时才暴露问题。

## 7. 鉴权与权限

**鉴权**

- 登录成功 → Redis `session:{id}` → httpOnly + Secure + SameSite=Lax Cookie；
- 密码 argon2id；不做 JWT（可即时吊销、无 token 泄漏面）；
- 写请求校验 `Origin` 头，配合 SameSite=Lax 覆盖 CSRF。

**权限**

- 权限矩阵只在 Go（`internal/domain/shared/permissions.go`）；
- 前端只读 `GET /api/v1/me` 返回的能力位（如 `canPinTopic`）做 UI 显隐；
- **前端隐藏按钮不构成安全边界**，Go 侧一律二次校验；
- 授权判定放在 service 层，因此 CLI 与 job 能复用同一套规则。

## 8. 内容渲染与审核

**渲染管线（只在服务端）**

```
content_md ──goldmark（禁 raw HTML 直通，只放行语义标签白名单）──▶ bluemonday 白名单 ──▶ content_html
```

- URL 仅允许 `http` / `https`，拦 `javascript:` / `data:`；
- 前端**只插入服务端渲染好的 HTML**，绝不在客户端二次解析用户 Markdown（否则两套渲染器行为迟早不一致）；
- 允许的标签与禁止项见[数据模型](/architecture/data-model) §4。

**审核前置（不做纯事后清理）**

发帖链路：

1. 同步：AC 自动机过敏感词，命中则拒绝或替换；
2. 同步：`posts` 写入与可疑内容进 `reports` 在同一事务内；
3. 异步：入 River 审核队列，供人工处理（`web/admin` 只保留审核队列这类真正需要可视化的页面）。

## 9. 测试

| 层 | 方式 | 依赖 |
| --- | --- | --- |
| domain | 表驱动单测 | 无 |
| service | 注入 fake Repository / fake Transactor | 无 |
| persistence | 连 `neonya_test` 库执行真 SQL | Postgres |
| transport | `httptest` + fake service | 无 |
| 冒烟 | `api/smoke/smoke_test.go`（`//go:build smoke`） | `SMOKE_BASE_URL` |

不使用 testcontainers：测试只需一个 `TEST_DATABASE_URL`，跑前 `goose up`，结束 truncate。

## 10. 相关决策

- [扁平流 + 引用即回复](/decisions/flat-reply-model)
- [决策记录索引](/decisions/)
