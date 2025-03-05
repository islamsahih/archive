<script setup lang="ts">
import type {NavItem} from "@nuxt/content";
import {htmlToText} from "html-to-text";
import MiniSearch from "minisearch";
import type {Ref} from "@vue/reactivity";

const appConfig = useAppConfig()
const router = useRouter()
const route = useRoute()
const {$viewport} = useNuxtApp()
const $t = useTranslate()
const {navPageFromPath} = useContentHelpers()

const navigation = inject<Ref<NavItem[]>>('navigation', ref([]))

const page = navPageFromPath(route.path, navigation.value)
if (!page) {
  throw createError({statusCode: 500, fatal: true})
}
useHead({
  ...page,
  meta: [{
    name: 'robots',
    content: 'noindex, nofollow'
  }]
})

const dirList = dirListFromNavigation(navigation.value)
const sectionList = [{label: $t('search.allSections'), _path: '/'}, ...dirList.map(link => ({
  ...link,
  label: pageSpecialTitle(link, 'full', undefined, navigation.value),
}))]
const searchSection = ref(sectionList[0]._path)

const inTitle = ref(true)
const inText = ref(true)
const searchFields = computed(() => [
  inTitle.value && 'title',
  inText.value && 'content',
].filter(v => v))
const searchStoreFields = ['path', 'title', 'date', 'description', 'icon', 'content', 'audio', 'video']

const searchQuery = ref(route.query.q as string | undefined)
const searchResult = ref(undefined)
const searchInProgress = ref(false)

const {data: articles, refresh: fetchArticles} = await useAsyncData(
  'search-articles',
  async () => searchResult.value ?
    (await queryContent(searchSection.value).where({$not: {_file: {$contains: 'index.json'}}}).find())
      .map((article, index) => ({
        id: index + 1,
        path: article._path,
        title: article.title,
        description: article.description,
        date: article.date,
        audio: article.audio,
        video: article.video,
        content: htmlToText(String(article.body || ''), {
          selectors: [
            {selector: 'a', options: {ignoreHref: true}},
            {selector: 'img', format: 'skip'}
          ]
        }).replace(/\n+/g, ' ')
      })) : []
)

async function search() {
  const query = searchQuery.value?.trim()
  import.meta.client && await router.push({name: route.name, query: query && {q: query}})
  if (!query) {
    searchResult.value = undefined
    return
  }
  searchResult.value = []
  searchInProgress.value = true
  setTimeout(async () => {
    let miniSearch = new MiniSearch({
      fields: searchFields.value,
      storeFields: searchStoreFields,
      searchOptions: {
        prefix: true,
        fuzzy: 0.1,
      }
    })
    await fetchArticles()
    miniSearch.addAll((articles.value ?? []))
    searchResult.value = miniSearch.search(query, {combineWith: 'AND'})
    searchInProgress.value = false
  }, 1000)
}

const orderByOptions = [
  {value: 'score', type: 'number', defaultDirection: -1},
  {value: 'index', type: 'index', defaultDirection: 1},
  {value: 'title', type: 'string', defaultDirection: 1},
  {value: 'date', type: 'string', defaultDirection: -1}
].map(v => ({...v, label: $t('search.orderByOptions.' + v.value)}))
const orderBy = ref(orderByOptions[0])
const orderDirection = ref(orderBy.value.defaultDirection)
watch(orderBy, (value) => orderDirection.value = value.defaultDirection)

const resultPage = ref(1)
watch(resultPage, () => {
  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  });
})
const resultPageCountOptions = [25, 50, 100]
const resultPageCount = ref(resultPageCountOptions[0])
const resultRows = computed(() =>
  searchResult.value?.sort((x, y) => {
      const order = orderBy.value
      const direction = orderDirection.value
      if (order.type === 'index') {
        const xIndex = navPageFromPath(x.path, navigation.value).index ?? 0
        const yIndex = navPageFromPath(y.path, navigation.value).index ?? 0
        return direction * (xIndex - yIndex)
      }
      if (order.type === 'number') {
        return direction * (x[order.value] - y[order.value])
      }
      if (order.type === 'string') {
        return direction * (x[order.value].toLowerCase().localeCompare(y[order.value].toLowerCase()))
      }
      return 0
    }
  ).slice((resultPage.value - 1) * resultPageCount.value, (resultPage.value) * resultPageCount.value)
    .map(result => {
      const page = navPageFromPath(result.path, navigation.value)
      const settings = pageSettings(page, undefined, navigation.value)
      return {
        ...highlightMatches(useCloneDeep(result)),
        breadcrumbs: pageBreadcrumbs(page, navigation.value, settings, {
          titleType: $viewport.isLessThan('sm') ? 'short' : 'full', selfTitleType: 'short', noLinks: true,
        }),
        icon: result.icon ?? settings.icon ?? router.resolve(result.path)?.meta?.icon
      }
    })
)

await search()

const showSections = ref(true)
const showMatches = ref(true)
const showDates = ref(false)
const showAudios = ref(false)
const showVideos = ref(false)

const optionsVisible = ref(false)

const ui = {
  header: {
    wrapper: 'border-b-0',
    title: 'font-medium',
  },
  aside: {
    container: 'flex flex-row flex-wrap lg:flex-col gap-x-5 gap-y-3.5 mt-3.5 lg:mt-0',
    section: 'w-full max-w-60'
  },
  results: {
    list: 'flex flex-col flex-wrap py-4 gap-2.5 border-t border-gray-200 dark:border-gray-700 break-words',
    link: 'not-prose text-primary-500 dark:text-primary-400 font-medium',
    icon: 'flex-shrink-0 w-5 h-5 align-sub mr-2',
    title: 'text-primary-500 dark:text-primary-400 text-base font-semibold dark:font-medium',
    description: '',
    media: {
      container: 'flex flex-col sm:flex-row items-start sm:items-center gap-4',
      audio: 'w-full h-9',
      video: 'bg-red-500 hover:bg-red-600 dark:bg-red-400 dark:hover:bg-red-500',
    },
    text: 'text-base',
    date: 'text-sm font-semibold text-gray-500 dark:text-gray-400',
  },
  status: {
    container: 'flex flex-col items-center justify-center flex-1 px-6 py-14 sm:px-14',
    icon: 'mx-auto mb-4 w-10 h-10 text-primary dark:text-primary',
    text: 'text-sm text-center text-gray-900 dark:text-white',
  }
}

</script>

<template>
  <Page>
    <template v-if="$viewport.isGreaterOrEquals('lg') || optionsVisible" #aside>
      <div :class="ui.aside.container">
        <LabelGroup :label="$t('search.inSection')" :class="ui.aside.section">
          <USelectMenu
            v-model="searchSection" :options="sectionList" value-attribute="_path">
            <template #option="{ option: section }">
              <span v-for="_ in section.depth" class="ml-2" aria-hidden="true"/>
              <span>{{ section.label }}</span>
            </template>
          </USelectMenu>
        </LabelGroup>
        <div class="flex flex-row flex-wrap gap-x-5 gap-y-3.5">
          <LabelGroup :label="$t('search.in')" class="">
            <UCheckbox id="in-titles" v-model="inTitle" :label="$t('search.inOptions.title')"/>
            <UCheckbox id="in-contents" v-model="inText" :label="$t('search.inOptions.text')"/>
          </LabelGroup>
          <LabelGroup :label="$t('search.show')" class="">
            <UCheckbox id="show-matches" v-model="showMatches" :label="$t('search.showOptions.matches')"/>
            <UCheckbox id="show-sections" v-model="showSections" :label="$t('search.showOptions.sections')"/>
            <UCheckbox id="show-dates" v-model="showDates" :label="$t('search.showOptions.dates')"/>
            <UCheckbox id="show-audios" v-model="showAudios" :label="$t('search.showOptions.audios')"/>
            <UCheckbox id="show-videos" v-model="showVideos" :label="$t('search.showOptions.videos')"/>
          </LabelGroup>
        </div>
      </div>
    </template>

    <template #header>
      <PageHeader :title="page.title"
                  :icon="page.icon ?? appConfig.ui.icons.search"
                  :ui="ui.header"
      />
      <SearchInput v-model="searchQuery" @search="search" @reset="search" :disabled="searchInProgress">
        <template v-if="$viewport.isLessThan('lg')" #trailing>
          <UButton
            icon="i-mdi-tune-variant"
            color="primary"
            variant="link"
            :padded="false"
            @click="optionsVisible = !optionsVisible"
          />
        </template>
      </SearchInput>
    </template>

    <template #default>
      <PageBody>
        <div v-if="searchResult" class="flex flex-col gap-2.5">
          <template v-if="resultRows.length">
            <div class="flex flex-row flex-wrap justify-between gap-2.5">
              <Label class="self-center" :text="`${$t('search.found')}: ${searchResult.length}`"/>
              <div class="flex flex-row gap-x-2.5">
                <UButton variant="ghost"
                         :icon="orderDirection === 1 ? appConfig.ui.icons.sort_asc : appConfig.ui.icons.sort_desc"
                         @click="orderDirection *= -1"
                />
                <USelectMenu v-model="orderBy" :options="orderByOptions" :ui-menu="{ width: 'min-w-max' }"/>
              </div>
            </div>
            <div v-for="row in resultRows" :class="ui.results.list">
              <UBreadcrumb v-if="showSections" :links="row.breadcrumbs" :ui="appConfig.ui.breadcrumbInactive"/>

              <NuxtLink :to="row.path" :class="ui.results.link">
                <UIcon v-if="row.icon" :name="row.icon" :class="ui.results.icon"/>

                <span v-if="showMatches && row.title_match" v-html="row.title_match" :class="ui.results.title"/>
                <span v-else :class="ui.results.title">{{ row.title }}</span>
              </NuxtLink>

              <p v-if="showMatches && row.description_match" :class="ui.results.description"
                 v-html="row.description_match"/>
              <p v-else-if="row.description" :class="ui.results.description">{{ row.description }}</p>

              <div v-if="(showAudios && row.audio) || (showVideos && row.video)" :class="ui.results.media.container">
                <audio v-if="showAudios && row.audio" preload="none" controls :class="ui.results.media.audio">
                  <source :src="row.audio">
                </audio>
                <UButton v-if="showVideos && row.video"
                         :to="row.video"
                         label="YouTube"
                         variant="solid"
                         icon="i-simple-icons-youtube"
                         :class="ui.results.media.video"
                         target="_blank"
                />
              </div>

              <p v-if="showMatches && row.content_match" :class="ui.results.text" v-html="row.content_match"/>

              <time v-if="showDates" :class="ui.results.date">{{ row.date }}</time>
            </div>
            <div v-if="resultRows.length < searchResult.length"
                 class="flex flex-col sm:flex-row items-end justify-end gap-2.5">
              <LabelGroup row :label="$t('search.resultsPerPage')">
                <USelect v-model="resultPageCount" :options="resultPageCountOptions"/>
              </LabelGroup>
              <UPagination v-model="resultPage"
                           :page-count="resultPageCount"
                           :total="searchResult?.length"
                           :max="5"
              />
            </div>
          </template>
          <template v-else>
            <UProgress v-if="searchInProgress" animation="carousel" size="xs"/>
            <div :class="ui.status.container">
              <UIcon :name="searchInProgress ? appConfig.ui.icons.loading : 'i-lets-icons-sad'"
                     :class="[ui.status.icon, {'animate-spin': searchInProgress}]"
              />
              <p :class="ui.status.text">
                {{ searchInProgress ? $t('search.progress') : $t('search.notFound') }}
              </p>
            </div>
          </template>
        </div>
      </PageBody>
    </template>
  </Page>
</template>