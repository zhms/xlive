<template>
	<div :class="classObj" class="app-wrapper">
		<sidebar class="sidebar-container" />
		<div :class="{ hasTagsView: true }" class="main-container">
			<div :class="{ 'fixed-header': true }">
				<navbar />
				<tags-view v-if="true" />
			</div>
			<app-main />
		</div>
	</div>
</template>

<script>
import { mapState } from 'vuex'
import { AppMain, Sidebar, TagsView, Navbar } from './index.js'
export default {
	name: 'Layout',
	components: {
		AppMain,
		Sidebar,
		TagsView,
		Navbar,
	},
	computed: {
		...mapState({
			sidebar: (state) => state.app.sidebar,
		}),
		classObj() {
			return {
				hideSidebar: !this.sidebar.opened,
				openSidebar: this.sidebar.opened,
				withoutAnimation: false,
			}
		},
	},
}
</script>

<style lang="scss" scoped>
@import '~@/styles/mixin.scss';
@import '~@/styles/variables.scss';
.app-wrapper {
	@include clearfix;
	position: relative;
	height: 100%;
	width: 100%;
	&.mobile.openSidebar {
		position: fixed;
		top: 0;
	}
}
.drawer-bg {
	background: #000;
	opacity: 0.3;
	width: 100%;
	top: 0;
	height: 100%;
	position: absolute;
	z-index: 999;
}
.fixed-header {
	position: fixed;
	top: 0;
	right: 0;
	z-index: 9;
	width: calc(100% - #{$sideBarWidth});
	transition: width 0.28s;
}
.hideSidebar .fixed-header {
	width: calc(100% - 54px);
}
.mobile .fixed-header {
	width: 100%;
}
</style>
