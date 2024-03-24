<template>
	<div class="chat">
		<div class="message-list" ref="messageListDom">
			<template v-for="item in MsgList" :key="item.k">
				<div class="message-item">
					<div class="time">{{ moment(item.time).format('YYYY-MM-DD HH:mm:ss') }} {{ isHongbao(item.msg) ? '' : item.from }}</div>
					<div class="message-content">
						<span v-if="!isHongbao(item.msg)">{{ item.msg }}</span>
						<img v-if="isHongbao(item.msg)" src="@/images/hongbao1.png" @click="openHongbao(item.msg)" style="cursor: pointer" />
					</div>
				</div>
			</template>
		</div>
		<div class="input-box flex">
			<div class="input">
				<input placeholder="send message" v-model="inputMessage" @keydown.enter="sendMessage" />
			</div>
			<Button type="success" @click="sendMessage" :disabled="inputMessage == ''"> Send </Button>
		</div>
		<Popup v-model:show="hongbaoShow" closeable :close-on-click-overlay="false">
			<div class="hongbao-box">
				<div class="title">Congratulations</div>
				<div class="title">{{ user.account }}</div>
				<div class="money">₹ {{ hongbaoAmount }}</div>
			</div>
		</Popup>
	</div>
</template>
<script setup>
import { ref, nextTick, computed } from 'vue'
import useMyFetch from '@/script/fetch.js'
import { Button, Icon, showConfirmDialog, Popup, Field, showToast } from 'vant'
import { sleep, checkFetchError, MsgList, setScoll, sendChatMsg } from '@/script/base'
import { useStorage, useIntervalFn } from '@vueuse/core'
import moment from 'moment'

const inputMessage = ref('')
const messages = ref([])
const messageListDom = ref(null)
const user = JSON.parse(useStorage('user').value)

const hongbaoShow = ref(false)
const hongbaoAmount = ref(0)

async function scroll() {
	await nextTick()
	messageListDom.value.scrollTop = 100000000
}

function sendMessage() {
	if (!inputMessage.value.trim()) return
	const text = inputMessage.value
	sendChatMsg(text)
	inputMessage.value = ''
}

function isHongbao(msg) {
	return msg.indexOf('__hongbao__') >= 0
}

function openHongbao(id) {
	id = id.replace('__hongbao__', '')
	useMyFetch('/api/v1/open_hongbao', {
		immediate: true,
		afterFetch: (res) => {
			if (res.data.data.amount == -1) {
				showToast('bonus has been opened')
				return
			}
			if (res.data.data.amount == -2) {
				showToast('bonus not found')
				return
			}
			if (res.data.data.amount == -3) {
				showToast('bonus has been expired')
				return
			}
			if (res.data.data.amount == -4) {
				showToast('bonus finished')
				return
			}
			hongbaoAmount.value = Math.floor(res.data.data.amount * 100) / 100
			hongbaoShow.value = !hongbaoShow.value
		},
	}).post(() => ({
		id: Number(id),
	}))
}

setScoll(scroll)
</script>
<style lang="scss" scoped>
.chat {
	background: #2e4068;
	position: relative;
	.message-list {
		padding: 10px 0;
		.message-item {
			padding: 8px 2px;
			font-weight: bold;
			word-break: break-all;

			.time {
				color: #fff;
				font-size: 12px;
				margin-bottom: 4px;
			}

			.message-content {
				color: #fff;
			}
			.revoke {
				margin-left: 8px;
				font-weight: bold;
				color: #ed6a0c;
				cursor: pointer;
			}
		}
	}
	.input-box {
		.input {
			width: 100%;
			height: 42px;
			input {
				width: 100%;
				border: none;
				height: 100%;
			}
		}
	}

	.controls {
		position: absolute;
		top: 10px;
		right: 20px;
		img {
			display: block;
			margin-top: 10px;
		}
	}
}

.message-item {
	margin-left: 5px;
}

.hongbao-box {
	width: 290px;
	height: 354px;
	background: url('@/images/hongbao2.png');
	text-align: center;
	padding-top: 80px;
	background-color: blueviolet;
	.title {
		font-weight: bold;
		color: #a56e2d;
		line-height: 19px;
	}
	.money {
		font-size: 40px;
		font-weight: bold;
		color: #956122;
		margin-top: 20px;
	}
}
</style>
