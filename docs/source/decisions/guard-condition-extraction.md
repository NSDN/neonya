# 2026-06-17 Guard 条件提取模式

## 背景

`useBoards()` 中 `handleGetBoards` 的防重复请求有两个独立的 guard 条件：
`boardsStore.loading`（防止并发）和 `boards.value.length > 0`（已获取过数据）。
两者行为相同（提前 return），但本质不同——前者是临时状态，后者是永久状态。

## 讨论脉络

- **方案 A：`||` 直接合并** — `if (a || b) return`；
  简洁但语义混在一起，后续不易扩展。
- **方案 B：提取为命名函数** — `function shouldSkipFetch() { ... }`；
  调用处一行，加条件只改函数体。
- **方案 C：predicate 数组** — `[a, b].some(Boolean)`；
  条件可任意增减，结构不变。

## 决策

**B + C 结合：提取为 `shouldSkipFetch()` 函数，内部用 `conditions.some(Boolean)` 数组**。

```typescript
function shouldSkipFetch() {
  const conditions = [boardsStore.loading, boards.value.length > 0]
  return conditions.some(Boolean)
}
```

调用处：`if (shouldSkipFetch()) return`。

新增条件时只需在 `conditions` 数组里加一行，调用处不受影响。

## 影响

- `apps/web/src/features/board/composables/useBoards.ts` — guard 条件以此模式书写。
- 同项目内其他组合式函数可复用此模式。

---

## 2026-10-30 修订：已失效

新版前端不再手写「该不该请求」的守卫：首屏取数在 Server Component 里直接 `await`，客户端状态由 TanStack Query 管理，防重复请求用 `enabled` / `staleTime`，见[前端架构](/architecture/frontend) 与[前端状态分层](./frontend-state-layers)。

`shouldSkipFetch()` + `conditions.some()` 这套模式已经没有对象了。

**仍可沿用的部分**：把多个同类条件提取成命名函数或 predicate 数组，是一条通用编码习惯，将来遇到类似场景可以照此写，但不再是项目级约定。
