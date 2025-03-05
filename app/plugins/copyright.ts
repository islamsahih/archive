export default defineNuxtPlugin(() => {
  const seo: any = useAppConfig().seo
  const currentYear = new Date().getFullYear()
  const dates = (seo.initialYear && seo.initialYear !== currentYear) ?
    `${seo.initialYear}-${currentYear}` : `${currentYear}`
  return {
    provide: {
      copyright: `Copyright © ${dates} ${seo.author ?? seo.title}`,
    }
  }
})
