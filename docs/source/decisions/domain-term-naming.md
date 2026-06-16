# 2026-06-17 领域术语命名

## 背景

代码库中版块和帖子的英文对应词需要确定。
当时的命名（`plate`、`topic`/`post`）既非论坛领域的惯用词，前后端包名也存在不一致：
前端用 `topic`，Go 后端包名却用 `post`。

## 讨论脉络

### 版块

- `Plate`（当前）→ 英文论坛软件中几乎不用，语义偏离。
- `Section` → 过于泛用，任何场景都可能出现，歧义多。
- `Board` → 有明确 BBS 历史渊源（Bulletin Board System），日本 2ch／5ch 用"板"称呼各分类，中文"版块"的"版"直接对应 Board。

选择 **Board**：更精确、更有论坛语感、与中文概念对应最自然。

### 帖子

- `Topic`（当前）→ Discourse 的术语，强调讨论"主题"的内容属性。
- `Thread` → 传统论坛（phpBB、vBulletin、XenForo）的通用术语，强调首帖 + 回帖串联的"线程"结构。

喵玉殿本质是传统 BBS 论坛，选择与 `Board` 搭配经典的 **Thread**：  
Board + Thread = 版面 + 帖子线程，是论坛领域的经典组合，语义连贯。

同时将后端包名 `post` 统一为 `thread`，消除包名与领域实体名不一致的问题。

### 文章 / 漫画

- `Article`、`Comic` —— 作为 Thread 的子类型，语义清晰，保持不变。

## 决策

| 中文 | 旧名    | 新名       |
| ---- | ------- | ---------- |
| 版块 | Plate   | **Board**  |
| 帖子 | Topic   | **Thread** |
| 文章 | Article | 不变       |
| 漫画 | Comic   | 不变       |

## 影响

- **后端**：`internal/plate/` → `internal/board/`，`internal/post/` → `internal/thread/`。  
  API 路由 `/plates` → `/boards`，`/topics` → `/threads`。  
  数据库表 `plates` → `boards`，`topics` → `threads`；列 `plate_id` → `board_id`，`topic_type` → `thread_type`。
- **前端**：`features/plate/` → `features/board/`，`features/topic/` → `features/thread/`。  
  所有 Store / API / Composable / 常量 / 路由 / 组件 / CSS class 同步重命名。
- **文档**：AGENTS.md、决策记录、前端概述文档同步更新。
