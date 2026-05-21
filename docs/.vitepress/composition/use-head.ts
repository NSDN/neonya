import type { HeadConfig } from 'vitepress'

// See: <https://vitepress.dev/zh/reference/site-config#head>
export default function useHead(): HeadConfig[] {
  return [['link', { rel: 'icon', href: '/images/logo.ico' }]]
}
