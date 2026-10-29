import type { DefaultTheme } from 'vitepress'

export default function useSidebar(): DefaultTheme.Sidebar {
  return {
    ...useHomePageSidebar(),
    ...useArchitectureSidebar(),
    ...usePrdSidebar(),
    ...useDecisionsSidebar(),
    ...useDesignSidebar(),
    ...useFrontendSidebar(),
    ...useArchiveSidebar(),
    ...useDevelopDescriptionSidebar(),
    ...useMaintenanceSidebar(),
  }
}

function useHomePageSidebar(): DefaultTheme.Sidebar {
  return {
    '/': [
      {
        text: '架构',
        collapsed: false,
        items: [
          { text: '总览', link: '/architecture/overview' },
          { text: '后端', link: '/architecture/backend' },
          { text: '数据模型', link: '/architecture/data-model' },
          { text: '前端', link: '/architecture/frontend' },
          { text: '部署与运行', link: '/architecture/deployment' },
        ],
      },
      {
        text: '设计',
        collapsed: false,
        items: [
          { text: '条目布局与交互', link: '/design/post-entry-layout' },
          { text: '图片管理', link: '/design/picture' },
        ],
      },
      {
        text: '决策记录',
        collapsed: false,
        items: [{ text: '决策记录索引', link: '/decisions/' }],
      },
      {
        text: '产品需求文档（PRD）',
        collapsed: false,
        items: [{ text: '需求讨论记录', link: '/prd/discussion-record' }],
      },
      {
        text: '画面设计书',
        collapsed: false,
        items: [
          { text: '边栏', link: '/frontend/page-design/sidebar' },
          { text: '登入画面', link: '/frontend/page-design/login' },
        ],
      },
      {
        text: '维护',
        collapsed: false,
        items: [{ text: '维护手册', link: '/maintenance/' }],
      },
      {
        text: '文档',
        collapsed: false,
        items: [{ text: '文档开发说明', link: '/develop-description/' }],
      },
      {
        text: '资料',
        collapsed: false,
        items: [{ text: '参考资料', link: '/reference/reference' }],
      },
      {
        text: '归档',
        collapsed: true,
        items: [
          { text: '需求文档 v0.0.x', link: '/archive/document_v0.0.x' },
          { text: '2022.07.15 需求讨论（中午）', link: '/archive/discussion_07_15' },
          { text: '20221230 群内讨论', link: '/archive/discussion_20221230' },
          {
            text: '20240822 群内讨论 漫画版块',
            link: '/archive/discussion-20240822-commic-plate',
          },
        ],
      },
    ],
  }
}

function useArchitectureSidebar(): DefaultTheme.Sidebar {
  const pages = [
    { text: '总览', link: '/architecture/overview' },
    { text: '后端', link: '/architecture/backend' },
    { text: '数据模型', link: '/architecture/data-model' },
    { text: '前端', link: '/architecture/frontend' },
    { text: '部署与运行', link: '/architecture/deployment' },
  ]

  return {
    '/architecture/': [
      { text: '', items: [{ text: '回到首页', link: '/' }] },
      { text: '架构', collapsed: false, items: pages },
    ],
  }
}

function usePrdSidebar(): DefaultTheme.Sidebar {
  return {
    '/prd/': [
      {
        text: '',
        items: [
          { text: '回到首页', link: '/' },
          { text: '需求讨论记录', link: '/prd/discussion-record' },
        ],
      },
    ],
  }
}

function useDecisionsSidebar(): DefaultTheme.Sidebar {
  return {
    '/decisions/': [
      {
        text: '',
        items: [
          { text: '回到首页', link: '/' },
          { text: '决策记录索引', link: '/decisions/' },
        ],
      },
    ],
  }
}

function useDesignSidebar(): DefaultTheme.Sidebar {
  return {
    '/design/': [
      {
        text: '',
        items: [
          { text: '回到首页', link: '/' },
          { text: '条目布局与交互', link: '/design/post-entry-layout' },
          { text: '图片管理', link: '/design/picture' },
        ],
      },
    ],
  }
}

function useFrontendSidebar(): DefaultTheme.Sidebar {
  return {
    '/frontend/': [
      {
        text: '',
        items: [
          { text: '回到首页', link: '/' },
          { text: '边栏', link: '/frontend/page-design/sidebar' },
          { text: '登入画面', link: '/frontend/page-design/login' },
        ],
      },
    ],
  }
}

function useArchiveSidebar(): DefaultTheme.Sidebar {
  return {
    '/archive/': [
      {
        text: '',
        items: [
          { text: '回到首页', link: '/' },
          { text: '需求文档 v0.0.x', link: '/archive/document_v0.0.x' },
          { text: '2022.07.15 需求讨论（中午）', link: '/archive/discussion_07_15' },
          { text: '20221230 群内讨论', link: '/archive/discussion_20221230' },
          {
            text: '20240822 群内讨论 漫画版块',
            link: '/archive/discussion-20240822-commic-plate',
          },
        ],
      },
    ],
  }
}

function useDevelopDescriptionSidebar(): DefaultTheme.Sidebar {
  return {
    '/develop-description/': [
      {
        text: '',
        items: [
          { text: '回到首页', link: '/' },
          { text: '文档开发说明', link: '/develop-description/' },
        ],
      },
    ],
  }
}

function useMaintenanceSidebar(): DefaultTheme.Sidebar {
  return {
    '/maintenance/': [
      {
        text: '',
        items: [
          { text: '回到首页', link: '/' },
          { text: '维护手册', link: '/maintenance/' },
        ],
      },
      {
        text: '流程',
        collapsed: false,
        items: [
          {
            text: '依赖与工具更新流程',
            link: '/maintenance/dependency-update',
          },
        ],
      },
    ],
  }
}
