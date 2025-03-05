<script setup lang="ts">
import {twMerge} from "tailwind-merge";

const props = defineProps({
    label: {
      type: String,
      default: undefined
    },
    class: {
      type: String,
      default: ''
    },
    labelClass: {
      type: String,
      default: ''
    },
    row: {
      type: Boolean,
      default: undefined,
    },
    col: {
      type: Boolean,
      default: true,
    },
    variant: {
      type: String,
      default: undefined
    }
  }
)

const variant = computed(() => {
  const variant = props.variant ?? props.row ? 'row' : props.col ? 'col' : undefined
  return config.variant[variant as keyof typeof config.variant] ?? {base: '', label: ''}
})

const classList = computed(() => {
  return twMerge(config.base, variant.value.base, props.class)
})

const labelClassList = computed(() => {
  return twMerge(config.label, variant.value.label, props.labelClass)
})

const config = {
  base: '',
  label: '',
  variant: {
    col: {
      base: '',
      label: 'mb-1',
    },
    row: {
      base: 'flex flex-row gap-x-2.5',
      label: 'self-center',
    }
  }
}
</script>
<template>
  <div :class="classList">
    <Label :class="labelClassList" :text="props.label" />
    <slot/>
  </div>
</template>