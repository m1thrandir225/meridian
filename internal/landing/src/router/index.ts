import { usePostHog } from '@/composables/usePostHog'
import HomeView from '@/views/HomeView.vue'
import { createWebHistory, createRouter, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: HomeView,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

usePostHog()

export default router
