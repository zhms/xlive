<template>
	<section class="app-main">
		<transition name="fade-transform" mode="out-in">
			<keep-alive>
				<router-view :key="key" />
			</keep-alive>
		</transition>
	</section>
</template>

<script>
export default {
	name: 'AppMain',
	computed: {
		key() {
			return this.$route.path
		},
	},
}
</script>

<style lang="scss" scoped>
.app-main {
	min-height: calc(100vh - 50px);
	width: 100%;
	position: relative;
	overflow: hidden;
}

.fixed-header + .app-main {
	padding-top: 75px;
}
</style>

<style lang="scss">
.el-popup-parent--hidden {
	.fixed-header {
		padding-right: 15px;
	}
}
</style>
