import Vue from 'vue'
import Router from 'vue-router'
import Layout from '@/layout/index.vue'
Vue.use(Router)
let routers = [
	{ path: '/login', component: () => import('@/layout/login') },
	{ path: '/404', component: () => import('@/layout/404') },
	{ path: '/401', component: () => import('@/layout/401') },
	{
		path: '/dashboard',
		component: Layout,
		name: '系统首页',
		meta: { title: '系统首页', icon: 'el-icon-s-home' },
		children: [
			{
				path: 'index',
				component: () => import('../../views/Dashboard/index'),
				meta: { title: '系统首页', name: '系统首页', icon: 'el-icon-s-home' },
			},
		],
	},
	{
		path: '/user',
		component: Layout,
		meta: { title: '会员管理', icon: 'el-icon-s-tools' },
		children: [
			{
				path: 'list',
				component: () => import('../../views/User/index'),
				meta: { title: '会员列表', icon: 'el-icon-setting' },
			},
		],
	},
	{
		path: '/system',
		component: Layout,
		meta: { title: '系统管理', icon: 'el-icon-s-tools' },
		children: [
			{
				path: 'setting',
				component: () => import('../../views/System/Setting/index'),
				meta: { title: '系统设置', icon: 'el-icon-setting' },
			},
			{
				path: 'role',
				component: () => import('../../views/System/Role/index'),
				meta: { title: '角色管理', icon: 'el-icon-set-up' },
			},
			{
				path: 'account',
				component: () => import('../../views/System/User/index'),
				meta: { title: '账号管理', icon: 'el-icon-office-building' },
			},
			{
				path: 'login',
				component: () => import('../../views/System/LoginLog/index'),
				meta: { title: '登录日志', icon: 'el-icon-notebook-1' },
			},
			{
				path: 'action',
				component: () => import('../../views/System/ActionLog/index'),
				meta: { title: '操作日志', icon: 'el-icon-notebook-1' },
			},
		],
	},
]

const createRouter = () => {
	return new Router({
		scrollBehavior: () => ({ y: 0 }),
		routes: routers,
	})
}
const router = createRouter()
export function getRouters() {
	return routers
}
export default router
