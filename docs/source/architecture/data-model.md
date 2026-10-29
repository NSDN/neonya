# 数据模型

> 为什么这样建模见[决策记录](/decisions/)。本文是表、字段与索引的**唯一权威**。
> 界面上的呈现方式见[条目布局与交互](/design/post-entry-layout)。

## 1. 表清单

| 领域 | 表 |
| --- | --- |
| 账号 | `users`、`profiles`、`roles`、`permissions`、`user_roles` |
| 内容 | `boards`、`board_moderators`、`topics`、`posts`、`post_revisions`、`tags`、`topic_tags` |
| 阅读 | `topic_reads` |
| 互动 | `reactions`、`bookmarks`、`follows` |
| 运营 | `notifications`、`reports`、`audit_logs`、`uploads`、`settings` |

- 会话存 Redis，不建表；`river_*` 表由 River 自动迁移。
- 私信（`conversations` / `conversation_participants` / `conversation_messages`）**预留到 M3**，本次不建表。

## 2. 主题与条目

`topics` 只放元数据，**不存正文**；正文全在 `posts`，`post_number = 1` 的那条即主题正文（首条）。

```sql
topics
├─ id            uuid pk
├─ board_id      uuid
├─ author_id     uuid
├─ title         text
├─ kind          smallint   -- 0 = normal, 1 = comic（只有漫画区的展示形态不同）
├─ is_pinned     bool
├─ is_locked     bool
├─ post_count    int
├─ view_count    int
├─ last_post_at  timestamptz
├─ created_at / updated_at / deleted_at
```

漫画主题的额外信息放侧表，避免 `topics` 为普通主题长出一堆 NULL：

```sql
comic_topic_details
├─ topic_id       uuid pk → topics(id)
├─ original_title text      -- 原作名
├─ author_name    text      -- 原作作者
├─ chapter        text      -- 章节 / 话数
├─ group_name     text      -- 汉化组
└─ page_count     int
```

## 3. 引用模型（扁平流）

**回复 = 引用。没有嵌套，最大缩进恒为 1。**

| 场景 | `reply_to_post_id` | `quoted_excerpt` | 显示引用块 | 通知谁 |
| --- | --- | --- | --- | --- |
| 回复主题 | `NULL` | `NULL` | ❌ | 主题作者 |
| 引用主题正文（全文/划句） | 首条 `#1` 的 id | 片段 | ✅ | 主题作者 |
| 引用某条（全文/划句） | `X` | 片段 | ✅ | `X` 的作者 |
| 回复某条但删掉引用块 | `X` | `NULL` | ❌ | `X` 的作者 |

```sql
posts
├─ id               uuid pk
├─ topic_id         uuid
├─ post_number      int          -- 主题内序号，从 1 开始；锚点 #N 与分页游标
├─ author_id        uuid
├─ reply_to_post_id uuid null    -- 主引用目标，单边
├─ quoted_excerpt   text null    -- 被引文本快照，纯文本 ≤ 500 字
├─ quoted_is_full   bool
├─ quoted_offsets   int4range null  -- 划句引用位置（M2 用于原帖内高亮）
├─ content_md / content_html / content_plain / render_version / content_format
├─ created_at / edited_at / deleted_at
UNIQUE (topic_id, post_number)
```

规则：

- **存快照而不只存 ID**：原帖被编辑或删除后引用仍可读，否则引用会变死链。
- **MVP 只支持单引用**；多引用留到 M2，用 `post_quotes` 侧表扩展，可由 `reply_to_post_id` 回填。
- **被引用索引**：`SELECT post_number, author_id FROM posts WHERE reply_to_post_id = $1 AND deleted_at IS NULL ORDER BY post_number`；默认显示前 3 条，超过 8 条折叠。
- **通知去重**：同一 `(user, topic)` 15 分钟内只发一次（River job）。
- **分页**：`?after={postNumber}` 游标，服务端返回 30 条；没有 `?page=N`。

```sql
topic_reads                       -- 无限流的前提
├─ user_id               uuid
├─ topic_id              uuid
├─ last_read_post_number int      -- 「跳到未读」的落点
├─ notification_level    smallint  -- watching / tracking / normal / muted
└─ PRIMARY KEY (user_id, topic_id)
```

## 4. 内容格式

**Markdown 源文是唯一真相，HTML 只是渲染缓存。**

| `posts` 字段 | 用途 |
| --- | --- |
| `content_md` | Markdown 源文（可含白名单语义标签），**唯一真相** |
| `content_html` | 渲染缓存，可从 `content_md` 全量重建，允许为空 |
| `render_version` | 渲染器版本号；升级后按 `WHERE render_version < N` 精准重渲染 |
| `content_plain` | 抽取的纯文本，供 `pg_trgm` 搜索与 RSS 摘要 |
| `content_format` | `0` = markdown（默认），`1` = trusted_html（仅管理员公告） |

三条原则：

1. **能用扩展语法的用扩展语法**（spoiler、图片尺寸）——属于内容，必须入库；
2. **能渲染期算的一律不存**（@提及、链接预览卡片、emoji）——否则改名/失效后数据腐烂；
3. **加新语法优于开 HTML**——语法在白名单代码里，HTML 在用户输入里。

扩展语法清单：`:::spoiler`（折叠）、`![](url){width=300 align=center}`（图片尺寸/对齐）、`@name`（提及）、`:emoji:`（表情）、`attachment://uuid`（附件）、`$...$`（公式，后置可选）。

允许的**语义标签**白名单（无 `style`）：`<details> <summary> <sup> <sub> <kbd> <mark> <ruby> <rt> <video controls>`（仅白名单域 `src`）。
一律禁止：`style` / `class` / `on*` 属性，`<script> <iframe> <form> <object> <embed>`，`javascript:` / `data:` 协议。

按角色分级：

| 内容 | 作者 | 允许 |
| --- | --- | --- |
| 主题 / 回复 | 所有用户 | Markdown + 扩展语法 + 语义标签白名单 |
| 公告 / 版规 / 关于页 | 管理员 | 受限 HTML（仍禁 `script` / `iframe` / `style`） |
| 拖拽排版的富文本页 | 管理员 | 独立 `pages` 表，不混进 `posts` |

**渲染时机：写时渲染。** 发帖/编辑时在同一事务内算出 `content_html` / `content_plain` / `render_version`，读路径只读不写。

- 渲染器升级后跑 `api admin rerender --to-version N`，分批重算，不要全表刷；
- 兜底：读路径发现 `content_html` 为空时内存里即时渲染返回，**不落库**；
- 三个字段都是派生物，随时可从 `content_md` 重建，选错也不丢数据。

## 5. 主键与 URL

- 主键统一 **UUID v7**（时间有序，避免 B-tree 随机写，便于外部引用）；
- 公开 URL 用不可变 ID：`/t/{topicId}`、`/f/{boardSlug}`、`/u/{userId}`，条目锚点 `#post-{postNumber}`；
- `users.username` 不可变（登录与提及用）；`profiles.display_name` 可改，不进 URL；
- 不用 slug 做 URL 主键：标题一改外链全断。

## 6. 索引

| 索引 | 覆盖场景 |
| --- | --- |
| `topics(board_id, is_pinned DESC, last_post_at DESC)` | 版块主题列表 |
| `posts(topic_id, post_number)` | 扁平流与游标分页 |
| `posts(reply_to_post_id)` | 被引用索引 |
| `posts(author_id, created_at DESC)` | 个人主页 |
| `pg_trgm` GIN on `topics.title` | 标题模糊搜索 |

## 7. 通用约定

- **软删除**：统一 `deleted_at`，回收站查 `deleted_at IS NOT NULL`；外键 `ON DELETE RESTRICT`。
- **计数冗余**：`topics.post_count / last_post_at / view_count`，写入时更新 + Redis 暂存浏览数定时落库；读路径禁止 `COUNT(*)`。
- **反应唯一**：`reactions` 上 `UNIQUE (user_id, post_id)`——一条内容只有一个反应，换表情即覆盖。
- **编辑历史**：`post_revisions` 保存改前内容 + 操作者，只在编辑时写。
- **审计日志**：所有管理端写操作记录 who / what / before / after。
- **sqlc 约定**：`queries/*.sql` 按域分文件，生成的 querier 不直接暴露给 service，由 `persistence` 包一层窄接口。

## 8. 相关决策

- [扁平流 + 引用即回复](/decisions/flat-reply-model)
- [领域术语命名](/decisions/domain-term-naming)
- [决策记录索引](/decisions/)
