# 2026-06-17 API 路径命名约定

## 背景

前端 `API_URLS` 常量中的路径写法不统一（有的带 `/`、有的不带），
实际运行时 base URL 末尾是否带 `/` 也不确定，
字符串拼接 URL 容易出错（`/api` + `register` → `/apiregister`，`/api` + `/register` → `//register`）。

## 讨论脉络

- `new URL()` 的行为：path 以 `/` 开头时视为绝对路径，会直接替换 base 的整个 pathname。
- 如果 base 末尾不带 `/`（如 `http://localhost:10127/api`），
  其最后一节 `/api` 被 URL 解析器视为文件而非目录，
  相对路径 `boards` 会替换掉 `/api`，得到 `http://localhost:10127/boards`。

## 决策

**API 路径使用相对路径（不带前导 `/`）**。

- `API_URLS` 中所有值为不带前导 `/` 的相对路径：`'register'`、`'boards'`、`'threads'`。
- `request.ts` 中用 `new URL(path, baseURL)` 构造完整 URL。
- base URL 约定以 `/` 结尾（如 `http://localhost:10127/api/`），
  在 `request.ts` 中做防御性处理（`base.endsWith('/') ? base : base + '/'`）。

## 影响

- `apps/web/src/shared/constants/api.ts` — 所有 API 路径为无前导 `/` 的相对路径。
- `apps/web/src/shared/services/request.ts` — 使用 `new URL()` 构造 URL。
