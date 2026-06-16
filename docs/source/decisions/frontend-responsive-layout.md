# 2026-05-24 前端响应式布局设计

## 背景

当前布局为桌面端单一设计（Header + Sidebar + Content + Footer，固定 100vh），  
无移动端适配。需要重新设计一套响应式布局，同时满足桌面和手机两个端的使用场景。

## 讨论脉络

### 布局构成

逐个审视了 CONTENT 以外各部分的必要性：

- **Header**（Logo + 用户信息）：核心功能，**必须保留**。  
  但高度从 5rem 缩减到 3rem，内部元素同步缩小（padding 1rem→0.5rem、头像 3.5rem→2rem）。
- **Sidebar**（版块导航）：核心功能，**必须保留**。但宽度从 24rem 缩减到 16rem。
- **Footer**：当前为占位符，无实际内容。**决定删除**。
- **背景图**：东方论坛的氛围需求，**保留**。

### 响应式策略

- **方案 A：一套响应式布局**（一个 DefaultLayout + CSS 媒体查询）→ 内容一致、维护成本低。
- **方案 B：电脑/手机各自布局**（两套 Layout 组件）→ 灵活但重复代码多，内容相同无必要复制。

### Header 的归属

`Header.vue` 放在 `DefaultLayout/` 下但被 `NoneSidebarLayout` 也引用，语义不对。

- **方案 A：共享零件**（Header 挪到 `layouts/components/`，各 layout 按需引入）→ 灵活、语义清晰。
- **方案 B：全局固定**（Header 提入 App.vue 层）→ 控制粒度粗，登录页不需要完整 header。

### 手机端版块导航

- **方案 A：底部汉堡** → 拇指友好度中等。
- **方案 B：底部 tab 导航** → 拇指友好度高，但版块数量多时放不下。
- **方案 C：底部 sheet**（浮动按钮触发 bottom sheet 面板选择版块）→ 拇指友好，版块列表可滚动。

### 返回按钮的取舍

- **桌面端**：浏览器自带返回按钮，无需额外 UI 占用空间。
- **手机端**：系统滑动返回 + 浏览器后退 + Header Logo 已提供三条返回路径，无需额外按钮。

决定**全端删除** `BackButton` 组件。

### 回到顶部按钮

长列表滚动场景（版块帖子列表）需要回到顶部按钮。  
决定**新增** `ScrollToTopButton`，放在 `layouts/components/` 下，  
作为 DefaultLayout content 区域右下角的浮动按钮，两端共用。

### 断点单位

像素断点（768px）不考虑用户缩放和 DPI 设置。改用 `em` 单位，值从布局尺寸推导：

- Sidebar 16rem + 内容区至少 32rem = 48rem
- 媒体查询中 `em` 基于浏览器默认字体大小，随缩放自然变化

约定断点为 **`48em`**。

## 决策

1. **布局结构**：Header(3rem) + Sidebar(16rem, 可折叠) + Content(flex:1) + 无 Footer。
2. **响应式**：一套 `DefaultLayout`，`48em` 断点区分桌面/手机。  
   桌面端 Sidebar 常驻 + 可折叠，手机端 Sidebar 不挂载，改为浮动按钮 + BottomSheet。
3. **目录结构**：提取共用 Header 到 `layouts/components/`，Sidebar 相关零件留在 `DefaultLayout/` 下。
4. **Header**：作为共享零件（方案 A），各 layout 按需引入。
5. **手机端导航**：底部浮动按钮 + BottomSheet（方案 C）。
6. **返回按钮**：全端删除。
7. **回到顶部**：新增 `ScrollToTopButton`，右下浮动，两端共用。
8. **断点**：`48em`，源码写 `@media (--mobile)`，Vite 插件构建时替换为 `@media (max-width: 48em)`。  
   断点值定义在 `vite.config.ts` 的 `MOBILE_BREAKPOINT` 常量。

### 断点实现细节

CSS 自定义属性不能在 `@media` 条件使用。  
尝试了 `postcss-custom-media` 插件，但 v12 版本只处理单文件内的 `@custom-media` 声明，  
Vue SFC scoped style 各自独立处理无法跨文件共享声明。
最终在 `vite.config.ts` 中写 `enforce: 'pre'` transform 插件，对 `.vue` / `.css` 文件做字符串替换。

### `useBoards()` 重复请求防御

桌面端 Sidebar 和 BottomSheet 由 CSS 控制显隐，两个组件都始终挂载，各自在 `onMounted` 调用 `handleGetBoards()`。

为避免重复请求：

- Boards store 新增 `loading: ref<boolean>` 标记。
- `handleGetBoards()` 首行加 `if (loading) return; if (boards.length > 0) return`。
- 模板顺序保证 Sidebar 先于 BottomSheet 挂载，Sidebar 设置 `loading = true` 后 BottomSheet 调用直接返回。

## 影响

### 目录变更

```text
layouts/components/  ← 新建
  Header.vue         ← 从 DefaultLayout/ 移入
  UserSimpleInfo.vue ← 从 DefaultLayout/ 移入
  ScrollToTopButton.vue ← 新增

layouts/DefaultLayout/
  BottomSheet.vue    ← 新增
  Sidebar.vue        ← 宽度 24rem → 16rem
  SidebarController.vue ← 桌面端始终可见
  DefaultLayout.vue  ← 删 footer, 断点控制, 集成新组件

layouts/NoneSidebarLayout.vue ← 删 footer, 更新 import

features/board/stores/boards.ts    ← 新增 loading 标记
features/board/composables/useBoards.ts ← 新增防重复守卫

vite.config.ts ← 新增 customMediaPlugin（断点替换）

删除:
  shared/components/buttons/BackButton.vue
  layouts/DefaultLayout/Header.vue (原始)
  layouts/DefaultLayout/UserSimpleInfo.vue (原始)
```

### 路由 / App 层

- `App.vue`、`useLayout()`、路由配置：**不动**。
- `BoardContent.vue`、`BoardOnSidebar.vue`：**不动**。
- `Login.vue`、`Register.vue`：删除 `BackButton` 引用。

### `useBoards()` composable

Sidebar 和 BottomSheet 各自调用 `useBoards()`。  
两者由 CSS 控制显隐，始终同时挂载，各自在 `onMounted` 调用 `handleGetBoards()`。  
通过 store 层 `loading` 标记 + composable 防重复守卫保证只请求一次 API。
