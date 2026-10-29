# 2026-10-30 前端状态分层

## 背景

新版前端用 Next.js App Router（Server Component）与 TanStack Query，不再有 Pinia。
原 [Store 与 Composable 职责边界](./store-composable-separation) 的对象消失，需要一套替代分层。

## 讨论脉络

- **全部用 TanStack Query**：写法统一；但首屏要等客户端 JS，丢掉 SSR 的 SEO 与首屏速度优势。
- **全部用 Server Component**：首屏最快；但无限滚动、点赞、发帖这些交互必须回到客户端。
- **三层混合**（采纳）：首屏走服务端，交互走查询缓存，纯 UI 状态留在组件。
- **引入全局状态库**（Zustand / Jotai）：灵活；但会与查询缓存重复。当前真正跨页面共享的客户端状态只有「当前用户与能力位」，而它可以由服务端注入。

## 决策

| 层 | 用途 | 规则 |
| --- | --- | --- |
| Server Component + `await` | 首屏数据、SEO 内容 | 只读不写；不调用 `useState` |
| TanStack Query | 客户端缓存、失效、乐观更新、无限滚动 | 所有客户端写操作的唯一入口 |
| `useState` / `useReducer` | 组件局部表单与 UI 状态 | 不跨组件共享 |
| 服务端注入（`/api/v1/me`） | 当前用户与能力位 | 通过 props / context 下发，不在客户端重复请求 |

**不引入全局状态库**，直到出现两个以上跨页面共享的客户端状态。

判断标准沿用原决策：不看「有没有副作用」，看**数据的共享范围**。

## 影响

- `web/src/features/*`：每个域导出两类入口——服务端取数函数（供 Server Component 调用）与客户端 query hook。
- 权限显隐从服务端下发的能力位读取，前端不做授权判断，见[后端架构](/architecture/backend)。
- 无限滚动用 `?after={postNumber}` 游标 + `useInfiniteQuery`，未读落点见[数据模型](/architecture/data-model)。
