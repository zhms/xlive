<template>
	<div class="chat">
		<div class="message-list" ref="messageListDom">
			<template v-for="item in messages" :key="item.time">
				<div class="message-item" v-if="!deleteMessageList.includes(item.idClient) && item.text">
					<div class="time">{{ dayjs(item.time).format('MM-DD HH:mm:ss') }} {{ item.fromNick }}:</div>
					<div class="message-content">
						<span>{{ item.text }}</span>
					</div>
				</div>
			</template>
		</div>
		<div class="input-box flex" v-if="canSendMessage">
			<div class="input">
				<input placeholder="send message" v-model="inputMessage" @keydown.enter="sendMessage" />
			</div>
			<Button type="success" @click="sendMessage" :disabled="!chatRoomOk"> Send </Button>
		</div>
	</div>
</template>
<script setup>
import { ref, nextTick, computed } from 'vue'
import useMyFetch from '@/script/fetch.js'
import SDK from '@yxim/nim-web-sdk'
import { Button, Icon, showConfirmDialog, Popup, Field, showToast } from 'vant'
import { sleep, checkFetchError } from '@/script/base'
import { useStorage, useIntervalFn } from '@vueuse/core'
import dayjs from 'dayjs'

const inputMessage = ref('')
const chatRoomOk = ref(false)
const messages = ref([
	{
		chatroomId: '5991320983',
		idClient: '4cd5235279c18b144c798969a95f5809',
		from: 'blackrocktest_visitor_1039',
		fromNick: 'Visitor-65cf7e3808387',
		fromAvatar: 'https://blackrock.bochats.com/public/upload/images/userAvatar.jpg',
		fromCustom: '',
		userUpdateTime: 1708097374339,
		fromClientType: 'Web',
		time: 1708097390601,
		type: 'text',
		text: '34',
		resend: false,
		status: 'success',
		flow: 'in',
	},
	{
		chatroomId: '5991320983',
		idClient: '10a4c69c144d0d0afb6c725c6317a35a',
		from: 'blackrocktest_visitor_1044',
		fromNick: 'Visitor-65d03f81db050',
		fromAvatar: 'https://blackrock.bochats.com/public/upload/images/userAvatar.jpg',
		fromCustom: '',
		userUpdateTime: 1708146563963,
		fromClientType: 'Web',
		time: 1708146575178,
		type: 'text',
		text: '白',
		resend: false,
		status: 'success',
		flow: 'in',
	},
	{
		chatroomId: '5991320983',
		idClient: 'be1b3195165215f3def6b530450c152e',
		from: 'blackrocktest_visitor_1048',
		fromNick: 'A3352',
		fromAvatar: 'https://blackrock.bochats.com/public/upload/images/userAvatar.jpg',
		fromCustom: '',
		userUpdateTime: 1708278217475,
		fromClientType: 'Web',
		time: 1708278236690,
		type: 'text',
		text: '1111',
		resend: false,
		status: 'success',
		flow: 'in',
	},
])
const messageListDom = ref(null)
const blackWords = ref([])
const deleteMessageList = ref([])
const user = JSON.parse(useStorage('user').value)
const canSendMessage = ref(true)

let chatRoom
function initChatRoom(data) {
	chatRoom = SDK.Chatroom.getInstance({
		appKey: data.appKey,
		account: data.user.account,
		token: data.user.token,
		chatroomId: data.roomid,
		chatroomAddresses: data.chatRoom_Addr,
		onconnect: (obj) => {
			chatRoomOk.value = true
			getHistory()
			console.log('connect ok!', obj)
		},
		onerror: (error, obj) => {
			console.log('error', error, obj)
		},
		onwillreconnect: () => {},
		ondisconnect: () => {},
		onmsgs: (msgs) => {
			msgs.forEach((msg) => {
				if (msg.attach?.type === 'deleteChatroomMsg') {
					deleteMessageList.value.push(msg.attach.msgId)
				} else {
					mergerMessage([msg])
				}
			})
		},
	})
}

function getHistory() {
	chatRoom.getHistoryMsgs({
		done: (error, obj) => {
			mergerMessage(obj.msgs.reverse())
		},
	})
}

function mergerMessage(msgs) {
	messages.value = [...messages.value, ...msgs]
	if (messages.value.length > 200) {
		messages.value = messages.value.slice(messages.value.length - 200)
	}
	scroll()
}

async function scroll() {
	await nextTick()
	messageListDom.value.scrollTop = 100000000
}

function __sendMessage(text) {
	chatRoom.sendText({
		text,
		done: (error, msgObj) => {
			mergerMessage([msgObj])
		},
	})
}

function sendMessage() {
	if (!inputMessage.value.trim()) return
	const text = inputMessage.value
	__sendMessage(text)
	inputMessage.value = ''
}
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
