# 喵玉殿论坛 (neonya)

喵玉殿新版论坛，monorepo 结构。

## 技术栈

| 部分 | 目录           | 技术                            |
| ---- | -------------- | ------------------------------- |
| 后端 | `apps/server/` | Go + Gin + GORM + PostgreSQL    |
| 前端 | `apps/web/`    | Vue 3 + Vite + Naive UI + Pinia |
| 文档 | `docs/`        | VitePress                       |

## 环境要求

- Node.js >= 20
- pnpm >= 10
- Go（用于后端开发）
- Podman（用于运行 PostgreSQL）

## 快速开始

### 1. 安装依赖

```sh
pnpm install
```

### 2. 配置环境变量

```sh
cp apps/server/.env.example apps/server/.env
cp apps/server/.env.postgres.example apps/server/.env.postgres
```

编辑 .env 和 .env.postgres 按需修改。

### 3. 启动数据库

不懂 Podman 的话：

```sh
./apps/server/podman-db.ps1 start
```

另外，`stop` 可关闭，`remove` 可移除。

懂 Podman 的话：

```sh
# 生成并启动
podman kube play ./apps/server/pod.yaml

# 停止并移除
podman kube down ./apps/server/pod.yaml
```

### 4. 运行数据库迁移

```sh
migrate -source file://apps/server/database/migrations \
  -database "postgres://forum_user:password@localhost:5432/forum?sslmode=disable" up
```

### 5. 启动前端

```sh
pnpm dev
```

### 6. 启动后端（另开终端）

```sh
cd apps/server && go run ./cmd/server
```

## 开发命令

```sh
# 前端
pnpm dev                     # 启动 Vite 开发服务器（端口 10123）
pnpm build                   # 类型检查 + 构建

# 后端
go run ./cmd/server          # 启动服务器（端口 10127）
go build ./...               # 编译检查

# 数据库
./apps/server/podman-db.ps1 start     # 启动 pod + 容器
./apps/server/podman-db.ps1 stop      # 停止并移除 pod
./apps/server/podman-db.ps1 status    # 查看状态
./apps/server/podman-db.ps1 remove    # 移除 pod 及所有资源

# 文档
pnpm docs:dev                # VitePress 开发服务器（端口 10126）
pnpm docs:build              # 构建文档
```

## 项目结构

```
neonya/
├── apps/
│   ├── server/               # Go 后端
│   │   ├── cmd/server/       # 入口
│   │   ├── internal/         # 业务逻辑（auth/plate/post）
│   │   ├── database/         # 迁移文件
│   │   └── pod.yml           # Podman pod 声明
│   └── web/                  # Vue 前端
│       └── src/
│           ├── features/     # 功能模块
│           └── shared/       # 共享组件/工具
├── docs/                     # VitePress 文档
├── scripts/                  # 辅助脚本
└── pnpm-workspace.yaml
```
