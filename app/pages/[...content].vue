<script setup lang="ts">
import type {NavItem} from "@nuxt/content";
import Mustache from "mustache";

const {navPageFromPath, navDirFromPath} = useContentHelpers()

definePageMeta({
  // path: '/:categories+/:index(\\d+)',
  // icon: 'i-lets-icons-book-open-duotone-line'
})

const appConfig = useAppConfig()
const router = useRouter()
const route = useRoute()
const $t = useTranslate()
const navigation = inject<Ref<NavItem[]>>('navigation', ref([]))

const page = navPageFromPath(route.path, navigation.value)
if (!page) {
  throw createError({statusCode: 500, fatal: true})
}
const parent = navPageFromParentPath(route.path, navigation.value)
const settings = pageSettings(page, parent)
const children = navDirFromPath(route.path, navigation.value)?.sort((x, y) => {
  if (x.dir_index && y.dir_index) {
    return x.dir_index - y.dir_index
  }
  if (x.index && y.index) {
    return x.index - y.index
  }
  return x.title.toLowerCase().localeCompare(y.title.toLowerCase())
}).map(link => ({
  ...link,
  label: pageSpecialTitle(link, 'list', undefined, undefined, page),
  ...(settings.article?.lowercase_descriptions && {lowercaseDescription: lowercaseFirstLetter(link.description)})
}))
const isList = !settings.no_children_list && children

const {data: doc} = await useAsyncData(route.path, () => queryContent(route.path).findOne())
if (!doc.value) {
  throw createError({statusCode: 404, fatal: true})
}
useContentHead(doc)
const articleRenderData = {
  children: children,
  source_ar: doc.value.meta?.source?.ar?.join(', '),
  source_ru: doc.value.meta?.source?.ru?.join(', '),
  social: appConfig.seo.social,
  admin: appConfig.seo.social.filter(link => link.admin),
  book: appConfig.seo.social.filter(link => link.book),
}
const article = doc.value.body && Mustache.render(String(doc.value.body), articleRenderData)
const editNotificationText = settings.notify_edit && Mustache.render(String($t('common.editNotification.text')), articleRenderData)

const references = doc.value.references?.map(path => {
  const page = navPageFromPath(path, navigation.value)
  const settings = pageSettings(page, undefined, navigation.value)
  return {
    breadcrumbs: pageBreadcrumbs(page, navigation.value, settings, {
      titleType: 'full', selfTitleType: 'short', noLinks: true,
    }),
    label: page.title,
    to: path,
    icon: page.icon ?? settings.icon ?? router.resolve(page._path)?.meta?.icon
  }
})
const tags = doc.value.tags

const breadcrumbs = pageBreadcrumbs(page, navigation.value, settings, {
  titleType: 'full', selfTitleType: 'short',
})
const surround = parent && !children?.length && await (async () => {
  let prev = undefined;
  let next = undefined;
  (await useAsyncData(`${route.path}-surround`, () =>
          queryContent().where({
            $and: [
              {navigation: {$ne: false}},
              {_path: {$regex: `^${parent._path}/.*`}},
              {dir_index: {$in: [page.dir_index - 1, page.dir_index + 1]}}
            ]
          }).limit(2).find(),
      {default: () => []}
  )).data.value.forEach(link => {
    if (link.dir_index < page.dir_index) {
      prev = link
    } else {
      next = link
    }
  })
  if (!prev) {
    prev = parent
  }
  if (!next) {
    next = getNextChild(parent, page.dir_index, navigation.value)
  }
  return (next ? [prev, next] : [prev]).map(link => ({
    ...link,
    ...((short) => short != link.title ? {description: short} : undefined)
    (pageSpecialTitle(link, 'short', undefined, navigation.value)),
  }))
})()

const headerLinks = surround && surround?.map((link, i) => ({
  to: link._path,
  icon: i === 0 ? appConfig.ui.icons.arrow_left : appConfig.ui.icons.arrow_right,
  target: '_self'
}))

const searchQuery = ref('')

async function search() {
  await router.push({name: 'search', params: {query: searchQuery.value.trim()}})
}

const ui = {
  header: {
    base: {
      title: 'font-medium',
      links: 'flex flex-row flex-nowrap gap-2'
    },
    links: {
      base: 'inline-flex items-center rounded-full p-1.5 text-gray-900 dark:text-white hover:text-primary bg-gray-100 dark:bg-gray-800 hover:bg-primary/10 ring-0',
      // base: 'inline-flex items-center rounded-full p-1.5 bg-gray-100 dark:bg-gray-800 group-hover:bg-primary/10 ring-1 ring-gray-300 dark:ring-gray-700 mb-4 group-hover:ring-primary/50',
      icon: 'w-5 h-5'
    },
  },
  editNotification: 'mt-4 sm:mt-8',
  tags: {
    wrapper: 'flex flex-row flex-wrap items-start justify-start my-[2rem] gap-2',
  },
  references: {
    wrapper: 'flex flex-col items-start justify-start my-[2rem] gap-2',
    title: 'font-medium text-gray-900 dark:text-white',
    list: 'flex flex-col',
    link: 'text-primary-500 dark:text-primary-400 font-semibold dark:font-medium',
    icon: 'flex-shrink-0 w-5 h-5 align-sub mr-2'
  }
}
</script>

<template>
  <Page>
    <!--    <template #left>-->
    <!--      <Aside>-->
    <!--        <template #top>-->
    <!--          <SearchInput v-model="searchQuery" @search="search"/>-->
    <!--        </template>-->

    <!--        &lt;!&ndash;        <UNavigationTree :links="asideLinks" default-open/>&ndash;&gt;-->
    <!--      </Aside>-->
    <!--    </template>-->

    <template #default>
      <PageHeader :title="page.title"
                  :icon="page.icon ?? settings.icon ?? route.meta?.icon"
                  :description="doc.description ?? page.description"
                  :ui="ui.header.base"
      >
        <template #headline>
          <UBreadcrumb :links="breadcrumbs"/>
        </template>
        <template v-if="headerLinks" #links>
          <NuxtLink v-for="link in headerLinks" :to="link.to" :class="ui.header.links.base">
            <UIcon :name="link.icon" :class="ui.header.links.icon"/>
          </NuxtLink>
        </template>
      </PageHeader>

      <UAlert v-if="settings.notify_edit"
        :title="$t('common.editNotification.title')"
        :icon="appConfig.ui.icons.edit_notification"
        :class="ui.editNotification"
      >
        <template #description>
          <div class="text-sm font-medium" v-html="editNotificationText">
          </div>
        </template>
      </UAlert>

      <PageBody v-if="isList" prose>
        <ul class="list-none pl-0">
          <li v-for="link in children">
            <NuxtLink :to="link._path">{{ link.label }}</NuxtLink>
          </li>
        </ul>
      </PageBody>

      <PageBody v-else prose>
        <div v-if="doc.audio || doc.video"
             class="flex flex-col sm:flex-row items-end sm:items-center sm:justify-end w-full md:w-2/3 xl:w-1/2 ml-auto gap-4 mb-4 sm:mb-8">
          <audio v-if="doc.audio" preload="none" controls class="w-full h-9">
            <source :src="doc.audio">
          </audio>
          <UButton v-if="doc.video"
                   :to="doc.video"
                   label="YouTube"
                   variant="solid"
                   class="bg-red-500 hover:bg-red-600 dark:bg-red-400 dark:hover:bg-red-500"
                   icon="i-simple-icons-youtube"
                   target="_blank"
          />
        </div>

        <article v-if="article" v-html="article"/>

        <div v-if="tags" :class="ui.tags.wrapper">
          <UBadge v-for="tag in tags" color="primary" variant="subtle">{{ tag }}</UBadge>
        </div>

        <div v-if="references" :class="ui.references.wrapper">
          <span :class="ui.references.title">{{ $t('common.relatedPages') }}</span>
          <div v-for="link in references" :class="ui.references.list">
            <UBreadcrumb :links="link.breadcrumbs" :ui="appConfig.ui.breadcrumbInactive"/>
            <NuxtLink :to="link.to" :class="ui.references.link">
              <UIcon :class="ui.references.icon" :name="link.icon"/>
              {{ link.label }}
            </NuxtLink>
          </div>
        </div>

        <hr/>

        <ContentSurround v-if="surround" class="mt-8" :surround="surround"/>
      </PageBody>
    </template>
  </Page>
</template>
