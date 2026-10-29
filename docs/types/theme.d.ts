import type { DefaultTheme } from 'vitepress'

interface ThemeLocaleConfig {
  nav: DefaultTheme.NavItem[]
}

export namespace CustomTheme {
  interface Config extends DefaultTheme.Config {
    locales?: Record<string, ThemeLocaleConfig>
  }
}
