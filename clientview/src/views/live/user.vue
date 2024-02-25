<template>
	<div class="users">
		<div class="user flex" v-for="user in UserList" :key="user.username">
			<div>{{ user.username }}</div>
		</div>
	</div>
</template>
<script setup>
import useMyFetch from '@/script/fetch.js'
import { UserList } from '@/script/base'
// const { data: users } = useMyFetch('/api/yunxin/charRoom/user_online', {
// 	immediate: true,
// }).post()
</script>
<style lang="scss" scoped>
.users {
	background: #2e4068;
	color: #ffff;
	.user {
		height: 40px;
		border-bottom: 1px solid #ebedf0;
		padding: 10px;
		cursor: pointer;

		&:hover {
			font-weight: bold;
		}
	}
}
</style>
