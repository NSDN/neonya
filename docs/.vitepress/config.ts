import type { DefaultTheme } from 'vitepress'
import type { CustomTheme } from '../types'

import { defineConfigWithTheme } from 'vitepress'
import {
  useSidebar,
  useHead,
  useNavigation,
  useSocialLinks,
} from './composition'

/**
 * NOTE: Need use defineConfig when i18n type fixed,
 * defineConfigwithTheme is used with custom theme.
 */
export default defineConfigWithTheme<DefaultTheme.Config & CustomTheme.Config>({
  srcDir: './source', // relative to project root(vitepress-blog/docs/).
  lang: 'zh-CN',
  title: '喵玉殿新版论坛开发文档',
  description: '喵玉殿新版论坛开发文档',
  head: useHead(),
  appearance: 'dark', // 外观
  lastUpdated: true,

  markdown: {
    lineNumbers: true,
  },

  themeConfig: {
    logo: '/images/logo.ico',
    nav: useNavigation(),
    sidebar: useSidebar(),
    socialLinks: useSocialLinks(),
    outline: { label: '目录' },
  },
})
