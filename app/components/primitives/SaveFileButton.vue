<script setup lang="ts">
import type { PropType } from 'vue'

const model = defineModel({
  type: Object,
  default: () => ({text: '', data: null, filename: null}),
})

const props = defineProps({
  label: {
    type: String,
  },
  placeholder: {
    type: String,
  },
  icon: {
    type: String,
  },
  color: {
    type: String as PropType<'white' | 'gray' | 'black' | 'primary' | 'secondary' | 'error' | 'success'>,
  },
  tabindex: {
    type: Number,
    default: undefined
  },
  leadingIcon: {
    type: String,
  },
  trailingIcon: {
    type: String,
  },
  ui: {
    type: Object,
  }
})

const id = useId()

function click() {
  if (import.meta.client && model.value.filename) {
    const a = document.getElementById(id) as HTMLAnchorElement
    if (a) {
      a.href = URL.createObjectURL(new Blob([model.value.data ?? model.value.text], {type: 'application/octet-stream'}))
      a.download = model.value.filename
      a.click()
    }
  }
}
</script>

<template>
  <div class="flex flex-row flex-wrap justify-between items-center gap-x-4">
    <UButton v-bind="props" :disabled="!model.filename" @click="click"/>
    <a :id="id" class="hidden"></a>
  </div>
</template>