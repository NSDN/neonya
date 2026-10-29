import type { DefaultTheme } from 'vitepress'

export function useNavigation(): DefaultTheme.NavItem[] {
  return [
    { text: '首页', link: '/' },
    { text: '架构', link: '/architecture/overview' },
    { text: '设计', link: '/design/post-entry-layout' },
    { text: '决策记录', link: '/decisions/' },
    { text: '维护手册', link: '/maintenance/' },
    { text: '文档开发说明', link: '/develop-description/' },
  ]
}
