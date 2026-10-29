# 喵玉殿新版论坛开发文档

本仓库同时承载代码与文档：

| 路径     | 内容                                   |
| -------- | -------------------------------------- |
| `api/`   | Go 后端（chi + Huma + sqlc + goose）   |
| `web/`   | Next.js 前端                           |
| `docs/`  | VitePress 文档站（本页面）             |
| `draft/` | 待整理的草稿，不对外链接               |

## 本地预览

```shell
pnpm install     # 仓库根目录
pnpm dev:docs    # http://localhost:10126/
```

参与文档写作请先读[文档开发说明](/develop-description/)。

## 文档导航

- [架构总览](/architecture/overview)：技术栈、目录结构、取舍原则（拆为后端 / 数据模型 / 前端 / 部署四页）
- [条目布局与交互](/design/post-entry-layout)：主题详情页的线框图、元素编号与术语表
- [决策记录](/decisions/)：一篇一个决策，含讨论脉络与被否决的方案
- [画面设计书](/frontend/page-design/sidebar)：边栏与登入画面
- [需求讨论记录](/prd/discussion-record)
- [归档](/archive/document_v0.0.x)：旧版需求文档与群内讨论

## 约定

- 架构与设计决策记在 `source/decisions/`，**歧义时先查这里**，再提新方案。
- 架构正文（`source/architecture/`）只写「是什么」，不写「为什么」。
- 中文文案术语见[条目布局与交互](/design/post-entry-layout) 第 0 节（主题 / 条目 / 首条 / 作者）。
