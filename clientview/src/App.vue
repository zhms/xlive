<style lang="scss" scoped></style>

<template>
	<UseScreenSafeArea top right bottom left class="inner-app">
		<router-view v-slot="{ Component }">
			<component :is="Component" />
		</router-view>
	</UseScreenSafeArea>
</template>

<script setup>
import { UseScreenSafeArea } from '@vueuse/components'
import { rootScale, bodyWidth } from './base'
import { useCssVar } from '@vueuse/core'
import { watchEffect } from 'vue'
import useMyFetch from '@/fetch.js'

watchEffect(() => {
	useCssVar('--rootScale').value = rootScale.value
	useCssVar('--bodyWidth').value = bodyWidth.value + 'px'
})

// topbar z-index 是1w，所以遮罩从2w开始
useCssVar('--van-overlay-z-index').value = 2000

// 状态上报
const { execute: onlineExecute } = useMyFetch('/api/v1/user/user_login', {
	immediate: true,
}).post({
	account: 'fuck',
	password: 'fuck',
	is_visitor: 1,
})
// setInterval(onlineExecute, 60 * 1000)
</script>
