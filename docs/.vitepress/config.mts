import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  lang: "en-US",
  title: "Vite Go",
  description: "Connect your Go application with Vite-powered frontends",
  base: "/vite/",
  head: [
    ['link', { rel: 'apple-touch-icon', sizes: '180x180', href: '/vite/apple-touch-icon.png' }],
    ['link', { rel: 'icon', type: 'image/png', sizes: '32x32', href: '/vite/favicon-32x32.png' }],
    ['link', { rel: 'icon', type: 'image/png', sizes: '16x16', href: '/vite/favicon-16x16.png' }]
  ],
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Documentation', link: '/guide/getting-started' },
      { text: 'Examples', link: '/guide/examples' }
    ],
    logo: '/images/logo.webp',
    sidebar: [
      {
        text: 'Introduction',
        items: [
          { text: 'Getting Started', link: '/guide/getting-started' }
        ]
      },
      {
        text: 'Usage',
        items: [
          { text: 'Helper Function', link: '/guide/usage#option-1-helper-function' },
          { text: 'HTTP Handler', link: '/guide/usage#option-2-http-handler' },
          { text: 'Configuration', link: '/guide/usage#configuration-options' }
        ]
      },
      {
        text: 'Examples',
        items: [
          { text: 'Simple Helper', link: '/guide/examples#simple-helper-function' },
          { text: 'Basic Handler', link: '/guide/examples#basic-with-handler' },
          { text: 'Multi Page App', link: '/guide/examples#multi-page-application' },
          { text: 'Multiple Vite Instances App', link: '/guide/examples#multiple-vite-instances-application' },
          { text: 'Template Registration', link: '/guide/examples#template-registration' },
          { text: 'Inertia.js', link: '/guide/examples#inertia-js-integration' },
          { text: 'Router App', link: '/guide/examples#router-application' }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/olivere/vite' }
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2023-present'
    },

    search: {
      provider: 'local'
    },
  }
})
