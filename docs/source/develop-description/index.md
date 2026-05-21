# 文档开发说明

本文档框架为 [VitePress](https://vitepress.vuejs.org/)

## 仓库

<https://github.com/NSDN/nya-develop-document>

## 本地启动开发文档

```shell
# 1. 下载代码
git clone https://github.com/NSDN/nya-develop-document.git

# 2. 安装依赖
npm install

# 3.启动服务
npm run docs:dev

# 4.在浏览器中打开 URL
# http://localhost:10126/
```

## 更新开发文档流程

### 1. 在对应的目录内修改文件

```shell
# 进入目标目录
cd docs/source/prd/discussion-record
# 创建新文件
echo '# Sample' > sample.md
```

### 2. 修改 sidebar

```typescript
import type { DefaultTheme } from 'vitepress'

export default function useSidebar(): DefaultTheme.Sidebar {
  return {
    // 如果需要创建新边栏分支，请在下面创建方法后在此处按分类返回出去
    // ...: ES6 解构语法，此处为平铺边栏分支方法中 return 的返回值
    ...useHomePageSidebar(),

    // 产品需求文档
    ...useRequirementsDocumentSidebar(),
    ...useDiscussionRecordSidebar(),
  }
}

function useHomePageSidebar(): DefaultTheme.Sidebar {
  // ...
}

function useRequirementsDocumentSidebar(): DefaultTheme.Sidebar {
  // ...
}

function useDiscussionRecordSidebar(): DefaultTheme.Sidebar {
  const ROOT_PATH: string = '/prd/discussion-record'

  return {
    [ROOT_PATH]: [
      {
        text: '',
        items: [
          { text: '回到首页', link: '/' },
          { text: '需求讨论记录目录', link: `${ROOT_PATH}/` },
        ],
      },

      {
        text: '需求讨论记录',
        items: [
          {
            text: '2022.07.15 需求讨论（中午）',
            link: `${ROOT_PATH}/discussion_07_15.md`,
          },

          // 在此处添加新建的文章
          {
            text: 'Sample',
            link: `${ROOT_PATH}/sample.md`,
          },
        ],
      },
    ],
  }
}
```

### 3. 提交修改

```shell
# 1.
git add .
# 2.
git commit --message 'feat: update docs'
# 3.
git push origin <your branch>

# 4. 在 github 上创建 pull request 并联系管理员合并（或向管理员索取合并权限）。
```
