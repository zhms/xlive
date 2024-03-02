import Cookies from 'js-cookie'
import Element from 'element-ui'
import request from '@/api/request'
import router from './api/router'
import store from './api/store'
import Vue from 'vue'
import App from './App'
import md5 from 'js-md5'
import moment from 'moment'
import * as filters from './api/filters'
import 'normalize.css/normalize.css'
import './styles/element-variables.scss'
import '@/styles/index.scss'
import './icons'

Vue.prototype.$moment = moment
Vue.prototype.$md5 = md5
Vue.prototype.$get = request.get
Vue.prototype.$post = request.post
Vue.prototype.$patch = request.patch
Vue.prototype.$delete = request.delete
Vue.prototype.$download = request.download

router.beforeEach(async (to, from, next) => {
	document.title = '演示系统'
	if (to.path == '/login') return next()
	if (to.path == '/') return next(`/login`)
	let token = sessionStorage.getItem('x-token')
	if (!token || token == '') return next(`/login`)
	if (to.path == '/dashboard') return next(`/dashboard/index`)
	return next()
})
router.afterEach(() => {})
Vue.use(Element, {
	size: Cookies.get('size') || 'medium',
})
Object.keys(filters).forEach((key) => {
	Vue.filter(key, filters[key])
})
Vue.config.productionTip = false
new Vue({
	el: '#app',
	router,
	store,
	render: (h) => h(App),
})
