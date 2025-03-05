export default defineNuxtPlugin(() => {
  const appConfig = useAppConfig()
  const nuxtApp = useNuxtApp()

  const makeVariablesWithPrefix = (src, prefix) =>
    Object.entries(src).map(([key, value]) => `--${prefix}-${key}: ${value};`).join('\n')

  const root = computed(() => {
    return `:root {
  ${Object.entries(appConfig.ui.variables.root).map(([key, group]) => makeVariablesWithPrefix(group, key)).join('\n')}
}

.light {
  ${makeVariablesWithPrefix(appConfig.ui.variables.light, 'ui')}
}

.dark {
  ${makeVariablesWithPrefix(appConfig.ui.variables.dark, 'ui')}
}`
  })

  // Head
  const headData: any = {
    style: [{
      innerHTML: () => root.value,
      tagPriority: -2,
      id: 'nuxt-ui-variables'
    }]
  }

  // SPA mode
  if (import.meta.client && nuxtApp.isHydrating && !nuxtApp.payload.serverRendered) {
    const style = document.createElement('style')

    style.innerHTML = root.value
    style.setAttribute('data-nuxt-ui-variables', '')
    document.head.appendChild(style)

    headData.script = [{
      innerHTML: 'document.head.removeChild(document.querySelector(\'[data-nuxt-ui-variables]\'))'
    }]
  }

  useHead(headData)
})
