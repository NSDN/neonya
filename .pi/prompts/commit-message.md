---
description: 查看改动并生成 Conventional Commits 提交信息（不执行提交）
argument-hint: '[scope]'
---

1. 运行 `git status` 与 `git diff --staged --stat`；  
   如果没有暂存内容，改用 `git diff --stat` 并说明。
2. 生成**一个**英文 commit message，符合 Conventional Commits：
   - type 取 feat / fix / docs / refactor / perf / test / build / ci / chore
   - scope 用 ${1:-从改动内容推断}，可省略
   - subject ≤ 72 字符，祈使句，结尾不加句号
   - 需要时补 body：3~6 条 bullet，写「改了什么」，不写空话
3. 另给 1~2 个备选 subject，一句话说明取舍。
4. 仓库尚无 commit 时，按「首次基线提交」处理，并在输出里点明这一点。
5. **不要执行 `git commit`**，只输出可直接粘贴的内容。
