<template lang="pug">
.chat
  .controls(v-if="isAdmin")
    img.hongbao(
      src="https://static.lotterybox.com/game/live/hongbao/2.png",
      @click="() => (hongbaoCreaterShow = true)"
    )
    img.choujiang(
      src="https://static.lotterybox.com/game/live/hongbao/4.png",
      @click="() => sendChoujiang()"
    )

  .message-list(ref="messageListDom")
    template(v-for="item in messages")
      .message-item(
        v-if="!deleteMessageList.includes(item.idClient) && item.text"
      )
        .time {{ dayjs(item.time).format("MM-DD HH:mm:ss") }}
        .message-content
          span {{ item.fromNick }}:
          img.hongbao(
            v-if="item.text.indexOf('__hongbao__') == 0",
            src="https://static.lotterybox.com/game/live/hongbao/2.png",
            @click="() => getHongbao(item.text)"
          )
          img.choujiang(
            v-else-if="item.text.indexOf('__choujiang__') == 0",
            src="https://static.lotterybox.com/game/live/hongbao/4.png",
            @click="() => (choujiangShow = true)"
          )
          span(v-else) {{ item.text }}
        Icon.revoke(
          name="revoke",
          @click="confirmRevoke(item)",
          v-if="isAdmin"
        )

  .input-box.flex(v-if="canSendMessage")
    .input
      input(
        placeholder="send message",
        v-model="inputMessage",
        @keydown.enter="sendMessage"
      )
    Button(type="success", @click="sendMessage", :disabled="!chatRoomOk") Send

Popup(v-model:show="hongbaoCreaterShow", :style="{ width: '300px' }")
  .hongbao-create
    .title Send Bonus
    Field(
      v-model="hongbaoForm.amount",
      placeholder="Enter Amount",
      label="Amount"
    )
    Field(v-model="hongbaoForm.num", placeholder="Enter Number", label="Count")
    Field(
      v-model="hongbaoForm.max",
      placeholder="Enter Max Amount",
      label="Max Amt"
    )
    .buttons
      Button(type="primary", @click="sendHongbao") Send
      Button(@click="() => (hongbaoCreaterShow = false)") Cancel

Popup(v-model:show="hongbaoShow", closeable, :close-on-click-overlay="false")
  .hongbao-box
    .title Congratulations
    .money ₹ {{ currentHongbao }}

Popup(
  v-model:show="choujiangShow",
  closeable,
  :close-on-click-overlay="false",
  :style="{ background: 'transparent' }"
)
  Lottery
</template>
<script setup>
import { ref, nextTick, computed } from 'vue'
import useMyFetch from '@/fetch.js'
import SDK from '@yxim/nim-web-sdk'
import { Button, Icon, showConfirmDialog, Popup, Field, showToast } from 'vant'
import { sleep, checkFetchError } from '@/base'
import { useStorage, useIntervalFn } from '@vueuse/core'
import dayjs from 'dayjs'
import Lottery from './lottery.vue'
import qs from 'qs'
const urlQuery = qs.parse(location.search.slice(1))

// hongbao
const hongbaoCreaterShow = ref(false)
const hongbaoForm = ref({
	amount: 0,
	num: 0,
	max: 0,
})
const currentHongbao = ref(0)
const hongbaoShow = ref(false)
function sendHongbao() {
	if (hongbaoForm.value.amount == 0) {
		showToast('Amount is error!')
		return
	}
	if (hongbaoForm.value.num == 0) {
		showToast('Number is error!')
		return
	}
	hongbaoCreaterShow.value = false
	useMyFetch('/api/redPacket/create', {
		immediate: true,
		afterFetch: (res) => {
			__sendMessage(`__hongbao__${res.data.data.period}`)
		},
	}).post(() => ({
		amount: hongbaoForm.value.amount,
		num: hongbaoForm.value.num,
		max: hongbaoForm.value.max,
	}))
}
function getHongbao(text) {
	const period = text.replace('__hongbao__', '')
	useMyFetch('/api/redPacket/rob', {
		immediate: true,
		afterFetch: (res) => {
			if (checkFetchError(res)) {
				return
			}
			currentHongbao.value = res.data.data.money
			hongbaoShow.value = true
		},
	}).post(() => ({
		period,
	}))
}

// choujiang
const choujiangShow = ref(false)
function sendChoujiang() {
	__sendMessage(`__choujiang__`)
}

const inputMessage = ref('')
const chatRoomOk = ref(false)
const messages = ref([])
const messageListDom = ref(null)
const blackWords = ref([])
const deleteMessageList = ref([])
const user = JSON.parse(useStorage('user').value)
const isAdmin = computed(() => user?.is_admin === 1)
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
		// 消息
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
	const text = getWhiteWords()
	__sendMessage(text)
	inputMessage.value = ''
}

function getWhiteWords() {
	let words = inputMessage.value
	blackWords.value.forEach((item) => {
		words = words.replaceAll(item, '******************************'.slice(0, item.length))
	})
	return words
}

function confirmRevoke(item) {
	showConfirmDialog({
		confirmButtonText: 'Yes',
		cancelButtonText: 'no',
		message: 'Confirm revoke?',
	})
		.then(() => {
			revoke(item)
		})
		.catch(() => {})
}

function revoke(item) {
	deleteMessageList.value.push(item.idClient)
	useMyFetch('/api/yunxin/charRoom/recall_msg', {
		immediate: true,
	}).post(() => ({
		roomid: item.chatroomId,
		msgTimetag: item.time,
		msgId: item.idClient,
		fromAcc: item.from,
		operatorAcc: item.from,
	}))
}

useMyFetch('/api/yunxin/charRoom/roomid', {
	immediate: true,
	afterFetch: (res) => {
		initChatRoom(res.data.data)
	},
}).post(() => ({ roomid: urlQuery.roomid }))

useMyFetch('/api/yunxin/banword', {
	immediate: true,
	afterFetch: (res) => {
		blackWords.value = res.data.data.map((item) => item.replacefrom)
	},
}).post()

// 登录检测
const { execute: isLoginExecute } = useMyFetch('/api/user/islogin', {
	immediate: true,
	afterFetch: (res) => {
		if (res.data.data.status != 0) {
			canSendMessage.value = false
		}
	},
})

useIntervalFn(() => {
	isLoginExecute()
}, 5000)
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
.hongbao,
.choujiang {
	width: 100px;
}

.hongbao-create {
	.title {
		padding: 20px 0;
		text-align: center;
		font-weight: bold;
	}

	.buttons {
		padding: 10px 0;
		text-align: center;
		.van-button + .van-button {
			margin-left: 10px;
		}
	}
}

.hongbao-box {
	width: 290px;
	height: 354px;
	background: url('https://static.lotterybox.com/game/live/hongbao/1.png');
	text-align: center;
	padding-top: 80px;
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
