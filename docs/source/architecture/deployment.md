# 部署与运行

> 为什么用 Podman 而不是 Docker Desktop 见[决策记录](/decisions/)。

## 1. 开发环境：只跑两个基础设施容器

**应用进程跑在宿主机**（`go run` / `next dev`），容器只跑 Postgres 与 Redis。

- 不把仓库 bind mount 进 VM：跨文件系统 I/O 慢，rootless 下还有 UID 映射问题；
- 不需要 `podman compose`：Windows 上它要额外 provider，且语义与 Docker 有差异。

| 任务 | 行为 | 数据 |
| --- | --- | --- |
| `task infra:up` | `podman machine start` + 起已有容器（不存在才 `podman run` 创建） | 保留 |
| `task infra:down` | `podman machine stop` | 保留 |
| `task infra:reset` | 删容器 + 删命名卷 + `podman machine rm` | **销毁** |

约定：**镜像名一律全限定**（`docker.io/library/postgres:16-alpine`），避免 short-name 解析要求交互确认。

```powershell
podman run -d --name neonya-pg -p 5432:5432 `
  -e POSTGRES_USER=neonya -e POSTGRES_PASSWORD=neonya -e POSTGRES_DB=neonya `
  -v neonya-pgdata:/var/lib/postgresql/data `
  docker.io/library/postgres:16-alpine

podman run -d --name neonya-redis -p 6379:6379 `
  -v neonya-redisdata:/data `
  docker.io/library/redis:7-alpine redis-server --appendonly yes
```

`podman machine` 建议 4GiB 内存（Postgres + Redis + `next build` 同时跑）：

```powershell
podman machine stop
podman machine set --memory 4096
podman machine start
```

## 2. 生产：Nginx + Quadlet

```
Nginx
 ├─ location /api/  → api:8080   (Go 静态二进制，CGO_ENABLED=0)
 └─ location /      → web:3000   (Next standalone)
```

同源分流带来三个结果：无需 CORS、Cookie 天然同父域、前端可用相对路径 `/api/v1/...`。

编排用 **Quadlet**（systemd user unit），不用 compose：

```
~/.config/containers/systemd/neonya-api.container
~/.config/containers/systemd/neonya-worker.container
~/.config/containers/systemd/neonya-web.container
~/.config/containers/systemd/neonya-pg.container
```

```sh
systemctl --user enable --now neonya-api
systemctl --user status neonya-api
```

开机自启与重启策略交给 systemd，不需要 compose 或外部 supervisor。

## 3. 环境变量

| 变量 | 用于 | 示例 |
| --- | --- | --- |
| `DATABASE_URL` | api / worker / migrate | `postgres://neonya:***@localhost:5432/neonya?sslmode=disable` |
| `REDIS_URL` | api / worker | `redis://localhost:6379/0` |
| `HTTP_ADDR` | api | `:8080` |
| `API_BASE_URL` | web（Server Component 取数） | `http://localhost:8080` 或 `http://api:8080` |
| `SESSION_COOKIE_NAME` | api | `neonya_session` |
| `STORAGE_DIR` | api | `./var/uploads` |
| `SMTP_URL` | worker | `smtp://user:pass@host:587` |
| `APP_ENV` | 全部 | `dev` / `prod` |

复制 `.env.example` → `.env` 后填写；`.env` 不入库。

## 4. 数据库迁移

```sh
task db:migrate        # goose up
task db:rollback       # goose down 一步
task db:seed           # 幂等种子数据
```

- 迁移文件在 `api/migrations/`，纯 SQL，命名 `0001_xxx.sql`；
- 生产不跑 `seed`；`goose` 以库内 `goose_db_version` 表追踪版本。

## 5. CI

两个 job，按路径触发：

| job | 内容 |
| --- | --- |
| `api` | `go vet` / `go test` / `sqlc diff` / `openapi diff` |
| `web` | `tsc --noEmit` / eslint / `next build` |

CI 用 GitHub Actions 的 `services`（真 Docker），与本地 Podman 无关；测试库地址由 `TEST_DATABASE_URL` 注入。

## 6. 常用任务一览

```sh
task --list          # 全部任务及说明
task infra:up        # 起机器的 pg/redis
task db:migrate      # 迁移
task gen             # sqlc + openapi.yaml + schema.d.ts
task api:dev         # go run ./cmd/api serve   (JSON 日志，可 | jq)
task web:dev         # next dev
task smoke           # go test -tags smoke ./smoke/...
task check           # vet + test + diff + tsc + lint
```
