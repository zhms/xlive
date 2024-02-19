<template>
	<div class="users">
		<div class="user flex" v-for="user in users?.data" :key="user.username">
			<div>{{ user.username }}</div>
		</div>
	</div>
</template>
<script setup>
import useMyFetch from '@/script/fetch.js'

let users = {
	data: [
		{
			username: 'Visitor-65d35a88d59d3',
		},
		{
			username: 'A1232',
		},
	],
}
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
