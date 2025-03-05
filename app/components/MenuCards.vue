<script setup lang="ts">
const props = defineProps({
  links: {
    default: () => useNuxtApp().$menu.links
  }
})
</script>

<template>
  <div>
    <template v-for="section in links">
      <div class="flex flex-col py-10">
        <p v-if="section.label || section.title" class="text-center text-2xl font-semibold">
          {{ section.title || section.label }}
        </p>
        <p v-if="section.description" class="text-center text-gray-500 dark:text-gray-400 mt-4">
          {{ section.description }}
        </p>
        <MenuSectionCards :links="section.children ?? [section]" class="mt-8"/>
      </div>
    </template>
  </div>
</template>