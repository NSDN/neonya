# 2026-06-20 Store 与 Composable 职责边界

## 背景

前端重构登入功能时，讨论 store 和 composable 各自应该承担什么职责。
`error.notify()` 是项目全局的报错通道（通过 Naive UI 的 `window.$message` 弹出用户可见的提示），
store 内是否可以使用它成为讨论的起点。

## 讨论脉络

初始直觉是"store 纯数据，composable 管副作用"。
但在实际代码中，`userStore.logout()` 需要同时清除 JWT 和用户信息，
涉及跨 store 协调；`jwtStore` 初始化时从 `storage.get()` 恢复 token 也可能失败。
如果禁止 store 调用 `notify()`，这些失败路径会丢失用户提示。

进一步查阅 Pinia 官方文档，
官方将 store action 定位为**业务逻辑中枢**：
API 调用、跨 store 协调、错误提示都属于 action 的正常职责。
社区并未划分"store 不能有副作用"这条线。

最终认识到 store 和 composable 的边界不在于"有没有副作用"，
而在于数据的**共享范围**。

## 决策

**Store = 全局单例状态 + 业务逻辑。Composable = 组件局部状态 + 胶水逻辑。**

具体规则：

1. **放在 Store 里**：状态本身（ref）、修改状态的同步 setter、跨 store 的业务逻辑（如 `logout()` 同时清理 jwt 和 user）、异步 API 调用（action 内可 `await`，可 `notify()`）。
2. **放在 Composable 里**：组件的局部响应式状态（如 `loginInfo` 表单、`formError`、`loading` 标记）、协调局部状态 + store + router 的编排逻辑（如 `executeLogin`、`handleGetBoards`）。
3. **API 函数**放在 `apis/`，不放在 store 或 composable 中，保持纯函数可复用。
4. **`error.notify()` 可在任意位置调用**——它是项目级的报错通道，store 内的失败同样需要告知用户。

```typescript
// Store — 全局状态 + 业务逻辑
export const useUserStore = defineStore(STORE_ID.USER, () => {
  const uid = ref('')
  const nickname = ref('游客')

  const setUserInfo = (info: UserInfo) => { /* 更新 ref */ }

  const logout = () => {
    // 跨 store 协调：同时清理 jwt 和 user
    const jwtStore = useJWTStore()
    jwtStore.setJWT(Option.none())
    uid.value = ''
    nickname.value = '游客'
  }
})

// Composable — 组件局部状态 + 胶水逻辑
export function useLogin() {
  const loginInfo = reactive<LoginInfo>({ username: '', password: '' })
  const formError = reactive<LoginInfoError>({ username: '', password: '' })

  const executeLogin = async () => {
    // 验证（局部状态）→ API → store 更新 → router（编排逻辑）
  }
}
```

## 影响

- Store 内可自由使用 `error.notify()`、API 调用、跨 store 引用，符合 Pinia 社区惯例。
- Composable 的存在意义是承载**不能全局共享**的状态（表单、loading 等），以及**特定 UI 场景的编排逻辑**。
- 新增功能模块时按此边界判断：数据是否全局唯一 → store；是否每个组件实例独立 → composable。
