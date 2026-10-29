# 依赖与工具更新流程

> 策略与取舍见[依赖更新的分阶段策略](/decisions/dependency-update-policy)。本文只写**怎么操作**。
> 本文按**当前仓库现状**编写：目前只有 `docs/` 与根配置，`api/`、`web/` 尚未创建。

## 1. 现在能更新什么

| 部分 | 现状 | 可更新对象 | 当前版本 |
| --- | --- | --- | --- |
| `docs/` | 已存在 | `vitepress`、`vue` | vitepress 1.6.x、vue 3.x |
| 根 `package.json` | 已存在 | `prettier` | 3.9.x |
| `web/` | **未创建** | — | — |
| `api/` | **未创建** | — | — |
| 宿主工具链 | 已安装 | pnpm、node、podman、task | pnpm 12.15.0、node 26、podman 6.1.4 |

`docs` 用 Vue、将来 `web` 用 React，两套框架互不影响，可以分别升级。

## 2. 操作流（五步）

```powershell
# 0. 前置：工作区必须干净，单独分支
git status --short          # 有输出就先处理，不要带着未提交改动升级
git switch -c chore/deps-YYYYMMDD

# 1. 看有哪些过期
pnpm outdated -r

# 2. 分阶段更新（一次只做一个阶段，见 §3）
pnpm --filter @neonya/docs add -D vitepress@latest vue@latest

# 3. 验收（见 §4）
pnpm install
pnpm build:docs
pnpm dev:docs               # 手点首页 / 决策记录 / 条目布局
pnpm format:check

# 4. 提交（一阶段一提交，锁文件必须一起提交）
git add -A
git commit -m "chore(deps): bump vitepress and vue"

# 5. 出问题时整阶段回滚
git revert <sha>
```

## 3. 阶段顺序（按风险递增）

| 阶段 | 范围 | 判断依据 |
| --- | --- | --- |
| 1 | 补丁 / 次要版本 | `pnpm -r update`，无代码改动 |
| 2 | 文档站构建链 | `vitepress` + `vue` + `vite`，三者成套升 |
| 3 | 根工具（prettier 等） | 与构建无关，随时可升 |
| 4 | `web` 构建链 | Next + React + TypeScript，peer 互相约束，**必须同一提交**（`web/` 出现后适用） |
| 5 | Go 侧 | `go get -u ./...` 后跑 `task check`（`api/` 出现后适用） |
| 6 | 宿主工具链 | pnpm / node / podman / task，见 §6 |
| 7 | 清理 | 弃用依赖、lockfile 冗余 |

## 4. 每阶段的验收清单

- [ ] `pnpm install` 无 `ERR_PNPM_IGNORED_BUILDS`
- [ ] `pnpm build:docs` 通过（**死链检查发生在这里**，构建失败往往意味着某个文档链接写错了）
- [ ] `pnpm dev:docs` 打开 `/`、`/decisions/`、`/design/post-entry-layout` 三个页面正常
- [ ] `pnpm format:check` 通过
- [ ] 一条 commit，前缀 `chore(deps):` 或 `chore(tooling):`

`api/` 出现后追加：`go vet ./...`、`go test ./...`、`task check`。

## 5. 现状特有的坑

1. **`onlyBuiltDependencies`**：`pnpm-workspace.yaml` 里显式放行了 `esbuild`。若新依赖也需要跑安装脚本，`pnpm install` 会报 `ERR_PNPM_IGNORED_BUILDS`——用 `pnpm approve-builds <pkg>` 批准，并把 `pnpm-workspace.yaml` 的改动一起提交。**不要**用 `--dangerously-allow-all-builds`。
2. **`pnpm-lock.yaml` 必须同提交**，否则不同机器的安装结果不一致。
3. **VitePress 手写侧边栏**：`docs/.vitepress/composition/use-sidebar.ts` 是自己维护的，主版本升级若改了配置 API，构建会直接失败——靠阶段 2 的 `pnpm build:docs` 兜住。
4. **`packageManager` 字段**：`pnpm self-update` 后要同步改根 `package.json` 的 `packageManager`，否则 CI 会用旧版本。
5. **`engines.node`**：升 Node 主版本时同步改这里与 CI 配置，否则本地能跑、CI 挂。

## 6. 宿主工具链更新

```powershell
pnpm self-update                    # 升完跑一次 pnpm install + pnpm build:docs
scoop update task                   # 升完跑 task --list 确认任务都在
winget upgrade Podman               # 升完 podman machine stop/start，再 task infra:up 验证
```

Node 主版本升级：改 `engines.node`、根 `package.json` 的 `packageManager`、以及将来 CI 的 node 版本，三者一起提交。

## 7. 尚未覆盖

以下等内容出现后再补本文件：

| 项 | 触发条件 |
| --- | --- |
| Next + React + TypeScript 升级流程 | `web/` 创建后 |
| Go 模块升级（pgx / chi / huma / sqlc / goose） | `api/` 创建后 |
| 数据库与缓存镜像钉版本（`postgres:16-alpine` / `redis:7-alpine`） | `Taskfile.yaml` 出现后 |
| 升级后回填「版本矩阵」与「已完成更新记录」 | 首次完整升级后 |
