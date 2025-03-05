<script setup lang="ts">
import {twMerge} from "tailwind-merge";

const props = defineProps({
    text: {
      type: String,
      default: undefined
    },
    class: {
      type: String,
      default: ''
    },
    success: {
      type: Boolean,
      default: undefined
    },
    error: {
      type: Boolean,
      default: undefined
    },
    variant: {
      type: String,
      default: undefined
    },
  }
)

const classList = computed(() => {
  const variant = props.variant ?? (props.error ? 'error' : props.success ? 'success' : undefined)
  return twMerge(config.base, config.variant[variant] ?? '', props.class)
})

const config = {
  base: 'bg-transparent text-sm font-medium text-gray-700 dark:text-gray-200',
  variant: {
    success: 'text-success-500 dark:text-success-400',
    error: 'text-error-500 dark:text-error-400',
  }
}
</script>
<template>
  <p :class="classList">
    <template v-if="props.text">{{ props.text }}</template>
    <template v-else>
      <slot/>
    </template>
  </p>
</template>