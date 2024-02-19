<template lang="pug">
.users
  .user.flex(v-for="user in users?.data")
    div {{ user.username }}
</template>
<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import useMyFetch from '@/script/fetch.js'

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
