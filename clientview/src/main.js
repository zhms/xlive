import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './assets/main.scss'
import 'vant/lib/index.css'
import { urlQuery, formatMoney } from './base'

import { setToastDefaultOptions } from 'vant'
// setToastDefaultOptions({ position: "top", duration: 200000 });

const app = createApp(App)
app.use(router)
app.mount('#app')

// 指令
app.directive('money', (el, binding) => {
	el.innerHTML = formatMoney(binding.value)
})
