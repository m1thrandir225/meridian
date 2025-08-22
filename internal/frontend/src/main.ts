import './assets/index.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import { VueQueryPlugin } from '@tanstack/vue-query'
import App from './App.vue'
import router from './router'
import websocketService from '@/services/websocket.service'
import { useAuthStore } from '@/stores/auth'

const app = createApp(App)
const pinia = createPinia()

pinia.use(piniaPluginPersistedstate)

app.use(pinia)
app.use(router)
app.use(VueQueryPlugin)

// Initialize WebSocket connection when app starts
const authStore = useAuthStore()
if (authStore.checkAuth()) {
  websocketService.connect()
}

app.mount('#app')
