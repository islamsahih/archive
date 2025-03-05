<script setup lang="ts">
const appConfig = useAppConfig()
const model = defineModel({
  type: String,
  default: ''
})
const props = defineProps({
  placeholder: String,
  loading: Boolean,
  disabled: Boolean,
})
const emit = defineEmits(["search", "reset"])

function reset() {
  model.value = ''
  emit('reset')
}
function search() {
  if (model.value) emit('search')
  else reset()
}

</script>

<template>
  <UInput
    v-model="model"
    color="white"
    :placeholder="placeholder ?? $t('search.placeholder')"
    :icon="model === '' ? appConfig.ui.icons.search : undefined"
    :loading-icon="appConfig.ui.icons.loading"
    :loading="loading"
    :disabled="disabled"
    autocomplete="off"
    @keyup.enter="search"
    :ui="{ icon: { leading: { pointer: '' }, trailing: { pointer: '' } } }"
  >
    <!--    <template v-if="!loading" #leading>-->
    <!--      <UButton-->
    <!--        color="gray"-->
    <!--        variant="link"-->
    <!--        icon="i-heroicons-magnifying-glass-20-solid"-->
    <!--        :padded="false"-->
    <!--        @click="emit('search')"-->
    <!--      />-->
    <!--    </template>-->
    <template v-if="$slots.trailing || (!loading && model !== '')" #trailing>
      <div class="flex flex-row gap-1">
        <template v-if="!loading && model !== ''">
          <UButton
            color="gray"
            variant="link"
            :icon="appConfig.ui.icons.close"
            :padded="false"
            @click="reset"
          />
          <UButton
            :icon="appConfig.ui.icons.search"
            color="primary"
            variant="link"
            :padded="false"
            @click="search"
          />
        </template>
        <slot name="trailing" />
      </div>
    </template>
  </UInput>
</template>
