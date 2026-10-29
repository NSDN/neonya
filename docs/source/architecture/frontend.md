# 前端架构

> 选型理由见[决策记录](/decisions/)。条目级别的线框图与元素编号见[条目布局与交互](/design/post-entry-layout)。

## 1. 目录

```
web/src/
├─ app/                  # 只有路由、布局、metadata（SEO 在这里）
│  ├─ (forum)/{page, f/[slug], t/[id]}
│  ├─ (auth)/{login, register, reset}
│  ├─ (user)/u/[id]/  settings/
│  ├─ admin/             # 后台路由组，独立布局
│  └─ search/
├─ features/             # 域内 hook + 组件 + 数据请求
│  └─ {topic, post, board, auth, notification, moderation}/
├─ components/           # 跨域复用（TopicCard、Editor、Markdown、Emoji）
│  └─ ui/                # shadcn/ui 装进来的组件（源码进仓库）
└─ lib/
   ├─ api/schema.d.ts    # 由 openapi.yaml 生成，禁止手改
   └─ api/client.ts      # openapi-fetch 实例 + 鉴权包装
```

`app/` 只做路由与布局，业务逻辑下沉到 `features/`；对外请求统一走 `lib/api/client.ts`，不在组件里手写 `fetch`。

## 2. 渲染策略

| 页面 | 策略 |
| --- | --- |
| 首页 / 版块页 | SSR，短 TTL 缓存 |
| 主题详情 | SSR 首屏 + 客户端续读（无限滚动）；发帖后 `revalidateTag(topic:id)` |
| 个人主页 | SSR，短 TTL |
| 设置 / 通知 | 客户端渲染（TanStack Query） |
| 后台 | 客户端渲染 + 服务端权限校验兜底 |

- 无限滚动用 `?after={postNumber}` 游标续读，配合 `topic_reads.last_read_post_number` 落点。
- **不使用 Next 的 Server Action 承载业务逻辑**：所有写操作打 Go API，避免两处实现权限与校验。

## 3. 数据层

- **类型来源**：`web/src/lib/api/schema.d.ts` 由 `task gen` 从 `api/openapi.yaml` 生成，提交进仓库。
- **调用**：`openapi-fetch` 类型化 client，浏览器用相对路径 `/api/v1/...`（同源，无 CORS）；Server Component 走容器内地址 `http://api:8080`。
- **服务端取数**：Server Component 直接 `await`，取数逻辑不进入客户端 bundle。
- **客户端状态**：TanStack Query 管缓存与失效；不引入全局 store。

## 4. 编辑器

**唯一存储格式是 Markdown，编辑器只是视图。**

| 阶段 | 编辑器 | 说明 |
| --- | --- | --- |
| MVP | CodeMirror + 实时预览 | 便宜、可控 |
| M1 后 | WYSIWYG（优先 Milkdown） | ProseMirror 内核 + Markdown 双向序列化 |

硬约束：编辑器只提供「能可靠序列化成 Markdown」的功能集；做不到无损往返的功能不放进编辑器。写入字段见[数据模型](/architecture/data-model) §4 内容格式。

## 5. 组件

- **shadcn/ui**：`pnpm dlx shadcn@latest add <component>`，源码装进 `components/ui/`，可改可 diff，无大版本升级压力。
- 样式统一走 Tailwind v4。
- 条目组件（`PostEntry`）的元素编号（`U1` / `T2` / `A3` / `Q3` / `B1` …）与数据字段的对应关系见[条目布局与交互](/design/post-entry-layout)。

## 6. 与后端的边界

| 职责 | 归属 |
| --- | --- |
| 参数校验（服务端权威） | Go（Huma struct tag） |
| 表单即时校验（仅体验） | 前端（Zod），不作为安全边界 |
| 权限判定 | Go；前端只读 `/api/v1/me` 的能力位 |
| Markdown 渲染 | Go（goldmark + bluemonday）；**前端只插入渲染好的 HTML**，绝不二次解析 |
| SEO / metadata | Next（`app/**/page.tsx` 的 `metadata`） |

## 7. 相关决策

- [扁平流 + 引用即回复](/decisions/flat-reply-model)
- [决策记录索引](/decisions/)
