import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import i18n from './locales'
import './style.css'

document.documentElement.dataset.theme = 'dark'
document.documentElement.classList.toggle('dark', true)

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(i18n)
app.mount('#app')
