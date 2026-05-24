import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import type { Plugin } from 'vite'

const MOBILE_BREAKPOINT = 'max-width: 48em'

function customMediaPlugin(): Plugin {
  return {
    name: 'custom-media',
    enforce: 'pre',
    transform(code: string, id: string) {
      if (id.endsWith('.vue') || id.endsWith('.css')) {
        return code.replace(
          /@media\s*\(\s*--mobile\s*\)/g,
          `@media (${MOBILE_BREAKPOINT})`
        )
      }
    }
  }
}

export default defineConfig({
  plugins: [customMediaPlugin(), vue()],

  resolve: {
    alias: { '@': '/src' }
  },

  server: {
    port: 10123
  }
})
