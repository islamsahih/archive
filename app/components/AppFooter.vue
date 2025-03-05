<script setup lang="ts">
import type {NavItem} from "@nuxt/content";

const appConfig = useAppConfig()
const { $copyright } = useNuxtApp()

const navigation = inject<Ref<NavItem[]>>('navigation', ref([]))

const footerLinks = mapMenuItems(appConfig.footer.links, navigation.value)

const ui = {
  wrapper: 'border-t border-gray-200 dark:border-gray-800 -mt-px',
  top: {
    wrapper: 'hidden'
  },
  bottom: {
    wrapper: '',
    container: 'py-8 md:py-4 md:flex md:items-center md:justify-between md:gap-x-0 md:h-[--footer-height]',
    left: 'flex items-center justify-center md:justify-start md:flex-1 gap-x-1.5 mt-4 md:mt-0 md:order-1',
    center: 'hidden',
    right: 'md:flex-1 flex items-center justify-center md:justify-end gap-x-1.5 md:order-3'
  },
  links: {
    wrapper: 'text-center'
  },
}
</script>

<template>
  <Footer :ui="ui">
    <template #left>
      <ClientOnly>
        <p class="text-sm text-gray-500 dark:text-gray-400 hover:text-gray-600 dark:hover:text-gray-300">
          {{ $copyright }}
        </p>
      </ClientOnly>
    </template>

    <template #right>
      <FooterLinks :links="footerLinks" :ui="ui.links"/>
    </template>
  </Footer>
</template>
