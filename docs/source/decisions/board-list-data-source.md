# 2026-06-17 版块列表的数据源与默认选中

## 背景

版块列表需要在多种场景下使用（侧边栏显示、路由重定向、帖子列表过滤），
关于版块数据应该写死在前端代码、localStorage 还是数据库，需要确定唯一数据源。

## 讨论脉络

- **方案 A：写死在前端代码** — 无需 API 调用，即时可用；但新增/修改版块需要重新部署。
- **方案 B：存在 localStorage** — 首次加载后免请求；但不是真正的数据源，各用户版本不一致。
- **方案 C：存在数据库** — 论坛业界的标准做法，
  版块属于共享内容而非前端配置；
  管理员可动态增删版块无需重新部署；
  前后端共用单一数据源。

## 决策

**选用数据库作为唯一数据源**。

- 数据库 migration 负责播种默认版块。
- 删除 Go 端的 `DefaultBoards` 硬编码变量和启动时的 `InitBoardList()` 函数（双重播种冗余）。
- 删除前端的 `BOARD_IDS` 硬编码常量，`BoardId` 类型从 union type 改为 `string`，
  新增版块无需修改 TypeScript 类型。
- 版块列表未被选中时的默认版块从硬编码 `'localization'` 改为动态取列表第一项，
  fallback 为 `boardsStore.boards[0]?.id ?? 'localization'`。

## 影响

- `apps/server/internal/board/model.go` — 不再有 `DefaultBoards`。
- `apps/server/internal/board/service.go` — 不再有 `InitBoardList()`。
- `apps/web/src/features/board/types.ts` — `BoardId = string`，删除了 `BOARD_IDS`。
- `apps/web/src/features/board/composables/useBoards.ts` — 默认选中取列表第一项。
- `apps/web/src/features/board/composables/useCurrentBoard.ts` — fallback 改为动态。
