import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

const config: Config = {
  title: 'DLC Desktop Application',
  tagline: 'Lesson videos for every classroom, even without internet.',
  favicon: 'img/dlc.png',

  // Set the production url of your site here
  url: 'https://deerwalkrnd.github.io/',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/dlc-desktop-app',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'deerwalkrnd', // Usually your GitHub org/user name.
  projectName: 'dlc-desktop-app', // Usually your repo name.

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'throw',

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/deerwalkrnd/dlc-desktop-app/tree/main/docs/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    image: 'img/dlc.png',
    navbar: {
      title: 'DLC Desktop App',
      logo: {
        alt: 'Deerwalk Learning Center',
        src: 'img/dlc.png',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'guideSidebar',
          position: 'left',
          label: 'Guide',
        },
        {
          href: 'https://github.com/deerwalkrnd/dlc-desktop-app/releases',
          label: 'Download',
          position: 'right',
        },
        {
          href: 'https://github.com/deerwalkrnd/dlc-desktop-app',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Guide',
          items: [
            {label: 'Getting started', to: '/docs/intro'},
            {label: 'Adding videos', to: '/docs/adding-videos'},
            {label: 'Troubleshooting', to: '/docs/troubleshooting'},
          ],
        },
        {
          title: 'Project',
          items: [
            {label: 'Releases', href: 'https://github.com/deerwalkrnd/dlc-desktop-app/releases'},
            {label: 'Developer guide', to: '/docs/developers'},
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Deerwalk Learning Center.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
