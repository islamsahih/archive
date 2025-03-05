<script setup lang="ts">
const { data: navigation } = await useAsyncData('navigation', () => fetchContentNavigation())
provide('navigation', navigation)
</script>

<template>
  <AppHeader/>

  <AppMain>
    <NuxtLayout>
      <NuxtPage/>
    </NuxtLayout>
  </AppMain>

  <AppFooter/>

  <UNotifications/>
</template>