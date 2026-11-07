# 架构总览

> 本文只写**是什么**。为什么这么选、否决了什么方案，见[决策记录](/decisions/)。
>
> 运行前提：日常操作全部走命令行（`docker` / `task` / `go` / `pnpm` / `psql`），不依赖常驻 GUI。

## 1. 全景

```
browser
   │ HTTPS
   ▼
 nginx ──┬── /api/* ──────► api     (Go, chi + Huma)  ──┬──► PostgreSQL 16
         │                                               │
         └── /  ──────────► web     (Next.js, SSR)       ├──► Redis 7
                                                         │
                               worker  (Go) ─── River jobs ┘
```

- 前后端**同源**：Nginx 按路径分流，浏览器视角只有一个域名，因此不需要 CORS，Cookie 天然同父域。
- `api` 与 `worker` 是同一个 Go module 的两个二进制，共用 `internal/service`。
- 异步任务（邮件、缩略图、通知、计数落库）通过 **River** 入队，与业务写同事务。

## 2. 技术栈

| 层          | 选型                                   | 备注                                                 |
| ----------- | -------------------------------------- | ---------------------------------------------------- |
| 后端语言    | Go 1.27                                | 单静态二进制                                         |
| HTTP 路由   | chi v5                                 | 纯 `net/http` 兼容                                   |
| API 层      | Huma v2（`humachi` 适配器）            | 从 Go struct 生成 OpenAPI 3.1 + 请求校验             |
| 契约        | OpenAPI 3.1（Huma 导出，提交进仓库）   | 前端类型由它生成                                     |
| 数据库      | PostgreSQL 16                          | 部分索引、JSONB、`pg_trgm`                           |
| 驱动 / 查询 | pgx v5 + sqlc                          | 手写 SQL + 类型安全生成                              |
| 迁移        | goose（纯 SQL，可嵌入二进制）          | 版本化、可回滚                                       |
| 任务队列    | River（基于 Postgres）                 | 与业务写同事务入队                                   |
| 缓存/会话   | Redis 7 + go-redis                     | TTL、计数器、限流                                    |
| 附件存储    | 本地磁盘 → S3 兼容                     | 走 `Storage` 接口                                    |
| 搜索        | `pg_trgm` + GIN（MVP）                 | 中文分词不足时再上 Meilisearch                       |
| 前端        | Next.js 16 App Router + React 19       | SSR/SEO，Server Component 直连 Go                    |
| 前端数据层  | openapi-fetch + TanStack Query         | 类型来自生成的 `schema.d.ts`                         |
| 样式        | Tailwind CSS v4                        | —                                                    |
| UI 组件     | shadcn/ui（源码进仓库）                | CLI 安装，可改可 diff                                |
| 容器        | WSL 发行版内的 Docker Engine + compose | Windows 侧不装 `docker.exe`；dev / prod 共用 compose |
| 包管理      | pnpm workspace（`web` + `docs`）       | 两个独立 JS 包共用一个 lockfile                      |
| 任务编排    | Taskfile                               | `task --list` 当文档用                               |
| 日志        | `slog` → JSON 到 stdout                | `jq` 可读                                            |

## 3. 目录结构

```
neonya/
├─ api/                      # Go module
│  ├─ cmd/{api,worker,migrate}/
│  ├─ internal/{domain,service,persistence,transport,platform,config}/
│  ├─ openapi.yaml           # Huma 导出，提交进仓库
│  ├─ gen/                   # sqlc 生成物
│  ├─ queries/  migrations/  smoke/
│  └─ go.mod
├─ web/                      # Next.js
│  ├─ src/app/  src/features/  src/components/{ui}/  src/lib/api/
│  └─ package.json
├─ docs/                     # VitePress 文档站（本站）
│  ├─ .vitepress/  types/  source/
│  └─ package.json
├─ deploy/{nginx,Dockerfile.api,Dockerfile.web,compose.prod.yaml}/
├─ draft/                    # 草稿区，不对外链接
├─ Taskfile.yaml
└─ .github/workflows/ci.yml
```

各层内部的目录约定见 [后端](/architecture/backend)、[前端](/architecture/frontend)。

## 4. 取舍原则

按「改起来的成本」分级，只在成本最低的地方冒险：

| 级别           | 内容                                    | 策略         |
| -------------- | --------------------------------------- | ------------ |
| 换掉 < 1 天    | 搜索、存储、队列、日志、样式、UI 组件   | 选最省事的   |
| 换掉 1~2 周    | 路由框架、前端框架、ORM 形态、单体↔拆分 | 选文档最多的 |
| 换掉要重做数据 | 数据模型、契约机制、鉴权模型、SEO 归属  | 选最保守的   |

## 5. 明确的「不做」

- 不用 ORM（GORM / Prisma / ent）：查询复杂度不可控，见[决策记录](/decisions/)。
- 不用 GraphQL / tRPC：论坛是典型 REST 资源模型。
- 不用 DI 框架（wire / fx）：`cmd/api/main.go` 手写装配，一屏看清依赖链。
- 不用 Docker Desktop，也不在 Windows 侧装 `docker.exe`：容器只在 WSL 发行版内跑，见 [部署](/architecture/deployment)。
- 不用 Postman 集合：冒烟测试写成 Go 代码。
- 不用 npm 包型组件库（MUI / AntD）：shadcn/ui 源码进仓库。
- 不用 testcontainers：测试连独立测试库即可。
- 不引入 MinIO / Meilisearch / 私信，直到确实需要。

## 6. 跨页约定

| 约定           | 内容                                                                                     |
| -------------- | ---------------------------------------------------------------------------------------- |
| 契约单一来源   | API 与前端类型都从 Huma 导出的 OpenAPI 生成，禁止手写重复定义                            |
| 权限唯一实现   | 权限判定只在 Go 侧，前端只消费 `/api/v1/me` 返回的能力位                                 |
| 迁移纪律       | schema 变更必须带 goose 迁移，禁止手改生产库                                             |
| 日志与可观测性 | 结构化 JSON 到 stdout（禁止彩色美化）；请求 ID 贯穿 API → River job；发帖 / 登录链路打点 |
| 提交前         | `task check`（go vet / test / sqlc / openapi `-check` / tsc / lint）                     |

## 7. 阶段路线

| 阶段    | 内容                                                                                                          | 完成标志                                                                |
| ------- | ------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| M0 骨架 | 目录、docker infra、迁移、Huma 起服务、Next 起页面、契约生成链路、Go 冒烟                                     | `task infra:up` → `task db:migrate` → `task dev` → 页面显示 health 数据 |
| M1 核心 | 注册登录、版块、发主题/回复、扁平流 + 引用块 + 被引用索引、`topic_reads` 未读定位、Markdown、图片上传、软删除 | 能在页面上发一条，看到引用块与 `#N`                                     |
| M2 运营 | 权限与版主、举报审核、敏感词、通知（@与回复）、审计日志、后台                                                 | 举报进队列；`api admin grant-role` 生效                                 |
| M3 优化 | 搜索升级、计数落库、缓存策略、SEO、RSS、邮件订阅                                                              | —（暂不排期）                                                           |

每个阶段结束必须满足：`task check` 通过、有可演示的页面或接口、文档与代码同一次提交。
阶段内部的单元切分记在本节，session 进度记在 `worknote.md`。

### M0 剩余切分（`web/`）

每单元做完即停，等审阅；一次交付一个单元。

| 单元 | 内容                                                                  | 验收                                       |
| ---- | --------------------------------------------------------------------- | ------------------------------------------ |
| W1   | Next 16 脚手架（App Router + TS + Tailwind v4），无 API 调用          | `pnpm --filter @neonya/web dev` 打开默认页 |
| W2   | `next.config.ts` 的 `rewrites`：`/api/v1/*` → `http://localhost:8080` | 页面内 `fetch('/api/v1/health')` 通        |
| W3   | `task gen` 第二步：`openapi-typescript` → `schema.d.ts`               | `task gen` 后 `tsc --noEmit` 通过          |
| W4   | `lib/api/client.ts`（浏览器相对路径 / 服务端 `API_BASE_URL`）         | `tsc --noEmit`                             |
| W5   | 首页 Server Component 拉 `/health` 并渲染                             | 页面显示 `database: up`                    |
| W6   | `task web:dev`、`task check` 接上 `tsc` + `lint`                      | `task check` 通过                          |
| W7   | `.github/workflows/ci.yml`（api / web 两个 job）                      | 首次运行绿                                 |

## 8. 子页

| 页面                                        | 内容                                       |
| ------------------------------------------- | ------------------------------------------ |
| [后端](/architecture/backend)               | 分层、事务、错误映射、契约生成、鉴权与权限 |
| [数据模型](/architecture/data-model)        | 表、字段、索引、内容格式、引用模型         |
| [前端](/architecture/frontend)              | 目录、渲染策略、数据层、编辑器、组件       |
| [部署](/architecture/deployment)            | Docker、compose、Nginx、环境变量、CI       |
| [条目布局与交互](/design/post-entry-layout) | 通栏线框图、元素编号、中文术语             |
