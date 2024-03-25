import { createRouter, createWebHashHistory } from 'vue-router'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import { useCssVar } from '@vueuse/core'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: '/',
			redirect: '/index',
		},
		{
			path: '/index',
			name: 'index',
			component: () => import('../views/index/index.vue'),
			meta: { title: 'Live' },
		},
		{
			path: '/live',
			name: 'live',
			component: () => import('../views/live/live.vue'),
			meta: { title: 'Live' },
		},
		{
			path: '/plive',
			name: 'plive',
			component: () => import('../views/live/plive.vue'),
			meta: { title: 'Live' },
		},
	],
})

router.beforeEach((to, from, next) => {
	document.title = to.meta.title
	NProgress.start()
	next()
})
router.afterEach((to) => {
	document.querySelector('.inner-app').scrollTop = 0
	NProgress.done()
})

export default router
