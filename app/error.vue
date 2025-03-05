<script setup lang="ts">
import type {NuxtError} from '#app'

const props = defineProps({
  error: Object as () => NuxtError
})
</script>

<template>
    <div class="flex flex-col items-center">
      <PageHeader :title="`${error.statusCode}`"
                   icon="i-lets-icons-alarm-duotone"
                   :description="error.statusMessage"
      />
      <div>
        {{error}}
      </div>
    </div>
</template>