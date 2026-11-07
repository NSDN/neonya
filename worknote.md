# Work Note

## 2026-11-08

现状：

- M0 后端骨架完成：`/api/v1/health` 连库探活已实测（200；停库 503），`task check` 全过。

下一步：

- `web/` 骨架，按 W1–W7 单元推进，切分见 `docs/source/architecture/overview.md`「阶段路线」。
  - 只做 W1（Next 脚手架，无 API 调用），做完停下等审阅。

## 2026-10-30

- 文档站迁移完成（VitePress、侧边栏、导航），`draft/` 已删除。
- 架构 5 页 + 条目布局线框图定稿。
- 决策记录：14 篇一题一篇。10 篇补了「2026-10-30 修订」节，消掉与现架构的冲突——4 篇标失效、4 篇部分更新、2 篇新增（API 版本与路径、前端状态分层）。
- 移动端发帖入口定案：浮动按钮组（FAB），不设常驻输入条。
- 提交信息规范进 `AGENTS.md`，并加了 `.pi/prompts/commit-message.md`（`/commit-message`）。
- 代码尚未开始：`api/`、`web/` 都还不存在。
- 下一步：S1 = M0 骨架，见 `docs/source/architecture/overview.md`「阶段路线」。
