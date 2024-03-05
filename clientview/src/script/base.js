import { ref, computed, onBeforeUnmount, watch } from 'vue'
import { useResizeObserver, useWebSocket, useStorage, useIntervalFn } from '@vueuse/core'
import mitt from 'mitt'
import { showToast } from 'vant'
import useMyFetch from './fetch.js'
import qs from 'qs'
import router from '@/router/index.js'
import ClipboardJS from 'clipboard'

const eventBus = mitt()
const bodyWidth = ref(document.body.clientWidth)
const bodyHeight = ref(document.body.clientHeight)
useResizeObserver(document.body, () => {
	bodyWidth.value = document.body.clientWidth
	bodyHeight.value = document.body.clientHeight
})
const rootScale = computed(() => bodyWidth.value / 360)
const urlQuery = qs.parse(location.search.slice(1))
let RoomId = urlQuery['r']
if (!RoomId) {
	RoomId = '1'
}
let SaleId = urlQuery['s']
if (!SaleId) {
	SaleId = '0'
}

let scoll_callback = null
let OnlineCount = ref(0)
let UserList = ref([])
let MsgList = ref([])

function setScoll(cb) {
	scoll_callback = cb
}

function sleep(time = 1000) {
	return new Promise((resolve) => {
		setTimeout(resolve, time)
	})
}

// [0 - a)随机一个整数
function random(a = 10) {
	return parseInt(Math.random() * a)
}

// 产生一个1/a的概率
function rand(a = 10) {
	return random(a) === 0
}

function checkFetchError(res) {
	if (!res.data.code == 0) {
		showToast(res.data.message)
		return true
	}
	return false
}

const vScrollInto = {
	mounted: (el) => {
		el.addEventListener('focus', () => {
			el.scrollIntoViewIfNeeded()
		})
	},
}

const vClipboard = {
	mounted: (el) => {
		const clipboard = new ClipboardJS(el)
		clipboard.on('success', () => {
			showToast('Copy Successful!')
		})
	},
}

function useBodyBgColor(color) {
	document.body.style.backgroundColor = color
	onBeforeUnmount(() => {
		document.body.style.backgroundColor = ''
	})
}

function deepClone(item) {
	return JSON.parse(JSON.stringify(item))
}

function formatMoney(money) {
	return (money || '0').toString().replace(/\B(?=(\d{3})+\b)/g, ',')
}

function padZero(number) {
	return ('0' + number).slice(-2)
}

function upFirstLetter(word) {
	return word.replace(/^(\w)/, (s) => s.toUpperCase())
}

// local id
let localId = +new Date()
function getLocalId() {
	return localId++
}

function Defer() {
	let resolve, reject
	const promise = new Promise((_resolve, _reject) => {
		resolve = _resolve
		reject = _reject
	})
	return {
		resolve,
		reject,
		promise,
	}
}

function preloadImage(url) {
	const defer = Defer()
	const img = document.createElement('img')
	img.onload = () => {
		defer.resolve()
	}
	img.src = url
	return defer.promise
}

// ws地址
const wsProtocol = location.protocol === 'http:' ? 'ws' : 'wss'
let wsUrl = `${wsProtocol}://${location.host}/api/v1/app/ws/` + `${getToken()}_${RoomId}`
let ws
function wsclose() {
	if (ws) {
		ws.close()
		ws = null
	}
}
function wsconn(token) {
	if (ws) return
	if (token) {
		wsUrl = `${wsProtocol}://${location.host}/api/v1/app/ws/` + `${token}_${RoomId}`
	}
	ws = useWebSocket(wsUrl, {
		onMessage: (ws, e) => {
			try {
				const data = JSON.parse(e.data)
				if (data.msg_id == 'user_count') {
					OnlineCount.value = data.msg_data
				} else if (data.msg_id == 'user_list') {
					UserList.value = []
					data.msg_data.forEach((item) => {
						UserList.value.push({
							username: item,
						})
					})
				} else if (data.msg_id == 'user_come') {
					UserList.value.unshift({
						username: data.msg_data,
					})
				} else if (data.msg_id == 'user_leave') {
					const index = UserList.value.findIndex((item) => item.username == data.msg_data)
					UserList.value.splice(index, 1)
				} else if (data.msg_id == 'chat') {
					let msgdata = JSON.parse(data.msg_data)
					if (MsgList.value.length > 200) {
						MsgList.value.shift()
					}
					msgdata.k = random(100000000)
					MsgList.value.push(msgdata)
					if (scoll_callback) {
						scoll_callback()
					}
				} else if (data.msg_id == 'chat_limit') {
					showToast('Chat too fast, please try again later')
				} else if (data.msg_id == 'chat_ban') {
					showToast('You have been banned from chatting')
				}
			} catch (e) {}
		},
		heartbeat: {
			message: 'ping',
			interval: 5000,
			pongTimeout: 30000,
		},
		onDisconnected: () => {},
		onError: (res) => {
			console.log('error', res)
		},
		onConnected: () => {
			console.log('ws connected')
		},
	})
}

function sendChatMsg(msg) {
	if (ws) {
		ws.send(
			JSON.stringify({
				msg_id: 'chat',
				msg_data: {
					msg: msg,
				},
			})
		)
	}
}

function getUser() {
	const { data } = useMyFetch('/game/user/info', { immediate: true })
	return data
}

function getLiveData() {
	let d = useStorage('user').value
	if (d) {
		return {
			data: JSON.parse(JSON.parse(d).live_data),
		}
	}
}

function getToken() {
	const token = useStorage('token').value
	if (token == '' && token == 'undefined' && token !== 'empty' && token !== undefined) return
	return token
}

function logout() {
	ws.close()
	ws = null
	useStorage('token').value = ''
	useStorage('user').value = ''
	router.push('/index')
}

function login() {
	router.push('/')
}

function miner_getStageTypeByGameType(gameType) {
	return Math.pow(2, ~~gameType / 10)
}

function miner_handleOrder(item) {
	if (item.mine_position) {
		item.pass_list.push({
			is_mine: 1,
			position: item.mine_position,
			is_cold: item.bonus != 0,
		})
	}
	const stageType = miner_getStageTypeByGameType(item.game_type)
	item.stageType = stageType
	item.blocks = miner_createBlocks(stageType, item.pass_list)
	return item
}

function miner_createBlocks(type, passList = []) {
	const passListObject = {}
	passList.forEach((item) => {
		passListObject[item.position] = item
	})

	const result = []
	for (var i = 0; i < type; i++) {
		const row = []
		for (var j = 0; j < type; j++) {
			const id = (type * i + j).toString()
			let status = 0
			if (passListObject[id]) {
				status = passListObject[id].is_mine == 0 ? 1 : -1
				if (passListObject[id].is_cold) {
					status = -2
				}
			}
			row.push({
				id,
				clicking: false,
				bonus: 0,
				// 0 默认状态，1 无雷， -1 雷, -2 冷雷
				status,
			})
		}
		result.push(row)
	}
	return result
}

function fly_getOddsColor(odds) {
	if (odds >= 100.01) return '#645cff'
	if (odds >= 5.01) return '#0082FE'
	if (odds >= 2.01) return '#00BFF2'
	if (odds >= 1.21) return '#00C67C'
	return '#F85169'
}

function parity_getColor(name) {
	return { green: '#00A280', violet: '#5841FF', red: '#FF6821' }[name]
}

export {
	bodyWidth,
	bodyHeight,
	rootScale,
	eventBus,
	RoomId,
	SaleId,
	urlQuery,
	vScrollInto,
	vClipboard,
	sleep,
	random,
	rand,
	checkFetchError,
	miner_createBlocks,
	miner_handleOrder,
	miner_getStageTypeByGameType,
	getLocalId,
	getUser,
	getToken,
	useBodyBgColor,
	deepClone,
	formatMoney,
	padZero,
	upFirstLetter,
	preloadImage,
	login,
	fly_getOddsColor,
	parity_getColor,
	logout,
	wsconn,
	ws,
	getLiveData,
	OnlineCount,
	UserList,
	MsgList,
	setScoll,
	sendChatMsg,
	wsclose,
}
