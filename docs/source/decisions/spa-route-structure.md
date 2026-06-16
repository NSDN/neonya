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
