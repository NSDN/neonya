# 决策记录

开发过程中做出的架构与设计决策，记录讨论脉络和最终结论。**遇到歧义时先查这里，再提新方案。**

新增决策请复制 [`_template.md`](./_template.md)，文件名用 kebab-case，并在下表登记。

## 一题一篇的决策记录

| 日期       | 决策                                                              | 简述                                                         |
| ---------- | ----------------------------------------------------------------- | ------------------------------------------------------------ |
| 2026-05-23 | [决策记录的文件组织形式](./decision-records-file-organization.md) | 一主题一文件，ADR 模板，AGENTS.md 引导查阅。                 |
| 2026-05-23 | [SPA 前端路由设计](./spa-route-structure.md)                      | ⚠️ 已失效（2026-10-30）：Next.js SSR 取代 SPA；保留「同路由 + 条件渲染」。 |
| 2026-05-24 | [前端响应式布局设计](./frontend-responsive-layout.md)             | 一套布局 + 48em 断点，移动端 BottomSheet（实现改 Next）。 |
| 2026-06-17 | [领域术语命名](./domain-term-naming.md)                           | Plate→Board；新版 `topic` / `post`（2026-10-30 修订）。       |
| 2026-06-17 | [版块列表的数据源与默认选中](./board-list-data-source.md)         | 数据库为唯一数据源，删除前后端硬编码。                       |
| 2026-06-17 | [HTTP 层换用原生 fetch](./http-native-fetch.md)                   | ⚠️ 已失效（2026-10-30）：改用 `openapi-fetch`；「零 HTTP 依赖」结论不变。 |
| 2026-06-17 | [API 路径命名约定](./api-url-naming.md)                           | 相对路径 + `new URL()`；前缀与版本见《API 版本与路径》。 |
| 2026-06-17 | [Guard 条件提取模式](./guard-condition-extraction.md)             | ⚠️ 已失效（2026-10-30）：TanStack Query 的 `enabled` 取代。 |
| 2026-06-18 | [浮动按钮组设计](./floating-action-group.md)                      | 三按钮组 + BottomSheet；移动端发帖入口定案为 FAB。 |
| 2026-06-20 | [Store 与 Composable 职责边界](./store-composable-separation.md)  | ⚠️ 已失效（2026-10-30）：改为前端状态分层三层。 |
| 2026-10-29 | [依赖更新的分阶段策略](./dependency-update-policy.md)             | 一阶段一提交；流程正文移到 `maintenance/`。                  |
| 2026-10-30 | [扁平流 + 引用即回复](./flat-reply-model.md)                      | 取消楼中楼与树形链视图；回复=引用 + 被引用索引，嵌套恒为 1。 |
| 2026-10-30 | [API 版本与路径](./api-versioning.md)                             | `/api/v1` 前缀；浏览器用相对路径，Server Component 用容器内地址。 |
| 2026-10-30 | [前端状态分层](./frontend-state-layers.md)                        | Server Component / TanStack Query / 局部 state 三层，不引全局状态库。 |

## 技术选型一览

尚未展开成单篇的决策，先在此登记结论与被否决项。**需要追溯讨论过程时再补写一篇。**

| 决策 | 被否决的选项 | 结论理由 |
| --- | --- | --- |
| Go + TypeScript 双语言 | 全 TS（Next + NestJS） | 论坛是服务端重业务，Go 部署与并发更省心 |
| Huma 代码优先 | 手写 `openapi.yaml` + oapi-codegen | 无前后端并行需求，手写 YAML 是纯税且易漂移 |
| chi + Huma | Gin / Echo / Fiber | chi 纯粹 `net/http`，中间件与 `httptest` 直接可用 |
| sqlc + goose | Prisma / GORM / ent | 跨语言不可用；ORM 隐式行为多、查询不可控 |
| 事务：`Transactor` + ctx 传 tx | `pgx.Tx` 显式传入 repository 方法 | 保住 domain 零依赖；显式传参漏一个参数会让 SQL 静默跑到事务外 |
| River | asynq / Kafka / 无队列 | 与业务写同事务入队，量级不需要独立 MQ |
| `pg_trgm` 起步 | 直接上 Meilisearch | 少一个中间件；中文分词不足时再升级 |
| Next.js | Nuxt / SvelteKit | SSR 生态与资料量最大；Server Component 直连 Go |
| pnpm workspace（`web` + `docs`） | 各包独立 lockfile / Turborepo | 两个 JS 包共享一套依赖与格式化配置；包数少，不需要 task runner |
| Cookie 会话 | JWT | 可即时吊销、无 token 泄漏面、同源部署下更简单 |
| 扁平流 + 引用即回复 | 楼中楼两层 / 无限嵌套 / 树形链视图 | 见[独立 ADR](./flat-reply-model.md) |
| 引用存文本快照 | 只存被引 ID | 原帖被编辑/删除后引用仍可读 |
| 中文文案用「条目」 | 楼 / 帖子 / 评论 | 「楼」是楼层文化；「帖」与「主题」撞车；「评论」弱化指向性 |
| Podman | Docker Desktop / `docker compose` | 无宿主守护进程；`machine stop` 即彻底关闭，不需管理员 |
| dev 裸 `podman run` | compose provider | 开发只需两个基础设施容器，不值得引 provider |
| prod 用 Quadlet | docker compose / `podman-compose` | systemd 直接托管，自带开机自启与重启策略 |
| 测试用独立测试库 | testcontainers | rootless Podman 下需配 socket 并禁用 Ryuk |
| shadcn/ui | MUI / Ant Design | 源码进仓库，可 diff、可改，无大版本升级压力 |
| 运维 CLI 子命令 | 后台管理页面先行 | 终端里不必先开浏览器；后台只保留真正需要可视化的页面 |
| `slog` JSON | zap / zerolog / 彩色输出 | 标准库够用；JSON 便于 `jq` 与日志采集 |
| OpenAPI codegen | 手写 TS 类型 / gRPC | 跨语言唯一的廉价编译期检查，spec 是 Huma 副产物 |
| 术语 `topic` / `board` | `thread` / `category` / `forum` / `node` | 见[领域术语命名](./domain-term-naming.md) |
| 内容存 Markdown | ProseMirror JSON / 存 HTML | 可读可 diff；源文在即可全量重渲染 |
| 语义标签白名单 | 开放内联 HTML（含 `style`） | `style` / `class` 等于允许用户改站点 UI |
| WYSIWYG 输出 Markdown | 富文本另存 HTML | 同一内容两种表示必然漂移 |
| UUID v7 + `/t/{id}` | slug URL / 自增 bigint | 标题改动会断外链；自增 ID 暴露体量 |
