export const useTranslate = () => {
  const $t = useNuxtApp().$i18n.t ?? (key => key)
  return key => $t(key)
}
