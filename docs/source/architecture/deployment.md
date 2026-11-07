# 部署与运行

> 为什么把容器放进 WSL 内的 Docker、否决了哪些方案，见[决策记录](/decisions/)。

## 1. 开发环境：只跑两个基础设施容器

**应用进程跑在宿主机**（`go run` / `next dev`），容器只跑 Postgres 与 Redis。

容器与 docker 客户端都在 WSL 发行版（`Ubuntu`）内，Windows 侧不装 `docker.exe`：

- 客户端与 daemon 同处一个 Linux 内核，`compose.yaml` 的相对路径按 Linux 语义解析，不需要路径翻译；
- 不把仓库 bind mount 进容器：跨文件系统 I/O 慢，drvfs 的权限语义也不完整；
- Postgres / Redis 的数据一律放命名卷（原生 ext4），不落 `/mnt/...`。

| 任务               | 行为                                                                  | 数据     |
| ------------------ | --------------------------------------------------------------------- | -------- |
| `task infra:up`    | `docker compose up -d --wait` + 建测试库 `neonya_test`                | 保留     |
| `task infra:psql`  | `docker exec -it neonya-pg psql -U neonya -d neonya`                  | 只读     |
| `task infra:down`  | `docker compose stop`；发行版空闲约 15 秒后自行退出，dockerd 随之停止 | 保留     |
| `task infra:reset` | `docker compose down -v`                                              | **销毁** |

约定：**镜像名一律全限定**（`docker.io/library/postgres:16-alpine`）。

约定：**dockerd 不随 systemd 自启**——`docker.service` / `containerd.service` 均为 `disabled`，只留 `docker.socket` 做 socket 激活。
否则任意一次 `wsl` 命令都会把 daemon 拉起来，因此 compose 里容器一律 `restart: 'no'`。

内存上限交给 `%USERPROFILE%\.wslconfig`（默认取宿主内存的 50%，pg + redis 用不到那么多）：

```ini
[wsl2]
memory=16GB
[experimental]
sparseVhd=true
```

## 2. 生产：Nginx + Docker Compose

```
Nginx
 ├─ location /api/  → api:8080   (Go 静态二进制，CGO_ENABLED=0)
 └─ location /      → web:3000   (Next standalone)
```

同源分流带来三个结果：无需 CORS、Cookie 天然同父域、前端可用相对路径 `/api/v1/...`。

编排用 **Docker Compose**（`compose.prod.yaml`），由 systemd 的一个 oneshot unit 在开机时执行 `docker compose -f compose.prod.yaml up -d`，容器自身用 `restart: unless-stopped`：

```sh
sudo systemctl enable --now neonya.service
sudo systemctl status neonya.service
```

开机自启与重启策略交给 systemd + compose，不引外部 supervisor。

> 待核：`compose.prod.yaml` 与 api / worker 镜像尚未落地（M1 起才有）。`Quadlet` 是 Podman 专有，已随 Podman 一并弃用。

## 3. 环境变量

| 变量                  | 用于                         | 示例                                                                          |
| --------------------- | ---------------------------- | ----------------------------------------------------------------------------- |
| `DATABASE_URL`        | api / worker / migrate       | `postgres://neonya:***@localhost:5432/neonya?sslmode=disable`                 |
| `REDIS_URL`           | api / worker                 | `redis://localhost:6379/0`                                                    |
| `HTTP_ADDR`           | api                          | `:8080`                                                                       |
| `API_BASE_URL`        | web（Server Component 取数） | `http://localhost:8080` 或 `http://api:8080`                                  |
| `SESSION_COOKIE_NAME` | api                          | `neonya_session`                                                              |
| `TRUSTED_PROXY_COUNT` | api                          | `0`（直连，不解析 `X-Forwarded-For`）/ `1`（Nginx 后）；填小了客户端可伪造 IP |
| `STORAGE_DIR`         | api                          | `./var/uploads`                                                               |
| `SMTP_URL`            | worker                       | `smtp://user:pass@host:587`                                                   |
| `APP_ENV`             | 全部                         | `dev` / `prod`                                                                |

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

| job   | 内容                                               |
| ----- | -------------------------------------------------- |
| `api` | `go vet` / `go test` / `sqlc diff` / `spec -check` |
| `web` | `tsc --noEmit` / eslint / `next build`             |

CI 用 GitHub Actions 的 `services`（真 Docker），与本地 WSL 内的 Docker 无关；测试库地址由 `TEST_DATABASE_URL` 注入。

## 6. 常用任务一览

```sh
task --list          # 全部任务及说明
task infra:up        # 起 pg/redis 容器
task db:migrate      # 迁移
task gen             # sqlc + openapi.yaml + schema.d.ts
task api:dev         # go run ./cmd/api serve   (JSON 日志，可 | jq)
task web:dev         # next dev
task smoke           # go test -tags smoke ./smoke/...
task check           # 格式 + vet + test + 契约校验（tsc / lint 待 web/ 出现后接入）
```

### M0 自检

```powershell
task infra:up        # 起 pg / redis，等健康检查通过
task db:migrate      # 期望输出：current version: 1
task api:dev         # 另开终端
curl.exe -s http://localhost:8080/api/v1/health
```

期望 `200` + `{"state":"ok","database":"up",…}`；停掉容器再请求应变成 `503` + `degraded`。

`task check` 里的契约校验走 `go run ./cmd/api spec -check -output openapi.yaml`，**只读不写盘**：不一致时报出第几行，并提示跑 `task gen` 刷新。
