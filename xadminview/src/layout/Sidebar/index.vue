<template>
	<div>
		<logo v-if="true" :collapse="isCollapse" />
		<el-scrollbar wrap-class="scrollbar-wrapper">
			<el-menu :default-active="activeMenu" :collapse="isCollapse" :background-color="variables.menuBg" :text-color="variables.menuText" :unique-opened="false" :active-text-color="variables.menuActiveText" :collapse-transition="false" mode="vertical">
				<sidebar-item v-for="route in routes" :key="route.path" :item="route" :base-path="route.path" />
			</el-menu>
		</el-scrollbar>
	</div>
</template>

<script>
import { mapGetters } from 'vuex'
import Logo from './Logo'
import SidebarItem from './SidebarItem'
import variables from '@/styles/variables.scss'
import { getRouters } from '@/api/router'
export default {
	data() {
		return {
			routes: [],
		}
	},
	components: { SidebarItem, Logo },
	computed: {
		...mapGetters(['userinfo']),
		...mapGetters(['sidebar']),
		activeMenu() {
			const route = this.$route
			const { meta, path } = route
			if (meta.activeMenu) {
				return meta.activeMenu
			}
			return path
		},
		variables() {
			return variables
		},
		isCollapse() {
			return !this.sidebar.opened
		},
	},
	created() {
		let AuthData = JSON.parse(JSON.parse(sessionStorage.getItem('userinfo')).auth_data)
		let routes = getRouters()
		let final_routes = []
		for (let i = 0; i < routes.length; i++) {
			if (!routes[i].meta) {
				final_routes.push(routes[i])
			} else {
				if (routes[i].children.length == 1) {
					let a = routes[i].meta.title
					if (AuthData[a] && AuthData[a]['查']) {
						final_routes.push(routes[i])
					}
				} else {
					let m = {
						path: routes[i].path,
						component: routes[i].Layout,
						meta: routes[i].meta,
						children: [],
					}
					for (let j = 0; j < routes[i].children.length; j++) {
						let a = routes[i].meta.title
						let b = routes[i].children[j].meta.title
						if (AuthData[a] && AuthData[a][b] && AuthData[a][b]['查']) {
							m.children.push(routes[i].children[j])
						}
					}
					if (m.children.length > 0) {
						final_routes.push(m)
					}
				}
			}
		}
		this.routes = final_routes
	},
}
</script>
