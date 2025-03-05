export default defineNuxtPlugin(async () => {
  const appConfig = useAppConfig()
  const router = useRouter()
  const {navPageFromPath} = useContentHelpers()
  const {data: navigation} = await useAsyncData('navigation', () => fetchContentNavigation())

  function makeLinks(items: any[]): any[] {
    const makeLink = (path: string, children?: any[]): any => {
      if (children) {
        children = makeLinks(children)
      }
      const doc = navPageFromPath(path, navigation.value)
      const route = router.resolve(path)
      return {
          to: doc?._path && route.path,
          title: doc?.title ?? route.meta.title ?? '',
          short: doc?.short ?? route.meta.short ?? '',
          label: doc?.short ?? doc?.title ?? route.meta.short ?? route.meta.title ?? '',
          description: doc?.description ?? route.meta.description ?? '',
          icon: doc?.icon ?? route.meta.icon ?? '',
          ...(children?.length ? {children: children} : {})
        }
    }
    let links = []
    for (const item of items) {
      switch (typeof item) {
        case 'string':
          links.push(makeLink(item))
          break
        case 'object':
          const keys = Object.keys(item)
          if (keys.length === 1) {
            links.push(makeLink(keys[0], item[keys[0]]))
          }
          break
      }
    }
    return links
  }

  function findSection(links: any[], path: string): any {
    for (const link of links) {
      if (link.to === path) {
        return link
      }
      if (link.children) {
        const child = findSection(link.children, path)
        if (child) {
          return child
        }
      }
    }
    return false
  }

  const links = makeLinks(appConfig.menu.links)
  const footerLinks = makeLinks(appConfig.menu.footerLinks)

  return {
    provide: {
      menu: {
        links: links,
        footerLinks: footerLinks,
        section: (path: string) => findSection(links, path),
        children: (path: string) => findSection(links, path)?.children ?? []
      }
    }
  }
})
