<template>
  <div :class="ui.wrapper" v-bind="attrs">
    <div v-if="!isLg && $slots.header" :class="ui.header.top">
      <slot name="header"/>
    </div>

    <div v-if="$slots.left || (isLg && $slots.aside)" :class="ui.left">
      <Aside v-if="isLg && $slots.aside" :class="ui.aside.left">
        <slot name="aside"/>
      </Aside>
      <slot name="left"/>
    </div>

    <div :class="centerClass">
      <div v-if="isLg && $slots.header" :class="ui.header.center">
        <slot name="header"/>
      </div>
      <div v-if="!isLg && $slots.aside" :class="ui.aside.center">
        <slot name="aside"/>
      </div>
      <slot/>
    </div>

    <div v-if="$slots.right" :class="ui.right">
      <slot name="right"/>
    </div>
  </div>
</template>

<script setup lang="ts">
import type {PropType} from 'vue'

const {$viewport} = useNuxtApp()
const isLg = computed(() => $viewport.isGreaterOrEquals('lg'))

const config = {
  wrapper: 'flex flex-col min-h-inherit lg:grid lg:grid-cols-10 lg:gap-8',
  left: 'lg:col-span-2 min-h-inherit',
  center: {
    narrow: 'lg:col-span-6',
    base: 'lg:col-span-8',
    full: 'lg:col-span-10'
  },
  right: 'lg:col-span-2 order-first lg:order-last',
  header: {
    top: '',
    center: ''
  },
  aside: {
    left: 'min-h-inherit',
    center: ''
  },
}

defineOptions({
  inheritAttrs: false
})

const props = defineProps({
  class: {
    type: [String, Object, Array] as PropType<any>,
    default: undefined
  },
  ui: {
    type: Object as PropType<Partial<typeof config>>,
    default: () => ({})
  }
})

const slots = useSlots()
const {ui, attrs} = useUI('page', toRef(props, 'ui'), config, toRef(props, 'class'), true)

const centerClass = computed(() => {
  if ((slots.left || slots.aside) && slots.right) {
    return ui.value.center.narrow
  } else if ((slots.left || slots.aside) || slots.right) {
    return ui.value.center.base
  }

  return ui.value.center.full
})
</script>
