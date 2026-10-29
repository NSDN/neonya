# 2026-10-30 API 版本与路径

## 背景

前后端同源部署（Nginx 按路径分流），需要确定 API 的路径前缀与版本策略，以及浏览器与 Server Component 各自怎么拿 base URL。

## 讨论脉络

- **不带版本前缀**（`/boards`）：URL 短；但破坏性变更时无法并存新旧接口。
- **带版本前缀 `/api/v1/*`**：路径可读、与静态资源区分开、将来可并存 v2。
- **版本放请求头**（`Accept: application/vnd.xxx`）：更"纯"的 REST；但 `curl` 调试不直观，缓存与访问日志也不友好。
- **浏览器用绝对地址**：可以跨域；但需要 CORS 与 Cookie 跨站处理，与「同源部署」的设计相悖。
- **Server Component 用公网地址**：实现简单；但取数会绕一圈公网，且需要额外的 DNS 与证书路径。

## 决策

1. API 前缀固定 **`/api/v1`**，路由集中在 `transport/http` 注册，Huma 导出的 OpenAPI 同样以该前缀为准。
2. **浏览器**：一律相对路径 `/api/v1/...`（`openapi-fetch` 的 `baseUrl = '/api/v1'`），同源、无 CORS。
3. **Server Component**：使用容器内地址 `API_BASE_URL`（开发 `http://localhost:8080`，生产 `http://api:8080`）。
4. 前端不手工拼 URL、不写前导斜杠判断——路径由生成 client 持有。
5. 破坏性变更时新增 `/api/v2`，旧版本保留一个发布周期再删。

## 影响

- `api/internal/transport/http`：所有路由挂在 `/api/v1` 路由组下。
- `web/src/lib/api/client.ts`：两个实例——浏览器（相对路径）与服务端（绝对地址），共用生成的类型。
- 旧 [API 路径命名约定](./api-url-naming) 中「路径不带前导 `/`、用 `new URL()` 拼接」的写法被本条取代。
- 统一域名下 `/api/*` 由 Nginx 转发，见[部署](/architecture/deployment)。
