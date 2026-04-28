import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', name: 'packages', component: () => import('@/pages/AllPackagesPage.vue') },
    { path: '/update-cleanup', name: 'update-cleanup', component: () => import('@/pages/UpdateCleanupPage.vue') },
    { path: '/taps', name: 'taps', component: () => import('@/pages/TapsPage.vue') },
    { path: '/services', name: 'services', component: () => import('@/pages/ServicesPage.vue') },
    { path: '/brewfile', name: 'brewfile', component: () => import('@/pages/BrewfilePage.vue') },
    { path: '/packages/:type/:name', name: 'package-detail', component: () => import('@/pages/PackageDetailPage.vue') },
    { path: '/install', redirect: '/' },
    { path: '/installed', redirect: '/' },
    { path: '/explore', redirect: '/' },
    { path: '/upgrade', redirect: '/update-cleanup' },
    { path: '/update', redirect: '/update-cleanup' },
    { path: '/maintain', redirect: '/update-cleanup' },
    { path: '/cleanup', redirect: '/update-cleanup' },
    { path: '/bundle', redirect: '/brewfile' },
    { path: '/settings', name: 'settings', component: () => import('@/pages/SettingsPage.vue') },
  ],
})

export default router
