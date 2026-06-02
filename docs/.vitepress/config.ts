import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'tfparams',
  description: 'Terraform parameter sheet generator',
  base: '/tfparams/',
  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/tfparams/favicon.svg' }],
  ],
  themeConfig: {
    // base 配下では themeConfig.logo に base が自動付与されないことがあるため明示する（VitePress issue #2981）
    logo: '/tfparams/logo.svg',
    nav: [
      { text: 'Guide', link: '/guide/getting-started' },
      { text: 'Reference', link: '/reference/cli' },
      { text: 'Changelog', link: '/changelog' },
    ],
    sidebar: {
      '/guide/': [
        {
          text: 'Guide',
          items: [
            { text: 'Getting Started', link: '/guide/getting-started' },
            { text: 'Installation', link: '/guide/installation' },
            { text: 'Configuration', link: '/guide/configuration' },
          ],
        },
      ],
      '/reference/': [
        {
          text: 'Reference',
          items: [
            { text: 'CLI', link: '/reference/cli' },
            { text: 'Config File', link: '/reference/config-file' },
          ],
        },
      ],
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/tfparam/tfparams' },
    ],
    footer: {
      message: 'Released under the MIT License.',
    },
  },
})
