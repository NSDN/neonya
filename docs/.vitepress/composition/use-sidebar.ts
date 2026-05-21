import type { DefaultTheme } from 'vitepress'

export default function useSidebar(): DefaultTheme.Sidebar {
  return {
    ...useHomePageSidebar(),

    // 前端文档
    ...useFrontendSidebar(),

    // 后端文档
    ...useBackendSidebar(),
  }
}

function useHomePageSidebar(): DefaultTheme.Sidebar {
  const REFERENCE_ROOT = '/reference/'
  const ARCHIVE_PATH = '/archive/'

  return {
    '/': [
      {
        text: '产品需求文档（PRD）',
        collapsed: false,
        items: [{ text: '需求讨论记录', link: '/prd/discussion-record' }],
      },
      {
        text: '方式设计',
        collapsed: false,
        items: [{ text: '图片管理', link: '/design/picture' }],
      },
      {
        text: '资料',
        items: [{ text: '参考资料', link: `${REFERENCE_ROOT}reference` }],
      },
      {
        text: '归档',
        collapsed: true,
        items: [
          { text: '需求文档 v0.0.4', link: `${ARCHIVE_PATH}/document_v0.0.x` },
          {
            text: '2022.07.15 需求讨论（中午）',
            link: `${ARCHIVE_PATH}discussion_07_15.md`,
          },
          {
            text: '20221230 群内讨论',
            link: `${ARCHIVE_PATH}discussion_20221230.md`,
          },
          {
            text: '20240822 群内讨论 漫画版块',
            link: `${ARCHIVE_PATH}discussion-20240822-commic-plate.md`,
          },
        ],
      },
    ],
  }
}

function useFrontendSidebar(): DefaultTheme.Sidebar {
  const ROOT_PATH = '/frontend/'
  const PAGE_DESIGN_PATH = `${ROOT_PATH}page-design/`
  const COMMON_COMPONENTS_PATH = `${ROOT_PATH}common-components/`
  const UTILS_PATH = `${ROOT_PATH}utils/`

  return {
    [ROOT_PATH]: [
      {
        text: '',
        items: [
          { text: '前端 - 回到首页', link: '/' },
          { text: '概述', link: `${ROOT_PATH}overview` },
          { text: '画面基本布局', link: `${ROOT_PATH}main-layout` },
        ],
      },

      {
        text: '画面设计书',
        collapsed: false,
        items: [
          { text: '边栏', link: `${PAGE_DESIGN_PATH}sidebar` },
          { text: '登入画面', link: `${PAGE_DESIGN_PATH}login` },
        ],
      },

      {
        text: '共通组件',
        collapsed: false,
        items: [{ text: 'TODO', link: `${COMMON_COMPONENTS_PATH}todo` }],
      },

      {
        text: '工具函数',
        collapsed: false,
        items: [
          {
            text: '认证授权',
            items: [
              {
                text: 'hashSecreate',
                link: `${UTILS_PATH}authorization#hashsecreate`,
              },
            ],
          },
        ],
      },
    ],
  }
}

function useBackendSidebar(): DefaultTheme.Sidebar {
  const ROOT_PATH: string = '/backend/'
  const API_DESIGN_PATH: string = `${ROOT_PATH}api-design/`

  return {
    [ROOT_PATH]: [
      {
        text: '',
        items: [
          { text: '后端 - 回到首页', link: '/' },
          { text: '概述', link: `${ROOT_PATH}overview` },
        ],
      },

      {
        text: 'API 设计书',
        collapsed: false,
        items: [
          {
            text: 'HTTP 连接测试 API',
            link: `${API_DESIGN_PATH}ping-test`,
          },

          {
            text: '认证授权',
            items: [
              {
                text: '登入 API',
                link: `${API_DESIGN_PATH}authorization/login`,
              },
            ],
          },

          {
            text: '论坛版块列表获取 API',
            link: `${API_DESIGN_PATH}get-plates`,
          },
          {
            text: '喵玉汉化馆帖文列表获取 API',
            link: `${API_DESIGN_PATH}get-localization-plate`,
          },
          {
            text: '创建漫画主题 API（草案）',
            link: `${API_DESIGN_PATH}create-commic-topic`,
          },
        ],
      },
    ],
  }
}
