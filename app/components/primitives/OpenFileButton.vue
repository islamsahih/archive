<script setup lang="ts">
import type { PropType } from 'vue'

const model = defineModel({
  type: String,
  default: '',
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
  loadingVariant: {
    type: String, // 'default', 'outside', 'disabled'
    default: 'default',
  },
  loadingIcon: {
    type: String,
  },
  loadingText: {
    type: String,
    default: '',
  },
  change: {
    type: Function,
  },
  onload: {
    type: Function,
  },
  onerror: {
    type: Function,
  },
  onabort: {
    type: Function,
  },
  onloadstart: {
    type: Function,
  },
  onloadend: {
    type: Function,
  },
  disabled: {
    type: Boolean,
  },
  ui: {
    type: Object,
  }
})

const id = useId()
const loading = ref(false)
const loadingDefault = computed(() => props.loadingVariant === 'default' && loading.value)
const loadingOutside = computed(() => props.loadingVariant === 'outside' && loading.value)

function change(event: any) {
  if (props.change) {
    props.change(event)
  } else {
    if (import.meta.client && event?.target?.files?.length === 1) {
      let file = event.target.files[0]
      let reader = new FileReader()
      reader.onload = (event) => {
        model.value = file.name
        if (props.onload) props.onload(event)
      }
      reader.onerror = (event) => {
        model.value = 'The requested file could not be read'
        if (props.onerror) props.onerror(event)
      }
      reader.onabort = (event) => {
        model.value = 'The requested file could not be read'
        if (props.onabort) props.onabort(event)
      }
      reader.onloadstart = (event) => {
        loading.value = props.loadingVariant != 'disabled' && file.size >= 10485760
        if (props.onloadstart) props.onloadstart(event)
      }
      reader.onloadend = (event) => {
        loading.value = false
        if (props.onloadend) props.onloadend(event)
      }
      reader.readAsArrayBuffer(file)
    }
  }
}

function click() {
  if (import.meta.client) {
    document.getElementById(id)?.click()
  }
}
</script>

<template>
  <div class="flex flex-row flex-wrap justify-between items-center gap-x-4">
    <input
      type="file"
      :id="id"
      :multiple="false"
      class="hidden"
      @change="change"/>
    <UButton v-bind="props" :loading="loadingDefault" @click="click"/>
    <Spinner v-if="loadingOutside"
             :icon="loadingIcon"
             :text="loadingText"
    />
    <p v-else-if="model?.length" class="truncate max-w-48 text-sm font-medium">{{ model }}</p>
    <p v-else-if="placeholder?.length" class="text-sm font-medium">{{ placeholder }}</p>
  </div>
</template>