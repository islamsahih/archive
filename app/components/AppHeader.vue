<script setup lang="ts">
import type {NavItem} from "@nuxt/content";

const appConfig = useAppConfig()
const route = useRoute()

const navigation = inject<Ref<NavItem[]>>('navigation', ref([]))

const headerLinks = mapMenuItems(appConfig.header.links, navigation.value, {topLabel: 'short', label: 'full'})
const barLinks = headerLinks.map(link => ({label: link.label, to: link.to}))
const panelLinks = headerLinks.map(
  link => [link].concat(link.children ?? [])
)

const SCROLL_HIDE_THRESHOLD = 100
const isHidden = ref(false);
const {$viewport} = useNuxtApp()
let lastScrollTop = 0;
let lastFirstScreen = false;
const handleScroll = useThrottle(() => {
  const scrollTop = window.scrollY || document.documentElement.scrollTop;
  const firstScreen = scrollTop < window.innerHeight
  if (Math.abs(scrollTop - lastScrollTop) > SCROLL_HIDE_THRESHOLD || (firstScreen !== lastFirstScreen)) {
    isHidden.value = $viewport.isLessThan('md') && !firstScreen && scrollTop > lastScrollTop;
    lastScrollTop = scrollTop <= 0 ? 0 : scrollTop;
  }
  lastFirstScreen = firstScreen;
}, 100);

if (import.meta.client) {
  onMounted(() => {
    window.addEventListener('scroll', handleScroll);
  });

  onUnmounted(() => {
    window.removeEventListener('scroll', handleScroll);
  });
}

const ui = {
  header: {
    //wrapper: 'static md:sticky',
    container: 'max-w-[var(--content-width)]'
  },
  logo: 'w-auto h-10',
  l: {
    wrapper: 'flex flex-row items-center gap-1 h-20',
    icon: 'text-primary-600 inline-block aspect-[1/1] w-full h-full',
    text: {
      wrapper: 'flex flex-col h-full font-black tracking-wide text-gray-700 dark:text-gray-200',
      top: 'grow',
      bottom: 'grow'
    }
  },
  links: {
    wrapper: 'hidden lg:flex gap-x-16',
    base: 'text-md font-medium',
    trailingIcon: {
      active: ''
    },
  },
  popover: {
    navigation: {
      wrapper: 'p-2 space-y-1',
      base: 'group/link items-start',
      inactive: 'text-gray-700 dark:text-gray-200',
      size: 'text-md',
      label: 'text-base font-medium inline-block relative',
      description: 'text-sm leading-snug text-gray-500 dark:text-gray-400 line-clamp-2 font-normal relative',
      icon: {
        base: 'w-5 h-5 flex-shrink-0 mt-0.5',
        active: 'text-primary-500 dark:text-primary-400',
        inactive: 'text-primary-500 group-hover:text-primary-500 group-hover/link:text-primary-500 dark:text-primary-400 dark:group-hover:text-primary-400 dark:group-hover/link:text-primary-400',
      }
    }
  },
  panel: {
    navigation: {
      base: 'gap-2',
      label: 'text-wrap relative',
      active: 'before:bg-gray-200',
      inactive: 'text-gray-700 dark:text-gray-200',
      size: 'text-lg',
      icon: {
        active: 'text-primary-500 dark:text-primary-400',
        inactive: 'text-primary-500 group-hover:text-primary-500 dark:text-primary-400 dark:group-hover:text-primary-400',
      }
    }
  }
}
</script>

<template>
  <Header id="AppHeader" :class="{ 'hidden-up': isHidden }" :ui="ui.header">
    <template #logo>
      <Logo :class="ui.logo"/>

      <!--      <Logo :class="ui.logo"-->
      <!--            icon="i-lets-icons-book-open-duotone-line"-->
      <!--            text-top="Islam"-->
      <!--            text-bottom="Sahih"-->
      <!--            :v="1"-->
      <!--      />-->
    </template>

    <template #center>
      <HeaderLinks :links="barLinks" :ui="ui.links">
        <template #panel="{ link, close }">
          <UVerticalNavigation
            :links="link.children"
            :ui="ui.popover.navigation"
            @click="close"
          >
            <template #default="{ link }">
              <p>
                <span :class="ui.popover.navigation.label">
                  {{ link.label }}
                </span>
                <span v-if="link.description" :class="ui.popover.navigation.description">
                  {{ link.description }}
                </span>
              </p>
            </template>
          </UVerticalNavigation>
        </template>
      </HeaderLinks>
    </template>

    <template #right>
      <ColorModeButton/>
    </template>

    <template #panel>
      <UVerticalNavigation
        :links="panelLinks"
        :ui="ui.panel.navigation"
      />
    </template>

    <template #bottom>
      <HeaderGradient />
    </template>
  </Header>
</template>

<style scoped>
#AppHeader {
  transition: transform 0.3s ease;
}
#AppHeader.hidden-up {
  transform: translateY(-100%)
}
</style>