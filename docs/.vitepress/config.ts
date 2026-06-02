import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'tfparams',
  description: 'Terraform parameter sheet generator',
  base: '/tfparams/',
  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/tfparams/favicon.svg' }],
  ],
  themeConfig: {
    // VitePress prepends `base` to themeConfig.logo automatically — pass a
    // root-relative path WITHOUT the base prefix, or it doubles to
    // /tfparams/tfparams/logo.svg and 404s.
    logo: '/logo.svg',
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
