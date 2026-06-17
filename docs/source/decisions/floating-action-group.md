# 2026-06-18 浮动按钮组设计

## 背景

`DefaultLayout` 右下角在移动端模式下有一个"版块"按钮，触发 BottomSheet 版块选择面板。  
需要扩展为按钮组，容纳三个操作：回到顶层、新建帖文、打开版块菜单。  
按钮组需保留可扩展性，后续可追加新按钮。

## 讨论脉络

### 按钮组的职责边界

- **方案 A：容器只做定位，按钮通过 slot/prop 传入** — 灵活但外部需要知道按钮组内部结构。
- **方案 B：容器只暴露自己，内部全自理** — 自闭环，DefaultLayout 只需 `<FloatingActionGroup />` 一行。

决定方案 B。`FloatingActionGroup` 内部编排三个子组件，各自自治。

### 子组件粒度

- `ScrollToTopButton`：已有独立组件，移入按钮组目录，去除自定位（改由容器统一布局）。
- `NewPostButton`：从内联按钮提取为独立子组件，遵循一致性。
- `BoardMenuButton`：将原 BottomSheet 合并入按钮组件，版块按钮与面板为一个自治整体。

### BottomSheet 的状态管理三层简化

1. `show` prop + `update:show` emit（外部驱动）→ 内部 `show` ref（自管理）。
2. `show` + `visible` 双变量（visible 作为"动画 latch"）→ 单 `show` ref，`<Transition>` 自行管理退场 DOM 保留。
3. `watch` 监听 `show` → `openSheet()` 函数同步赋值。

### 按钮形状

- **圆形**（ScrollToTopButton 原有风格，`border-radius: 50%`）：FAB 辨识度高，但容纳文字空间有限。
- **圆角矩形**（`border-radius: 0.5rem`）：三个按钮尺寸统一，图标居中。

选择圆角矩形，按钮统一为 `3rem × 3rem` 正方形。

### 模板处理器风格

组件内存在 `closeSheet`、`handleSelectBoard` 等函数调用式处理器。  
`@click="show = true"` 内联赋值打破了函数调用风格的一致性。

决定提取 `openSheet()` 函数，模板中统一函数调用。

### 底部面板的标题栏

BottomSheet 原有标题栏（"版块" 标题 + "✕" 关闭按钮）。

决定删除。遮罩层点击已提供关闭路径，面板内版块网格直接展示，减少视觉层级。

## 决策

1. **按钮组设计**：同一目录下拆为三个子组件 + 一个容器。  
   容器仅做 fixed 定位 + flex column 布局，子组件各自闭环。
2. **BottomSheet 状态**：单一 `show` ref，`openSheet()` / `closeSheet()` 对称操作。
3. **按钮形状**：圆角矩形 `3rem × 3rem`，图标居中。
4. **ScrollToTopButton**：去除自定位，scroll 阈值可见性逻辑不变。
5. **面板交互**：删除标题栏，版块网格直接呈现，遮罩点击关闭。
6. **模板处理器**：全部使用函数调用，不内联赋值。

## 影响

### 目录变更

```text
layouts/components/
  Header.vue
  UserSimpleInfo.vue
  FloatingActionGroup/           ← 新建
    FloatingActionGroup.vue      ← 按钮组容器（fixed 定位 + flex column）
    ScrollToTopButton.vue        ← 从 ../ 移入，去自定位，改圆角矩形
    NewPostButton.vue            ← 新建，➕ 按钮 + TODO message
    BoardMenuButton.vue          ← BottomSheet.vue 合并入此组件

layouts/DefaultLayout/
  DefaultLayout.vue             ← 删 ScrollToTopButton / board-trigger / BottomSheet，
                                   改为 <FloatingActionGroup />
  BottomSheet.vue               ← 删除（合并入 BoardMenuButton）
```

### DefaultLayout 侧

```vue
<FloatingActionGroup />
```

无 props、无 emits、无 slot，完全自闭环。

### 按钮组的"回到顶层"与"新建帖文"

- 回到顶层：`window.$message.warning('TODO')`，等待实装。
- 新建帖文：`window.$message.warning('TODO')`，等待实装。
- 版块菜单（仅移动端）：打开底部版块选择面板。
