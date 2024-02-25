<template>
	<div class="chat">
		<div class="message-list" ref="messageListDom">
			<template v-for="item in MsgList" :key="item.k">
				<div class="message-item">
					<div class="time">{{ item.time }} {{ item.from }}:</div>
					<div class="message-content">
						<span>{{ item.msg }}</span>
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
	</div>
</template>
<script setup>
import { ref, nextTick, computed } from 'vue'
import useMyFetch from '@/script/fetch.js'
import SDK from '@yxim/nim-web-sdk'
import { Button, Icon, showConfirmDialog, Popup, Field, showToast } from 'vant'
import { sleep, checkFetchError, MsgList, setScoll, sendChatMsg } from '@/script/base'
import { useStorage, useIntervalFn } from '@vueuse/core'

const inputMessage = ref('')
const messages = ref([])
const messageListDom = ref(null)
const user = JSON.parse(useStorage('user').value)

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
</style>
