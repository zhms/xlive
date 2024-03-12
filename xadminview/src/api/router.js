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
		path: '/live',
		component: Layout,
		meta: { title: '直播间', icon: 'el-icon-video-camera-solid' },
		children: [
			{
				path: 'live_list',
				component: () => import('../../views/Live/LiveRoom/index'),
				meta: { title: '直播间列表', icon: 'el-icon-video-camera' },
			},
			{
				path: 'chat_list',
				component: () => import('../../views/Live/ChatList/index'),
				meta: { title: '互动列表', icon: 'el-icon-chat-dot-round' },
			},
			{
				path: 'banip_list',
				component: () => import('../../views/Live/BanIp/index'),
				meta: { title: 'Ip封禁', icon: 'el-icon-close' },
			},
		],
	},
	{
		path: '/user',
		component: Layout,
		meta: { title: '会员管理', icon: 'el-icon-user-solid' },
		children: [
			{
				path: 'user_list',
				component: () => import('../../views/User/index'),
				meta: { title: '会员管理', icon: 'el-icon-user-solid' },
			},
		],
	},
	{
		path: '/hongbao',
		component: Layout,
		meta: { title: '红包管理', icon: 'el-icon-present' },
		children: [
			{
				path: 'hongbao_list',
				component: () => import('../../views/Hongbao/index'),
				meta: { title: '红包管理', icon: 'el-icon-present' },
			},
		],
	},
	{
		path: '/data',
		component: Layout,
		meta: { title: '数据分析', icon: 'el-icon-cpu' },
		children: [
			{
				path: 'online_list',
				component: () => import('../../views/User/index'),
				meta: { title: '在线管理', icon: 'el-icon-document' },
			},
			{
				path: 'online_chart',
				component: () => import('../../views/User/index'),
				meta: { title: '在线图表', icon: 'el-icon-location' },
			},
			{
				path: 'peak_chart',
				component: () => import('../../views/User/index'),
				meta: { title: '峰值图表', icon: 'el-icon-location-outline' },
			},
		],
	},
	{
		path: '/robot',
		component: Layout,
		meta: { title: '机器人管理', icon: 'el-icon-service' },
		children: [
			{
				path: 'robot_list',
				component: () => import('../../views/User/index'),
				meta: { title: '机器人列表', icon: 'el-icon-service' },
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
				meta: { title: '登录日志', icon: 'el-icon-notebook-2' },
			},
			{
				path: 'action',
				component: () => import('../../views/System/ActionLog/index'),
				meta: { title: '操作日志', icon: 'el-icon-reading' },
			},
			{
				path: 'tools',
				component: () => import('../../views/System/Tools/index'),
				meta: { title: '系统工具', icon: 'el-icon-magic-stick' },
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
