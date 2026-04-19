import { createI18n } from 'vue-i18n'
import en from './en.json'
import zh from './zh.json'

const preferred = localStorage.getItem('gobrew-lang')
const savedLang = preferred ?? (navigator.language.startsWith('zh') ? 'zh' : 'en')

const i18n = createI18n({
  legacy: false,
  locale: savedLang,
  fallbackLocale: 'en',
  messages: { en, zh },
})

export default i18n
