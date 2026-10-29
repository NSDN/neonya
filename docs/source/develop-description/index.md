# 文档开发说明

本文档站框架为 [VitePress](https://vitepress.vuejs.org/)，与代码同仓：`docs/`。

## 本地启动

```shell
# 仓库根目录执行（pnpm workspace）
pnpm install
pnpm dev:docs      # http://localhost:10126/
```

其他命令：

```shell
pnpm build:docs    # 构建到 docs/.vitepress/dist
pnpm preview:docs  # 预览构建产物
```

## 目录结构

```
docs/
├─ .vitepress/
│  ├─ config.ts          # 站点配置（标题、导航、侧边栏）
│  ├─ composition/       # use-sidebar / use-navigation 等
│  └─ theme/             # 自定义主题与图标
├─ types/                # 主题类型声明
└─ source/               # 所有页面（srcDir）
   ├─ index.md
   ├─ decisions/         # 决策记录（一篇一个决策）
   ├─ design/            # 界面与交互设计（含线框图）
   ├─ frontend/          # 画面设计书
   ├─ prd/               # 需求讨论记录
   ├─ reference/         # 参考资料
   ├─ archive/           # 历史归档（不再修改）
   └─ develop-description/
```

## 新增页面的流程

1. 在 `docs/source/` 下对应目录新建 `.md`（文件名 kebab-case）；
2. 在 `docs/.vitepress/composition/use-sidebar.ts` 里加条目（迁移期手写；等文档结构稳定后再换成 `vitepress-sidebar` 自动生成）；
3. **新增顶级分区时，导航与首页侧边栏要同时改**：`use-navigation.ts` 加一项，`useHomePageSidebar()` 加一组。只改一处就会出现「侧边栏有、导航没有」（或相反）；
4. `pnpm dev:docs` 本地确认无死链。

### 新增一条决策记录

1. 复制 `docs/source/decisions/_template.md`，命名为 `kebab-case-主题.md`；
2. 按模板写「背景 → 讨论脉络 → 决策 → 影响」，**讨论脉络要写被否决的方案和否决原因**；
3. 在 `decisions/index.md` 的表格登记一行。

## 写作约定

- 中文句子以「。」结尾；行过长时在句末换行，需要强制换行时用行尾两个空格。
- 术语统一：**主题（topic）/ 条目（post）/ 首条 / 作者**；禁用「楼、楼层、主楼、楼主」。
- 代码块标注语言；表格中避免超宽行。
- 格式由根目录 `.prettierrc` 约束，提交前跑 `pnpm format`。

## 提交

```shell
git add .
git commit -m 'docs: xxx'
git push
```

文档与代码同仓，因此**改行为的 PR 应同时更新受影响的页面与决策记录**。
