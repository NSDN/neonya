# 2026-05-23 SPA 前端路由设计

## 背景

重构版喵玉殿论坛为 SPA 架构。论坛的层级结构为版块（board）→ 帖子（thread）两级。版块分为普通文字版块（article）和漫画版块（comic），文字版块以标题列表显示帖子，漫画版块以缩略图形式显示帖子。

## 讨论脉络

### 路由结构

- **方案 A：按类型分路由**（如 `/articles/:id` vs `/comics/:id`）
  - 缺点：URL 不直观，用户被迫感知版块类型；同一概念资源被不同路径表达。
- **方案 B：同一路由 + 条件渲染**（`/boards/:id` 内根据数据的 `pageType` 切换组件）
  - 优点：URL 统一，版块类型由服务端数据决定，前端无需维护额外路由规则。

### 路由模式

- Hash 模式：浏览器只请求 `/`，无需服务端配合，URL 带有 `#`。
- History 模式：URL 干净自然，需服务端将未匹配路径 fallback 到 `index.html`。

## 决策

1. **路由结构**：版块列表统一用 `/boards/:id`，帖子详情统一用 `/threads/:id`（待建）。
   组件内部根据数据的类型字段（`pageType` / `threadType`）条件渲染不同视图。
2. **路由模式**：从 Hash 切换到 History。
   前端 `createWebHistory` 替换 `createWebHashHistory`；
   后端 Gin 添加 catch-all 路由，未匹配路径返回 `index.html`。

## 影响

- 前端路由：追加 `/threads/:id` 路由定义。
- 前端 router 实例：`createWebHashHistory` → `createWebHistory`。
- 后端：添加 catch-all fallback 路由。

---

## 2026-10-30 修订：已失效

新版前端改为 **Next.js App Router（SSR）**，本决策的三项内容全部失效：

| 原决策 | 现在 |
| --- | --- |
| SPA（单页应用） | Next.js App Router 服务端渲染；所有写操作打 Go API，见[前端架构](/architecture/frontend) |
| Hash → History 模式 | 不存在「前端路由模式」概念；路由即文件（`app/f/[slug]/page.tsx`） |
| 后端 Gin catch-all 落回 `index.html` | 页面由 Next 渲染，Nginx 只做 `/api/*` 与 `/` 的分流，见[部署](/architecture/deployment) |

**保留的内核**：「同一路由 + 按类型条件渲染」这一步仍然有效——URL 统一为 `/t/{topicId}`，视图由 `topics.kind`（`0 = normal` / `1 = comic`）决定，而不是拆成 `/articles/:id` 与 `/comics/:id`（该否决结论未变）。字段名由 `pageType` / `threadType` 改为 `topics.kind`，见[数据模型](/architecture/data-model)。
