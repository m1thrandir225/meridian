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

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // Chat routes
    {
      path: '/',
      name: 'home',
      component: ChatView,
    },
    {
      path: '/channel/:id',
      name: 'channel',
      component: ChatView,
    },
    {
      path: '/bot-registration',
      name: 'bot-registration',
      component: BotRegistrationView,
    },
    // Settings routes
    {
      path: '/settings/profile',
      name: 'settings-profile',
      component: ProfileView,
    },
    {
      path: '/settings/password',
      name: 'settings-password',
      component: PasswordView,
    },
    {
      path: '/settings/appearance',
      name: 'settings-appearance',
      component: AppearanceView,
    },
    // Redirect /settings to profile
    {
      path: '/settings',
      redirect: '/settings/profile',
    },
    // Auth routes
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterView,
    },
    {
      path: '/forgot-password',
      name: 'forgot-password',
      component: ForgotPasswordView,
    },
  ],
})

export default router
