import { createRouter, createWebHistory } from 'vue-router'
import ChatView from '../views/ChatView.vue'
import LoginView from '@/views/LoginView.vue'
import RegisterView from '@/views/RegisterView.vue'
import ForgotPasswordView from '@/views/ForgotPasswordView.vue'
import BotRegistrationView from '@/views/BotRegistrationView.vue'
// Settings pages
import ProfileView from '@/views/settings/ProfileView.vue'
import PasswordView from '@/views/settings/PasswordView.vue'
import AppearanceView from '@/views/settings/AppearanceView.vue'
import { useAuthStore } from '@/stores/auth'
import BotManagementView from '@/views/BotManagementView.vue'
import InviteAcceptView from '@/views/InviteAcceptView.vue'
import HomeView from '@/views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // Chat routes
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { requiresAuth: true },
    },
    {
      path: '/channel/:id',
      name: 'channel',
      component: ChatView,
      meta: { requiresAuth: true },
    },
    {
      path: '/bot-registration',
      name: 'bot-registration',
      component: BotRegistrationView,
      meta: { requiresAuth: true },
    },
    {
      path: '/bot-management',
      name: 'bot-management',
      component: BotManagementView,
      meta: { requiresAuth: true },
    },
    //Invite routes
    {
      path: '/invites/:inviteCode',
      name: 'invite-accept',
      component: InviteAcceptView,
      meta: { requiresAuth: true },
    },
    // Settings routes
    {
      path: '/settings/profile',
      name: 'settings-profile',
      component: ProfileView,
      meta: { requiresAuth: true },
    },
    {
      path: '/settings/password',
      name: 'settings-password',
      component: PasswordView,
      meta: { requiresAuth: true },
    },
    {
      path: '/settings/appearance',
      name: 'settings-appearance',
      component: AppearanceView,
      meta: { requiresAuth: true },
    },
    // Redirect /settings to profile
    {
      path: '/settings',
      redirect: '/settings/profile',
      meta: { requiresAuth: true },
    },
    // Auth routes
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { requiresAuth: false },
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterView,
      meta: { requiresAuth: false },
    },
    {
      path: '/forgot-password',
      name: 'forgot-password',
      component: ForgotPasswordView,
      meta: { requiresAuth: false },
    },
  ],
})

router.beforeEach((to, from, next) => {
  const isAuthenticated = useAuthStore().checkAuth()

  if (to.meta.requiresAuth && !isAuthenticated) {
    const redirect = to.fullPath

    next({ name: 'login', query: { redirect } })
  } else {
    next()
  }
})

export default router
