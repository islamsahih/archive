import type {Config} from 'tailwindcss'
import colors from 'tailwindcss/colors'

export default <Partial<Config>>{
  content: [
    './pages/**/*.{vue,js,ts}',
    './components/**/*.{vue,js,ts}',
    './layouts/**/*.{vue,js,ts}',
    './plugins/**/*.{js,ts}',
    './nuxt.config.{js,ts}',
  ],
  theme: {
    extend: {
      fontFamily: {
        // sans: ['ui-sans-serif', 'system-ui', 'sans-serif', "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji"]
         sans: ['app-sans', 'app-arabic', 'ui-sans-serif', 'system-ui', 'sans-serif', "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji"]
      },
      colors: {
        success: colors.green,
        error: colors.red,

        background: 'rgb(var(--ui-background) / <alpha-value>)',
        foreground: 'rgb(var(--ui-foreground) / <alpha-value>)'
      },
      minHeight: {
        inherit: 'inherit'
      },
      typography: (theme: any) => {
        return {
          DEFAULT: {
            css: {
              'h1, h2, h3, h4': {
                fontWeight: theme('fontWeight.bold'),
                'scroll-margin-top': 'var(--scroll-mt)'
              },
              'h1 a, h2 a, h3 a, h4 a': {
                borderBottom: 'none !important',
                color: 'inherit',
                fontWeight: 'inherit'
              },
              a: {
                fontWeight: theme('fontWeight.medium'),
                textDecoration: 'none',
              },
              'a:hover': {
                borderColor: 'var(--tw-prose-links)'
              },
              'a:has(> code)': {
                borderColor: 'transparent !important'
              },
              'a code': {
                color: 'var(--tw-prose-code)',
                border: '1px dashed var(--tw-prose-pre-border)'
              },
              'a:hover code': {
                color: 'var(--tw-prose-links)',
                borderColor: 'var(--tw-prose-links)'
              },
              pre: {
                borderRadius: '0.375rem',
                border: '1px solid var(--tw-prose-pre-border)',
                color: 'var(--tw-prose-pre-code) !important',
                backgroundColor: 'var(--tw-prose-pre-bg) !important',
                whiteSpace: 'pre-wrap',
                wordBreak: 'break-word'
              },
              code: {
                backgroundColor: 'var(--tw-prose-pre-bg)',
                padding: '0 0.375rem',
                display: 'inline-block',
                borderRadius: '0.375rem',
                border: '1px solid var(--tw-prose-pre-border)'
              },
              'code::before': {
                content: ''
              },
              'code::after': {
                content: ''
              },
              'blockquote p:first-of-type::before': {
                content: ''
              },
              'blockquote p:last-of-type::after': {
                content: ''
              },
              'input[type="checkbox"]': {
                color: 'rgb(var(--color-primary-500))',
                borderRadius: theme('borderRadius.DEFAULT'),
                borderColor: 'rgb(var(--color-gray-300))',
                height: theme('spacing.4'),
                width: theme('spacing.4'),
                marginTop: '-3.5px !important',
                marginBottom: '0 !important',
                '&:focus': {
                  '--tw-ring-offset-width': 0
                }
              },
              'input[type="checkbox"]:checked': {
                borderColor: 'rgb(var(--color-primary-500))'
              },
              'input[type="checkbox"]:disabled': {
                opacity: 0.5,
                cursor: 'not-allowed'
              },
              'p[lang="ar"]': {
                fontWeight: '700',
              },
              // 'p': {
              //   textIndent: '2rem'
              // },
              'ul.contains-task-list': {
                marginLeft: '-1.625em'
              },
              'ol': {
                paddingInlineStart: '0px',
                'list-style-position': 'inside'
              },
              'ol > li': {
                // textIndent: '2rem',
                paddingInlineStart: '0px',
              },
              'ol > li::marker': {
                color: 'rgb(var(--tw-prose-body))',
              },
              'ol.marker-bold > li::marker': {
                fontWeight: '700',
              },
              'ul ul': {
                paddingLeft: theme('padding.6')
              },
              'ul ol': {
                paddingLeft: theme('padding.6')
              },
              'ul > li.task-list-item': {
                paddingLeft: '0 !important'
              },
              'ul > li.task-list-item input': {
                marginRight: '7px'
              },
              'ul > li.task-list-item > ul.contains-task-list': {
                marginLeft: 'initial'
              },
              'ul > li.task-list-item a': {
                marginBottom: 0
              },
              'ul > li.task-list-item::marker': {
                content: 'none'
              },
              'ul > li > p': {
                margin: 0
              },
              'ul > li > span.issue-badge, p > span.issue-badge': {
                verticalAlign: 'text-top',
                margin: '0 !important'
              },
              'ul > li > button': {
                verticalAlign: 'baseline !important'
              },
              'blockquote': {
                fontStyle: 'normal',
              },
              'blockquote h6': {
                fontStyle: 'normal',
              },
              'blockquote table': {
                marginTop: '1.6em',
                marginBottom: '1.6em',
                fontSize: theme('fontSize.base'),
                fontWeight: theme('fontWeight.normal'),
                fontStyle: 'normal'
              },
              'blockquote thead': {
                display: 'none'
              },
              'blockquote tbody td': {
                paddingTop: '1rem',
                paddingBottom: '1rem',
                paddingInlineStart: '1rem',
                paddingInlineEnd: '1rem'
              },
              'blockquote blockquote': {
                borderLeft: 'none',
                backgroundColor: 'rgb(var(--color-gray-100))',
                paddingTop: '1rem',
                paddingBottom: '1rem',
                paddingRight: '1rem',
                paddingLeft: '1rem',
              },
              '.dark blockquote blockquote': {
                backgroundColor: 'rgb(var(--color-gray-800))',
              },
              'blockquote blockquote p': {
                textAlign: 'justify',
                marginTop: '0',
                marginBottom: '0'
              },
              'blockquote blockquote p:is(:has(+ h6))': {
                marginBottom: '0',
              },
              'blockquote blockquote h6': {
                // marginTop: '0.25rem',
                // marginBottom: '0.25rem',
                // fontStyle: 'italic',
                textAlign: 'right',
                fontSize: theme('fontSize.sm'),
                fontWeight: theme('fontWeight.normal'),
              },
              table: {
                display: 'block',
                overflowX: 'auto',
              },
              'table code': {
                display: 'inline-flex'
              },
              'searchmark': {
                backgroundColor: 'rgb(var(--color-primary-200))',
                fontWeight: theme('fontWeight.bold'),
              }
            }
          },
          primary: {
            css: {
              '--tw-prose-body': 'rgb(var(--color-gray-700))',
              '--tw-prose-headings': 'rgb(var(--color-gray-900))',
              '--tw-prose-lead': 'rgb(var(--color-gray-600))',
              '--tw-prose-links': 'rgb(var(--color-primary-500))',
              '--tw-prose-bold': 'rgb(var(--color-gray-900))',
              '--tw-prose-counters': 'rgb(var(--color-gray-500))',
              '--tw-prose-bullets': 'rgb(var(--color-gray-300))',
              '--tw-prose-hr': 'rgb(var(--color-gray-200))',
              '--tw-prose-quotes': 'rgb(var(--color-gray-900))',
              '--tw-prose-quote-borders': 'rgb(var(--color-gray-200))',
              '--tw-prose-captions': 'rgb(var(--color-gray-500))',
              '--tw-prose-code': 'rgb(var(--color-gray-900))',
              '--tw-prose-pre-code': 'rgb(var(--color-gray-900))',
              '--tw-prose-pre-bg': 'rgb(var(--color-gray-50))',
              '--tw-prose-pre-border': 'rgb(var(--color-gray-200))',
              '--tw-prose-th-borders': 'rgb(var(--color-gray-300))',
              '--tw-prose-td-borders': 'rgb(var(--color-gray-200))',
              '--tw-prose-invert-body': 'rgb(var(--color-gray-200))',
              '--tw-prose-invert-headings': theme('colors.white'),
              '--tw-prose-invert-lead': 'rgb(var(--color-gray-400))',
              '--tw-prose-invert-links': 'rgb(var(--color-primary-400))',
              '--tw-prose-invert-bold': theme('colors.white'),
              '--tw-prose-invert-counters': 'rgb(var(--color-gray-400))',
              '--tw-prose-invert-bullets': 'rgb(var(--color-gray-600))',
              '--tw-prose-invert-hr': 'rgb(var(--color-gray-800))',
              '--tw-prose-invert-quotes': 'rgb(var(--color-gray-100))',
              '--tw-prose-invert-quote-borders': 'rgb(var(--color-gray-700))',
              '--tw-prose-invert-captions': 'rgb(var(--color-gray-400))',
              '--tw-prose-invert-code': theme('colors.white'),
              '--tw-prose-invert-pre-code': theme('colors.white'),
              '--tw-prose-invert-pre-bg': 'rgb(var(--color-gray-800))',
              '--tw-prose-invert-pre-border': 'rgb(var(--color-gray-700))',
              '--tw-prose-invert-th-borders': 'rgb(var(--color-gray-700))',
              '--tw-prose-invert-td-borders': 'rgb(var(--color-gray-800))'
            }
          },
          invert: {
            css: {
              '--tw-prose-pre-border': 'var(--tw-prose-invert-pre-border)',
              'input[type="checkbox"]': {
                backgroundColor: 'rgb(var(--color-gray-800))',
                borderColor: 'rgb(var(--color-gray-700))'
              },
              'input[type="checkbox"]:checked': {
                backgroundColor: 'rgb(var(--color-primary-400))',
                borderColor: 'rgb(var(--color-primary-400))'
              }
            }
          }
        }
      }
    }
  },
  safelist: [
    {
      pattern: /ring-(primary|success|error)-(500|400)/,
      variants: [
        'hover', 'dark:hover',
        'hover-group', 'dark:hover-group',
        'focus-within', 'dark:focus-within',
        'focus-visible', 'dark:focus-visible',
      ],
    },
  ],
}

