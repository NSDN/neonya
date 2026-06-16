# Todolist

1. 修改 BottomSheet 中的版块列表。
2. 研究一下请求帖子列表的路由是否有更佳写法（目前为 `thread?boardid=x`）。
3. 处理前端的 build 里的类型检查报错。
4. 补全登录流程 — apps/web/src/features/authorization/composables/useLogin.ts 中  
   获取 token 和用户信息的代码被注释掉了，这是其他功能的基础
5. 实现后端 GET /api/threads/:id — 目前无法查看单个帖子内容
6. 实现帖子详情页 — 前端路由 /threads/:id + 组件，引用最多的 TODO
7. 实现前端发帖 UI — 后端 API 已就绪，只差前端表单
