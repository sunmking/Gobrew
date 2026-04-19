import { createRouter, createWebHashHistory } from 'vue-router'
import { DOMAIN_ENTRY_PATHS, resolveDomainFromPath } from '@/config/domains'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', name: 'home', component: () => import('@/pages/HomePage.vue'), meta: { domain: resolveDomainFromPath('/') } },
    {
      path: DOMAIN_ENTRY_PATHS.install,
      name: 'install',
      component: () => import('@/pages/InstallPage.vue'),
      meta: { domain: resolveDomainFromPath(DOMAIN_ENTRY_PATHS.install) },
    },
    {
      path: DOMAIN_ENTRY_PATHS.upgrade,
      name: 'upgrade',
      component: () => import('@/pages/UpdatePage.vue'),
      meta: { domain: resolveDomainFromPath(DOMAIN_ENTRY_PATHS.upgrade) },
    },
    {
      path: DOMAIN_ENTRY_PATHS.maintain,
      name: 'maintain',
      component: () => import('@/pages/MaintainPage.vue'),
      meta: { domain: resolveDomainFromPath(DOMAIN_ENTRY_PATHS.maintain) },
    },
    {
      path: DOMAIN_ENTRY_PATHS.services,
      name: 'services',
      component: () => import('@/pages/ServicesPage.vue'),
      meta: { domain: resolveDomainFromPath(DOMAIN_ENTRY_PATHS.services) },
    },
    { path: '/installed', name: 'installed', component: () => import('@/pages/InstalledPage.vue'), meta: { domain: resolveDomainFromPath('/installed') } },
    { path: '/explore', name: 'explore', component: () => import('@/pages/ExplorePage.vue'), meta: { domain: resolveDomainFromPath('/explore') } },
    { path: '/update', name: 'update-legacy', redirect: DOMAIN_ENTRY_PATHS.upgrade },
    { path: '/services', name: 'services', component: () => import('@/pages/ServicesPage.vue'), meta: { domain: resolveDomainFromPath('/services') } },
    { path: '/taps', name: 'taps', component: () => import('@/pages/TapsPage.vue'), meta: { domain: resolveDomainFromPath('/taps') } },
    { path: '/cleanup', name: 'cleanup', component: () => import('@/pages/CleanupPage.vue'), meta: { domain: resolveDomainFromPath('/cleanup') } },
    { path: '/bundle', name: 'bundle', component: () => import('@/pages/BundlePage.vue'), meta: { domain: resolveDomainFromPath('/bundle') } },
    { path: '/settings', name: 'settings', component: () => import('@/pages/SettingsPage.vue'), meta: { domain: resolveDomainFromPath('/settings') } },
  ],
})

export default router
