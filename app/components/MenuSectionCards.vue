<script setup lang="ts">
import type {PropType} from 'vue'

const props = defineProps({
  links: {
    default: () => useNuxtApp().$menu.children(useRoute().path)
  },
  class: {
    type: [String, Object, Array] as PropType<any>,
    default: undefined
  },
  ui: {
    type: Object as PropType<Partial<typeof config>>,
    default: () => ({})
  }
})

const config = {
  grid: "grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 2xl:grid-cols-4 gap-8",
  card: {
    base: 'rounded-lg divide-y divide-gray-200 dark:divide-gray-800 ring-1 ring-gray-200 dark:ring-gray-800 shadow bg-white dark:bg-gray-900 relative group hover:ring-2 hover:bg-gray-100/50 dark:hover:bg-gray-800/50 hover:ring-secondary-500 dark:hover:ring-secondary-400',
  }
}

const {ui, attrs} = useUI('CardsGrid', toRef(props, 'ui'), config, toRef(props, 'class'), true)
</script>

<template>
  <div :class="ui.grid">
    <UCard v-for="link in links" :ui="ui.card">
      <NuxtLink v-if="link.to" :to="link.to" :aria-label="link.label" class="focus:outline-none" tabindex="-1">
        <span class="absolute inset-0" aria-hidden="true" />
      </NuxtLink>

      <div class="flex flex-row gap-4 mb-2 ">
        <div v-if="link.icon">
          <UIcon :name="link.icon" class="w-10 h-10 flex-shrink-0 text-primary"/>
        </div>
        <p v-if="link.label || link.title"
           class="flex items-center text-2xl font-semibold">
          {{ link.label || link.title }}
        </p>
      </div>
      <div v-if="link.description" class="mt-4 text-gray-500 dark:text-gray-400">
        {{ link.description }}
      </div>
    </UCard>
  </div>
</template>
