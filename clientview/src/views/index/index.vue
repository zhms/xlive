<template>
	<div class="index">
		<div class="head flex">
			<h2>Welcome To Live</h2>
		</div>

		<div class="desc">
			<h4>Welcome to live</h4>
			<div class="tip">There are some interesting things in this studio. Come and play. You can also invite your good friends to come with you!</div>
		</div>

		<div class="form-box flex">
			<div class="form">
				<div class="form-item">
					<input type="text" placeholder="Enter Account" v-model="name" v-scrollInto />
				</div>
				<div class="form-item">
					<input type="password" placeholder="Enter Password" v-model="pwd" v-scrollInto />
				</div>

				<Button type="primary" block :disabled="btnDisabled" @click="login"> Login </Button>
				<div class="visitor" @click="visitorLogin">Visitor Login</div>
			</div>
		</div>
	</div>
</template>
<script setup>
import { ref, computed } from 'vue'
import { Button } from 'vant'
import { useStorage } from '@vueuse/core'
import useMyFetch from '@/script/fetch.js'
import { useRouter } from 'vue-router'
import { bodyWidth } from '@/script/base.js'
import { SaleId, wsconn } from '../../script/base'
const name = ref('')
const pwd = ref('')
const router = useRouter()

const mediaWidth = ref(640)
const isMobile = computed(() => mediaWidth.value > bodyWidth.value)

const btnDisabled = computed(() => !name.value.trim() || !pwd.value || loginIsFetching.value)

const { execute: loginExecute, isFetching: loginIsFetching } = useMyFetch('/api/v1/user/user_login', {
	immediate: false,
	afterFetch: (res) => {
		loginCallback(res)
	},
}).post(() => ({
	account: name.value,
	password: pwd.value,
	is_visitor: 2,
	sale_id: SaleId,
}))

function loginCallback(res) {
	useStorage('token').value = res.data.data.token

	setTimeout(() => wsconn(res.data.data.token), 1000)

	useStorage('user').value = JSON.stringify(res.data.data)
	router.push(isMobile.value ? '/live' : 'plive')
}

function login() {
	loginExecute()
}

function generateRandomString(length) {
	const characters = 'abcdefghijklmnopqrstuvwxyz0123456789'
	let result = ''

	for (let i = 0; i < length; i++) {
		const randomIndex = Math.floor(Math.random() * characters.length)
		result += characters.charAt(randomIndex)
	}

	return result
}

function visitorLogin() {
	let account = 'Visitor_u' + generateRandomString(15)
	try {
		let user = useStorage('user').value
		if (user) {
			user = JSON.parse(user)
			if (user.is_visitor == 1) {
				account = user.account
			}
		}
	} catch (e) {}
	useMyFetch('/api/v1/user/user_login', {
		immediate: true,
		afterFetch: (res) => {
			loginCallback(res)
		},
	}).post(() => ({
		account: account,
		is_visitor: 1,
		sale_id: SaleId,
	}))
}

// if (useStorage('token').value && useStorage('token').value != 'null' && useStorage('token').value != 'undefined') {
// 	router.push(isMobile.value ? '/live' : 'plive')
// }
</script>

<style lang="scss" scoped>
.index {
	overflow: hidden;
	max-width: 800px;
	margin: 0 auto;
}
.head {
	height: 80px;
	background: #000;
	color: #ffffff;
	font-size: 20px;
	padding: 0 20px;
	font-weight: bold;
}

.desc {
	padding: 20px;
	background: #fff;
	line-height: 24px;

	.tip {
		margin-top: 4px;
	}
}

.form-box {
	background: url('@/images/bg1.png');
	background-size: cover;
	padding: 100px 0;
	height: 490px;

	.form {
		width: 80%;
		background: #fff;
		padding: 60px 16px;
		margin: 0 auto;

		.form-item {
			margin-bottom: 16px;
			border: 1px solid #e0e3ec;
			padding: 4px;
			border-radius: 2px;
			input {
				border: none;
				font-size: 16px;
				line-height: 2;
				width: 100%;
			}
		}
	}
}
.visitor {
	text-align: right;
	padding-top: 20px;
	padding-right: 10px;
	cursor: pointer;
}
</style>
