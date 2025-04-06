// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router';
import MobileOverview from '@/components/MobileOverview.vue';
import MobileHostSettings from '@/components/MobileHostSettings.vue';
import MobileCarDetails from '@/components/MobileCarDetails.vue';

const routes = [
  {
    path: '/',
    name: 'Overview',
    component: MobileOverview
  },
  {
    path: '/host',
    name: 'HostSettings',
    component: MobileHostSettings
  },
  {
    path: '/car/:id',
    name: 'CarDetails',
    component: MobileCarDetails,
    props: true
  },
  {
    path: '/:catchAll(.*)',
    redirect: '/'
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;
