<template lang="pug">
.index
  .head.flex
    h2 Welcome To Live

  .desc
    h4 Welcome to live
    .tip There are some interesting things in this studio. Come and play. You can also invite your good friends to come with you!

  .form-box.flex
    .form
      .form-item
        input(
          type="text",
          placeholder="Enter Account",
          v-model="name",
          v-scrollInto
        )
      .form-item
        input(
          type="password",
          placeholder="Enter Password",
          v-model="pwd",
          v-scrollInto
        )

      Button(type="primary", block, :disabled="btnDisabled", @click="login") Login
      .visitor(@click="visitorLogin") Visitor Login
</template>
<script setup>
import { ref, computed } from 'vue'
import { Button } from 'vant'
import { useStorage } from '@vueuse/core'
import useMyFetch from '@/fetch.js'
import { checkFetchError, vScrollInto } from '@/base'
import { useRouter } from 'vue-router'

import qs from 'qs'
import { bodyWidth } from '@/base.js'
const urlQuery = qs.parse(location.search.slice(1))

const name = ref('')
const pwd = ref('')
const router = useRouter()

const mediaWidth = ref(640)
const isMobile = computed(() => mediaWidth.value > bodyWidth.value)

console.log(bodyWidth.value, isMobile.value)

const btnDisabled = computed(() => !name.value.trim() || !pwd.value || loginIsFetching.value)

const { execute: loginExecute, isFetching: loginIsFetching } = useMyFetch('/api/user/login', {
	afterFetch: (res) => {
		loginCallback(res)
	},
}).post(() => ({
	name: name.value,
	pwd: pwd.value,
	pid: urlQuery.pid,
}))

function loginCallback(res) {
	if (checkFetchError(res)) return
	useStorage('token').value = res.data.data.token
	useStorage('user').value = JSON.stringify(res.data.data)
	router.push(isMobile.value ? '/live' : 'plive')
}

function login() {
	loginExecute()
}

function visitorLogin() {
	useMyFetch('/api/user/visitorLogin', {
		immediate: true,
		afterFetch: (res) => {
			loginCallback(res)
		},
	}).post(() => ({
		pid: urlQuery.pid,
	}))
}
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
	background: url('@/assets/images/bg1.png');
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
