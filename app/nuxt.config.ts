import fs from 'fs'
import path from 'path'

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-12-20',
  modules: [
    '@nuxt/ui',
    '@nuxt/content',
    '@nuxtjs/tailwindcss',
    'nuxt-security',
    "nuxt-viewport",
    '@nuxtjs/i18n',
    'nuxt-lodash'
  ],
  // extends: ['@nuxt/ui-pro'],
  components: [{
    path: '~/components',
    pathPrefix: false,
  }],
  routeRules: {
    '/': {redirect: '/islam_kamil'}
  },
  // nitro: {
  //   prerender: {
  //     routes: getContentRoutes1()
  //   }
  // },
  css: ['@/assets/css/main.css', '@/assets/css/fonts.css'],
  ui: {
    safelistColors: ['primary', 'gray', 'success', 'error'],
  },
  content: {
    navigation: {
      fields: ['index', 'dir_index', 'icon', 'description', 'settings']
    }
  },
  i18n: {
    vueI18n: './i18n.config.ts'
  },
  viewport: {
    breakpoints: {
      xs: 320,
      sm: 640,
      md: 768,
      lg: 1024,
      xl: 1280,
      '2xl': 1536,
    },
    defaultBreakpoints: {
      desktop: 'lg',
      mobile: 'xs',
      tablet: 'md',
    },
  },
  security: {
    headers: {
      crossOriginEmbedderPolicy: process.env.NODE_ENV === 'development' ? 'unsafe-none' : 'require-corp',
    },
  },
  devtools: {
    enabled: true,
  },
  vite: {
    logLevel: 'info',
  },
})

function getContentRoutes(baseDir = path.resolve(__dirname, 'content'), prefix = '') {
  const routes = []
  const entries = fs.readdirSync(baseDir, { withFileTypes: true })

  for (const entry of entries) {
    if (entry.isDirectory()) {
      const newPrefix = prefix + '/' + entry.name
      routes.push(...getContentRoutes(path.join(baseDir, entry.name), newPrefix))
    } else {
      const ext = path.extname(entry.name)
      const nameWithoutExt = path.basename(entry.name, ext)
      if (ext === '.json') {
        let route = prefix
        if (nameWithoutExt !== 'index') {
          route += '/' + nameWithoutExt
        }
        if (route === '') {
          route = '/'
        }
        routes.push(route)
      }
    }
  }

  return routes
}

function getContentRoutes1() {
  const contentDir = path.resolve(__dirname, 'content')
  const files = fs.readdirSync(contentDir).filter(file => file.endsWith('.json'))
  return files.map(file => '/' + file.replace(/\.json$/, ''))
}