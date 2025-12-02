import { createApp } from 'vue'
import App from './App.vue'
// import router from './router'

// 引入 element plus
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import './assets/css/base.css'
import './assets/css/code.css'
import './assets/css/preview.css'
import 'highlight.js/styles/github.css'

const app = createApp(App)

// app.use(router)
app.use(ElementPlus)

app.mount('#app')
