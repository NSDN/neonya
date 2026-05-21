import type { DefaultTheme } from 'vitepress'

export function useNavigation(): DefaultTheme.NavItem[] {
  return [
    { text: '首页', link: '/' },
    { text: '前端文档', link: '/frontend/overview' },
    { text: '后端文档', link: '/backend/overview' },
    { text: '文档开发说明', link: '/develop-description/' },
  ]
}
