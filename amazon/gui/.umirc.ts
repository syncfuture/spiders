import { defineConfig } from 'umi';

export default defineConfig({
  title: "Amazon Scraper v1.0.0",
  history: { type: 'hash' },
  nodeModulesTransform: {
    type: 'none',
  },
  routes: [
    {
      path: '/',
      component: '@/_layout',
      routes: [
        { path: '/', component: '@/pages/index' },
        { path: '/reviews', component: '@/pages/reviews' },
        { path: '/items', component: '@/pages/items' },
      ]
    },
  ],
  fastRefresh: {},
  antd: {
    dark: false,
    compact: true,
  },
  dva: {
    immer: true,
    hmr: false,
  },
});
