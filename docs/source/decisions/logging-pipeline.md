# 2026-10-31 日志管线与 panic 处理

## 背景

服务端日志同时有两种消费者：人（扫读错误）和采集器（解析字段）。M0 落地 HTTP 层时，chi 自带的两个中间件与这个目标冲突：

- `middleware.RealIP` 无条件把 `X-Forwarded-For` 最左值写进 `r.RemoteAddr`，上游没设代理头时客户端可以直接伪造 IP（已在上游标 Deprecated）；
- `middleware.Recoverer` 把美化后的堆栈写到 **stderr**，响应体是空的——只读 stdout 的采集器会整条丢掉 panic。

本次一并把「日志长什么样、panic 怎么兜、请求 ID 和客户端 IP 从哪来」定下来。

## 讨论脉络

- **彩色输出 vs 结构化 JSON**：彩色只服务人读；ANSI 转义会让 `jq` 直接解析失败。想靠 `isatty` 自动分流，但这里链路是 `task` → Task 内建 shell → 进程，加上重定向、`docker -t`、journald 捕获，判出来的结果常常与「人是否在看」相反。
- **zap / zerolog**：性能更好、字段更花哨，但要额外依赖与一套配置；论坛量级不需要。
- **stdout vs stderr**：stdout 是被采集的那条流。把堆栈写 stderr，等于在「只采 stdout」的部署里静默丢失。
- **Recoverer 的位置**：放最外层时 panic 会穿透 `requestLogging`，后半段代码被跳过，**这条请求就没有访问日志**，只剩堆栈。
- **请求 ID 的来源**：沿用客户端传来的 `X-Request-Id`，外部随手伪造一个就能让多条请求共用同一条线索，排障与将来的审计都不可信。

## 决策

1. **日志唯一格式**：`slog` JSON 到 stdout；不做彩色、不做美化输出。
2. **人读视图交给下游**：需要颜色时用 `task api:dev | jq -C`，渲染逻辑不进二进制。
3. **中间件顺序**：`requestID → clientIP → requestLogging → recovering`。值只能通过「传下去的 request 的 context」被下游读到，所以注入型中间件必须在前；panic 兜底紧贴 handler，这样 500 也会被访问日志记下来。
4. **panic 兜底自写**：`recovering` 写一行 slog JSON（含 `stack`），并返回 `application/problem+json` 的 500，而不是空 body。
5. **请求 ID 由服务端生成**：忽略入站同名头，并回写 `X-Request-Id` 响应头，用户可以把它报给运维。
6. **客户端 IP 只按信任跳数取**：`TRUSTED_PROXY_COUNT > 0` 才从 `X-Forwarded-For` 取第 N 跳，否则只认 TCP 来源。填小了=可伪造，填大了=拿不到 IP。

## 影响

- `api/internal/transport/http/`：`request_id.go`、`recovering.go`、`logging.go`、`router.go` 四处的分工；以后新增中间件要放在 `requestLogging` 之前才能被日志读到。
- 采集器只订阅 stdout 即可，访问日志与 panic 记录同流，用 `request_id` 串联。
- 前端遇到 500 时，把响应头里的 `X-Request-Id` 交给运维即可定位。
- `TRUSTED_PROXY_COUNT` 的取值见[部署与运行](/architecture/deployment) §3。
- 纯人读的一次性输出（将来的 `api admin` 子命令、`task db:status`）不受此约束，那些地方上色是净收益。
