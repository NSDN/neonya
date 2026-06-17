# 决策记录

开发过程中做出的架构与设计决策，记录讨论脉络和最终结论。

| 日期       | 决策                                                              | 简述                                         |
| ---------- | ----------------------------------------------------------------- | -------------------------------------------- |
| 2026-05-23 | [决策记录的文件组织形式](./decision-records-file-organization.md) | 一主题一文件，ADR 模板，AGENTS.md 引导查阅。 |
| 2026-05-23 | [SPA 前端路由设计](./spa-route-structure.md)                      | 同路由 + 条件渲染，Hash → History。          |
| 2026-05-24 | [前端响应式布局设计](./frontend-responsive-layout.md)             | 一套布局 + 48em 断点，移动端 BottomSheet。   |
| 2026-06-17 | [领域术语命名](./domain-term-naming.md)                           | Plate→Board、Topic→Thread，统一前后端命名。 |
| 2026-06-17 | [版块列表的数据源与默认选中](./board-list-data-source.md)         | 数据库为唯一数据源，删除前后端硬编码。       |
| 2026-06-17 | [HTTP 层换用原生 fetch](./http-native-fetch.md)                   | ofetch → 原生 fetch，零第三方 HTTP 依赖。   |
| 2026-06-17 | [API 路径命名约定](./api-url-naming.md)                           | 相对路径 + `new URL()`，统一路径写法。       |
| 2026-06-17 | [Guard 条件提取模式](./guard-condition-extraction.md)             | `shouldSkipFetch()` + `conditions.some()`。 |
| 2026-06-18 | [浮动按钮组设计](./floating-action-group.md)                      | 三按钮子组件 + 容器，BottomSheet 合并简化。  |
