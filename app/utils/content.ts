import {splitByCase, upperFirst} from 'scule'
import type {NavItem, ParsedContent} from '@nuxt/content'
import type {NavigationTree} from '~/types/ui-pro-types'
import Mustache from "mustache";

const {navPageFromPath, navDirFromPath} = useContentHelpers()

export function mapContentNavigation(navigation: NavItem[], options?: {
  labelAttribute?: string,
  onlyAttributes?: string[],
  skipAttributes?: string[]
}): NavigationTree[] {
  const navMap = {
    [options?.labelAttribute || 'title']: 'label',
    '_path': 'to'
  }
  const onlyAttributes = options?.onlyAttributes
  const skipAttributes = options?.skipAttributes

  return navigation.map((navLink) => {
    const link = {} as NavigationTree
    for (const key in navLink) {
      if (onlyAttributes?.length && (!onlyAttributes.includes(key) && !onlyAttributes.includes(navMap[key]))) {
        continue
      }
      if (skipAttributes?.length && (skipAttributes.includes(key) || skipAttributes.includes(navMap[key]))) {
        continue
      }

      if (key === 'children') {
        link.children = navLink.children?.length ? mapContentNavigation(navLink.children, options) : undefined
        continue
      }
      if (navLink[key]) {
        // @ts-ignore
        link[navMap[key] || key] = navLink[key]
      }
    }
    if (!link.label && navLink.title) {
      link.label = navLink.title
    }
    return link
  })
}

export function findPageHeadline(page: ParsedContent): string {
  return page._dir?.title ? page._dir.title : splitByCase(page._dir).map(p => upperFirst(p)).join(' ')
}

function findPageBreadcrumb(navigation?: NavItem[], page?: ParsedContent | NavItem, options?: {
  noSelf: boolean
}): NavItem[] {
  if (!navigation || !page) {
    return []
  }

  return navigation.reduce((breadcrumb: NavItem[], link: NavItem) => {
    if (page._path && (page._path + '/').startsWith(link._path + '/')) {
      const self = page._path.length === link._path.length
      if (self && options?.noSelf) {
        return breadcrumb
      }
      if (link.children) {
        breadcrumb.push(link)
        if (!self) {
          breadcrumb.push(...findPageBreadcrumb(link.children, page, options))
        }
      }
    }
    return breadcrumb
  }, [])
}

export function pageBreadcrumbs(page: NavItem | ParsedContent, navigation: NavItem[], settings?: any, options?: {
  titleType?: string,
  selfTitleType?: string,
  noSelf?: boolean,
  noLinks?: boolean,
}) {
  let breadcrumbs = []
  let links = findPageBreadcrumb(navigation, page, {noSelf: true})
  for (let i = 0; i < links.length; i++) {
    const link = links[i]
    breadcrumbs.push({
      label: options?.titleType ?
        pageSpecialTitle(links[i], options.titleType, undefined, navigation, i ? links[i - 1] : undefined) :
        link.title,
      ...(!options?.noLinks && {to: link._path})
    })
  }
  if (breadcrumbs.length && !options?.noSelf) {
    if (!settings) {
      settings = pageSettings(page, undefined, navigation)
    }
    if (!settings.no_self_breadcrumb) {
      breadcrumbs.push({
        label: options?.selfTitleType ?
          pageSpecialTitle(page, options.selfTitleType, settings) :
          page.title,
        ...(!options?.noLinks && {to: page.path})
      })
    }
  }
  return breadcrumbs
}

export function pathBreadcrumbs(path: string, navigation: NavItem[], options?: {
  titleType?: string,
  selfTitleType?: string,
  noSelf?: boolean,
  noLinks?: boolean,
}) {
  return pageBreadcrumbs(navPageFromPath(path, navigation), navigation, undefined, options)
}

export function mapMenuItems(items: (string | Object)[], navigation: NavItem[], options?: {
  label?: string,
  topLabel?: string,
}, depth = 0, parent?: NavItem): NavItem[] {
  const label = (depth === 0 && options?.topLabel) ? options?.topLabel : options?.label

  const makeLink = (path, parent?: NavItem) => {
    const navLink = navPageFromPath(path, navigation)
    let link = {} as NavItem
    for (const key in navLink) {
      if (key === 'children') {
        continue
      }
      link[key] = navLink[key]
    }
    link.label = label ? pageSpecialTitle(link, label, false, navigation, parent) : link.title
    link.to = link._path
    link.depth = depth
    return link
  }

  let links = [] as NavItem[]
  for (const item of items) {
    switch (typeof item) {
      case 'string':
        links.push(makeLink(item, parent))
        break
      case 'object':
        const keys = Object.keys(item)
        if (keys.length === 1) {
          const link = makeLink(keys[0])
          links.push({
            ...link,
            children: mapMenuItems(item[keys[0]], navigation, options, depth+1, link)
          })
        }
        break
    }
  }
  return links
}

export function dirTreeFromNavigation(navigation: NavItem[], depth = 0): NavItem[] {
  let dirs = [] as NavItem[]
  for (const navLink of navigation) {
    if (!navLink.children) {
      continue
    }
    let link = {} as NavItem
    for (const key in navLink) {
      if (key === 'children') {
        const children = dirTreeFromNavigation(navLink.children, depth + 1)
        if (children.length) {
          link.children = children.sort((x, y) => (x.index ?? 0) - (y.index ?? 0))
        }
        continue
      }
      link[key] = navLink[key]
    }
    link.depth = depth
    dirs.push(link)
  }
  return dirs.sort((x, y) => (x.index ?? 0) - (y.index ?? 0))
}

export function dirListFromNavigation(navigation: NavItem[]): NavItem[] {
  const getDirList = (navigation: NavItem[], depth = 0) => {
    let dirs = [] as NavItem[]
    for (const navLink of navigation) {
      if (!navLink.children) {
        continue
      }
      let link = {} as NavItem
      for (const key in navLink) {
        if (key !== 'children') {
          link[key] = navLink[key]
        }
      }
      link.depth = depth
      dirs.push(link)
      if (navLink.children?.length) {
        dirs.push(...getDirList(navLink.children, depth + 1))
      }
    }
    return dirs
  }

  return getDirList(navigation).sort((x, y) => (x.index ?? 0) - (y.index ?? 0))
}

export function parentDir(path: string): string {
  return path.substring(0, path.lastIndexOf('/'))
}

export function navPageFromParentPath(path: string, navigation: NavItem[]): NavItem | undefined {
  return navPageFromPath(parentDir(path), navigation)
}

export function navPageFromParent(page: NavItem, navigation: NavItem[]): NavItem | undefined {
  return navPageFromParentPath(page._path, navigation)
}

export function navDirFromParentPath(path: string, navigation: NavItem[]): NavItem[] | undefined {
  return navDirFromPath(parentDir(path), navigation)
}

export function navDirFromParent(page: NavItem, navigation: NavItem[]): NavItem[] | undefined {
  return navDirFromParentPath(page._path, navigation)
}

// export function findPageIcon(navigation: NavItem[], page: ParsedContent, section?: NavItem): string {
//   if (page.icon) {
//     return page.icon
//   }
//   if (!section) {
//     section =
//   }
// }

export function pageSettingsFromPath(path: string, navigation: NavItem[]) {
  return pageSettings(navPageFromPath(path, navigation), undefined, navigation)
}

export function pageSettings(page: NavItem | ParsedContent, parent?: NavItem, navigation?: NavItem[]) {
  if (!parent && navigation) {
    parent = navPageFromParentPath(page._path, navigation)
  }
  return {
    short_title_template: '{{title}}',
    full_title_template: '{{title}}',
    list_title_template: '{{dir_index}}: {{title}}',
    ...parent?.settings?.children,
    ...page.settings,
  }
}

export function pageSpecialTitleByPath(path: string, type: string, navigation: NavItem[]): string {
  return pageSpecialTitle(navPageFromPath(path, navigation), type, undefined, navigation)
}

export function pageSpecialTitle(page: NavItem | ParsedContent, type: string, settings?: Object, navigation?: NavItem[], parent?: NavItem): string {
  if (!settings) {
    if (!parent && !navigation) {
      return ''
    }
    settings = pageSettings(page, parent, navigation)
  }
  return settings[`${type}_title`] ?? Mustache.render(settings[`${type}_title_template`], page) ?? page.title
}

export function lowercaseFirstLetter(str) {
  if (!str) return str;
  return str.charAt(0).toLowerCase() + str.slice(1);
}

export function getNextChild(dir: NavItem, dir_index: number, navigation: NavItem[]): NavItem | undefined {
  if (!dir) {
    return undefined
  }
  if (dir.children) {
    for (let page of [...dir.children].sort((x, y) => x.dir_index - y.dir_index)) {
      if (page.dir_index > dir_index) {
        return page
      }
    }
  }
  return getNextChild(navPageFromParent(dir, navigation), dir.dir_index, navigation)
}