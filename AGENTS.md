# Agents

喵玉殿新版论坛（neonya）。**Go 后端 + Next.js 前端 + VitePress 文档**，单仓。

## 规则

1. 不要擅自改文件，改之前先征求用户的意见。
2. 命名时应采取完整拼写，但在真的有必要时可以向用户提案缩写。

## 目录与归属

| 路径     | 内容                | 技术栈                                                     |
| -------- | ------------------- | ---------------------------------------------------------- |
| `api/`   | Go module           | chi + Huma + pgx + sqlc + goose + River                    |
| `web/`   | 前端                | Next.js 16 App Router + React 19 + Tailwind v4 + shadcn/ui |
| `docs/`  | 文档站（VitePress） | `@neonya/docs`                                             |
| `draft/` | 草稿区，不对外链接  | —                                                          |

## 开工前必读（顺序）

1. 本文件
2. `worknote.md` —— 本轮做到哪、下一步是什么
3. `docs/source/decisions/` —— **有歧义先查这里**，已有结论不要重新发明
4. `docs/source/architecture/` —— 结构与字段的唯一权威

**文档分工**：`decisions/` 写「为什么、否决了什么」；`architecture/` 只写「是什么」。不要把理由写进架构页，也不要把字段清单写进决策记录。

## 命令

```sh
task --list          # 全部任务（每条都有 desc）
task infra:up        # podman machine start + postgres/redis
task db:migrate      # goose up
task gen             # sqlc + api/openapi.yaml + web/src/lib/api/schema.d.ts
task api:dev         # go run ./cmd/api serve（JSON 日志，可接 | jq）
task web:dev         # next dev
task smoke           # go test -tags smoke ./smoke/...
task check           # vet + test + sqlc diff + openapi diff + tsc + lint
pnpm dev:docs        # VitePress（端口 10126）
```

## 硬约束

- **安装任何依赖或工具前必须先问**（`pnpm add` / `go get` / `go install` / `scoop` / `winget`）。
- 容器只用 **Podman**，不用 Docker Desktop、不用 `docker compose`。
- 数据库只在 Podman 里跑；应用进程跑在宿主机（`go run` / `next dev`）。
- 测试连独立测试库（`TEST_DATABASE_URL` → `neonya_test`），**不用 testcontainers**。
- 不引入 hurl、Postman 集合；冒烟测试是 Go 代码。
- 提交信息用英文 + Conventional Commits（`type(scope): subject`，subject ≤ 72 字符，祈使句，结尾不加句号）。
- 升级任何依赖或工具前，先读 `docs/source/maintenance/dependency-update.md` 并按其阶段/验收清单执行。

## 代码风格

| 项         | 规则                                                                                     |
| ---------- | ---------------------------------------------------------------------------------------- |
| 命名       | 用完整单词，不缩写（`authorization` 不用 `auth`，`pageNumber` 不用 `pageNo`）            |
| Go         | tab 缩进、宽度 4，`gofmt` 标准                                                           |
| TS/JSON/MD | 空格缩进 2，无分号，单引号，无尾随逗号（见根 `.prettierrc`）                             |
| 注释       | 只在逻辑不直观时写；不写 AI 味注释；不删已有描述性注释（除非改动本身涉及该代码）         |
| Markdown   | 中文句子以「。」结尾；行过长时在句末换行；需要硬换行时行尾两个空格；提交前 `pnpm format` |

## 术语（中文文案）

- `topic` = **主题**，`post` = **条目**（量词「条」），`post_number` = **`#N`**（第 N 条）
- 禁用「楼、楼层、主楼、楼主」；「只看作者」不用「只看楼主」
- 详见 `docs/source/design/post-entry-layout.md` 第 0 节

## 工程约定摘要

完整版见 `docs/source/architecture/`。

- **契约单一来源**：OpenAPI 由 Huma 导出，前端类型由它生成；不手写 interface、不手写 fetch。
- **权限只在 Go**：前端只读 `/api/v1/me` 的能力位做 UI 显隐，不做授权判断。
- **`internal/domain` 零依赖**：不 import 框架、驱动、HTTP 包。
- **迁移纪律**：schema 变更必须带 goose 迁移；seed 幂等。
- **日志**：`slog` → JSON 到 stdout，禁止彩色美化输出。
- **XSS**：渲染只在服务端（goldmark 禁 raw HTML → bluemonday），前端只插入渲染好的 HTML。

## 每个 session 结束前

1. `task check` 通过（或明确记录未通过的部分）；
2. 更新 `worknote.md`（**简要**记本轮做了什么、下一步）；长期信息写进 `docs/`，别往便签里堆；
3. 提交 git（一条 commit 说明这一轮做了什么）。

## 新 session 开场流程

1. 读 `worknote.md` —— 看最新一节写了什么、下一步是什么，按「下一步」开干；
2. 有歧义先查 `docs/source/decisions/`，已有结论就不要再发明一遍；
3. 动手前看 `git status`：有未提交改动先问清来源，别接着别人的半成品往上写；
4. 结束时按上面「每个 session 结束前」三条收尾。

### 开场自检

动手前先自报这四项，**报不出来说明没读，不要猜**：

1. `worknote.md` 最新一节的日期与内容（做到哪、下一步）；
2. `git log -1` 的输出；
3. `git status` 是否有未提交改动；
4. 本轮实际读过哪些文件。

答案与本文件、`worknote.md` 不一致时以磁盘为准；若磁盘与用户理解不一致，**先改文档再动代码**。

## 文档维护

- 本文件是**索引**，不是手册：内容一旦在 `docs/` 有权威版本，这里只留一句指针。**超过 150 行就拆。**
- 判据：任何「未来的我还会再问一次」的问题，答案必须落盘；**只在对话里说过 = 不存在**。
- 同一事实只有一个权威位置：

  | 事实                         | 权威                                      |
  | ---------------------------- | ----------------------------------------- |
  | 为什么这么选 / 否决了什么    | `docs/source/decisions/`                  |
  | 结构 / 字段 / 索引 / 部署    | `docs/source/architecture/`               |
  | 界面元素、元素编号、中文术语 | `docs/source/design/post-entry-layout.md` |
  | 本轮进度 / 下一步            | `worknote.md`（便签）                     |
  | 阶段路线                     | `docs/source/architecture/overview.md`    |

- 状态与进度 → `worknote.md`（便签，保持简短）；阶段路线 → `architecture/overview.md`；术语 → `design/post-entry-layout.md`。这三者不另设文件。
- 文档页涨到 200 行以上再考虑拆页，不要预先切格子。
- 文档与代码改在同一处时，**同一次提交里一起改**；做不到就在对应文档里留一句「待写 / 待核」，不要只在对话里说。
