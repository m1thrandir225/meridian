import './assets/index.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import { createHead } from '@unhead/vue/client'
import App from './App.vue'

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

const head = createHead()
const app = createApp(App)

app.use(pinia)
app.use(head)
app.mount('#app')
