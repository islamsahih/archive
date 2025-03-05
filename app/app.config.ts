export default defineAppConfig({
  seo: {
    title: 'Islam Sahih',
    initialYear: 2023,
    social: [
      {
        title: 'Задать вопрос по Исламу, сообщить админам об ошибке на сайте',
        icon: 'i-logos-telegram',
        link: 'https://t.me/sahih_islam',
        admin: true
      },
      {
        title: 'Ответы на вопросы',
        icon: 'i-logos-telegram',
        link: 'https://t.me/islamic_storage'
      },
      {
        title: 'Канал с нашими видео',
        icon: 'i-logos-youtube-icon',
        link: 'https://www.youtube.com/@Islam_sahih'
      },
    ],
  },
  menu: {
    links: [
      '/islam_kamil',
      '/islam_kamil_hadiths',
      '/qa'
    ],
    footerLinks: [
      '/about', '/contact', '/terms',
    ]
  },
  header: {
    links: [
      {
        '/islam_kamil': [
          '/islam_kamil/part_1',
          '/islam_kamil/part_2',
          '/islam_kamil/part_3'
        ]
      },
      {
        '/islam_kamil_hadiths': [
          '/islam_kamil_hadiths/part_1',
          '/islam_kamil_hadiths/part_2',
          '/islam_kamil_hadiths/part_3'
        ]
      },
      '/qa',
      '/search'
    ]
  },
  aside: {
    links: [
      {
        '/islam_kamil': [
          '/islam_kamil/part_1',
          '/islam_kamil/part_2',
          {'/islam_kamil/part_3': ['/islam_kamil/part_3/chapter_1', '/islam_kamil/part_3/chapter_2']},
        ]
      },
      {'/islam_kamil_hadiths': ['/islam_kamil_hadiths/part_1', '/islam_kamil_hadiths/part_2', '/islam_kamil_hadiths/part_3']},
      '/qa'
    ],
  },
  footer: {
    links: [
      '/about', '/contact', '/terms',
    ],
  },
  ui: {
    primary: 'sky',
    gray: 'cool',
    variables: {
      root: {
        content: {
          width: '1280px;'
        },
        header: {
          width: 'var(--content-width)',
          height: '4rem',
        },
        footer: {
          height: '4rem',
        },
      },
      light: {
        background: '255 255 255',
        foreground: 'var(--color-gray-700)'
      },
      dark: {
        background: 'var(--color-gray-900)',
        foreground: 'var(--color-gray-200)'
      },
    },
    icons: {
      dark: 'i-heroicons-moon-20-solid',
      light: 'i-heroicons-sun-20-solid',
      system: 'i-heroicons-computer-desktop-20-solid',
      search: 'i-ion-search',
      external: 'i-heroicons-arrow-up-right-20-solid',
      chevron: 'i-heroicons-chevron-down-20-solid',
      hash: 'i-heroicons-hashtag-20-solid',
      menu: 'i-heroicons-bars-3-20-solid',
      close: 'i-heroicons-x-mark-20-solid',
      check: 'i-heroicons-check-circle-20-solid',
      loading: 'i-mdi-progress-helper',
      sort_asc: 'i-mdi-sort-ascending',
      sort_desc: 'i-mdi-sort-descending',

      chevron_up: 'i-heroicons-chevron-up-20-solid',
      chevron_down: 'i-heroicons-chevron-down-20-solid',
      chevron_right: 'i-heroicons-chevron-right-20-solid',
      chevron_left: 'i-heroicons-chevron-left-20-solid',

      arrow_left: 'i-heroicons-arrow-left-20-solid',
      arrow_right: 'i-heroicons-arrow-right-20-solid',

      bookmark: 'i-lets-icons-bookmark-duotone',
      book_open_duo: 'i-lets-icons-book-open-duotone-line',
      book_open_alt_duo: 'i-lets-icons-book-open-alt-duotone',
      book_ltr: 'i-ooui-book-ltr',
      book_sharp: 'i-ion-book-sharp',

      journal_ltr: 'i-ooui-journal-ltr',
      ios_journal: 'i-ion-ios-journal',

      question: 'i-lets-icons-question',
      question_duo: 'i-lets-icons-question-duotone-line',

      chatboxes: 'i-ion-chatboxes',
    },
    main: {
      wrapper: 'min-h-[calc(100vh-var(--header-height))] md:min-h-[calc(100vh-var(--header-height)-var(--footer-height))]',
    },
    breadcrumb: {
      wrapper: 'not-prose',
      ol: 'flex-wrap'
    },
    breadcrumbInactive: {
      ol: 'gap-x-0',
      li: 'gap-x-0',
      active: 'text-gray-500 dark:text-gray-400'
    },
    content: {
      surround: {
        link: {
          icon: {
            wrapper: 'ring-0'
          }
        }
      }
    }
  }
})