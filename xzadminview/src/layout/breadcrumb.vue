<template>
	<el-breadcrumb class="app-breadcrumb" separator="/">
		<transition-group name="breadcrumb">
			<el-breadcrumb-item v-for="item in levelList" :key="item.path">
				<span>{{ item.meta.title }}</span>
			</el-breadcrumb-item>
		</transition-group>
	</el-breadcrumb>
</template>
<script>
import pathToRegexp from 'path-to-regexp'
export default {
	data() {
		return {
			levelList: null,
		}
	},
	watch: {
		$route(route) {
			if (route.path.startsWith('/redirect/')) {
				return
			}
			this.getBreadcrumb()
		},
	},
	created() {
		this.getBreadcrumb()
	},
	methods: {
		getBreadcrumb() {
			let matched = this.$route.matched.filter((item) => item.meta && item.meta.title)
			this.levelList = matched.filter((item) => item.meta && item.meta.title && item.meta.breadcrumb !== false)
		},
		isDashboard(route) {
			const name = route && route.name
			if (!name) {
				return false
			}
			return name.trim().toLocaleLowerCase() === 'Dashboard'.toLocaleLowerCase()
		},
		pathCompile(path) {
			const { params } = this.$route
			var toPath = pathToRegexp.compile(path)
			return toPath(params)
		},
		handleLink(item) {
			const { redirect, path } = item
			if (redirect) {
				this.$router.push(redirect)
				return
			}
			this.$router.push(this.pathCompile(path))
		},
	},
}
</script>

<style lang="scss" scoped>
.app-breadcrumb.el-breadcrumb {
	display: inline-block;
	font-size: 14px;
	line-height: 50px;
	margin-left: 8px;

	.no-redirect {
		color: #97a8be;
		cursor: text;
	}
}
</style>
