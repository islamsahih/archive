export default defineI18nConfig(() => ({
  legacy: false,
  locale: 'ru',
  messages: {
    ru: {
      common: {
        uncategorized: 'Без категории',
        references: 'Ссылки:',
        relatedPages: 'Связанные страницы:',
        editNotification: {
          title: 'Раздел редактируется',
          text: 'На настоящий момент на сайте представлены не все материалы раздела. Полный список уроков из этой главы вы можете найти в нашем {{#book}}\u003ca href=\"{{link}}\" target=\"_blank\"\u003e\u003cspan class=\"inline-flex items-baseline ml-1 gap-x-1\"\u003e\u003cimg class=\"not-prose w-5 h-5 self-center\" src=\"/icon/{{icon}}.svg\" alt=\"{{title}}\" /\u003eTelegram\u003c/span\u003e\u003c/a\u003e{{/book}}.'
        }
      },
      search: {
        title: 'Поиск',
        placeholder: 'Поиск...',
        progress: 'Идет поиск...',
        button: 'Поиск',
        found: 'Найдено',
        notFound: 'Не найдено',
        options: 'Опции',
        inSection: 'Поиск в разделе:',
        allSections: 'Все разделы',
        in: 'Область поиска:',
        inOptions: {
          title: 'Заголовки',
          description: 'Описания',
          text: 'Текст',
          tags: 'Теги',
        },
        results: 'Результаты поиска',
        resultsPerPage: 'Результатов на странице:',
        orderBy: 'Сортировать',
        orderByOptions: {
          score: 'По совпадению',
          index: 'По порядку',
          title: 'По заголовку',
          date: 'По дате',
        },
        show: 'Показывать:',
        showOptions: {
          sections: 'Разделы',
          matches: 'Совпадения',
          dates: 'Даты',
          audios: 'Аудио',
          videos: 'Видео',
        },
      }
    },
  },
  missing: (locale, key) => {
    return key
  },
}))