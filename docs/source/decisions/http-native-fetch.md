# 2026-06-17 HTTP 层换用原生 fetch

## 背景

前端 HTTP 层原本使用 ofetch 第三方库。
项目仅用到了 `ofetch.create()` 的 `baseURL`、`headers`、`timeout` 和 `onResponse` 拦截器，
功能与原生 `fetch` 高度重合，引入第三方依赖的必要性值得讨论。

## 讨论脉络

- **方案 A：保留 ofetch** — 已配置完成，无需改动；但增加了包体积和维护依赖。
- **方案 B：换回 axios** — 功能全面；但包体积更大，API 风格偏重。
- **方案 C：换用原生 fetch** — Node 18+ 与所有现代浏览器内置，零额外依赖；
  `optionizeDeep` 已在项目内的 `shared/utils/useful.ts` 实现，
  超时用 `AbortSignal.timeout()` 替代，
  URL 拼接用 `new URL()` 替代字符串拼接。

## 决策

**换用原生 fetch**。

- 删除 `apps/web/src/shared/services/ofetch.ts`。
- `request.ts` 重写：`fetch()` + `AbortSignal.timeout()` + `optionizeDeep`。
- URL 构造使用 `new URL(path, baseURL)` 替代字符串拼接，处理 base 与 path 的边界情况。
- 移除 `ofetch` 依赖（减少 4 个包）。

## 影响

- `package.json` — 不再依赖 `ofetch`。
- `apps/web/src/shared/services/request.ts` — 完全重写。
- `apps/web/src/shared/services/ofetch.ts` — 删除。

## 参见

- [API 路径命名约定](./api-url-naming.md)

---

## 2026-10-30 修订：已被取代

新版改用 **`openapi-fetch`**：它本身也是原生 `fetch` 的薄封装，但请求/响应类型由 Huma 导出的 OpenAPI 生成，见[契约生成](/architecture/backend)。

| 原决策 | 现在 |
| --- | --- |
| 手写 `shared/services/request.ts` | `web/src/lib/api/client.ts` 包一层生成的 client（浏览器与服务端各一个实例） |
| 手写 URL 拼接 + `new URL()` | client 的 `baseUrl` + 生成的路径常量，见 [API 版本与路径](./api-versioning) |
| 不引入 ofetch / axios | **结论不变**：仍然不引入任何第三方 HTTP 客户端库 |

保留的是「零重量级 HTTP 依赖」这个结论，失效的只是「手写 request 层」这个实现方式。
