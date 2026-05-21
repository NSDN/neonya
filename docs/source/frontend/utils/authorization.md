# 认证授权

## hashSecreate

使用 crypto-js 对明文密码进行预处理。

- 类型

  ```TypeScript
  function hashSecreate(secreate: string): string
  ```

- 示例

  ```TypeScript
  const password = 'password'
  
  const hashed = hashSecreate(password)
  
  console.log(hashed) // 5e884...
  ```
