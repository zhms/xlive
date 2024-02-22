import { ref, computed, onBeforeUnmount, watch } from 'vue'
import { useResizeObserver, useWebSocket, useStorage, useIntervalFn } from '@vueuse/core'
import mitt from 'mitt'
import { showToast } from 'vant'
import useMyFetch from './fetch.js'
import qs from 'qs'
import dsbridge from 'dsbridge'
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
let AppId = urlQuery['app']
if (!AppId) {
	AppId = '1'
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
const wsUrl = `${wsProtocol}://${location.host}/api/v1/app/ws/` + `${getToken()}_${AppId}`
let ws
function wsconn() {
	if (ws) return
	ws = useWebSocket(wsUrl, {
		onMessage: (ws, e) => {
			//console.log('ws', e.data)
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

// 封装websocket
function getWebSocket(params, callback = () => {}) {
	const result = ref([])
	params['Authorization'] = getToken()
	const { send, close, status, open } = useWebSocket(wsUrl, {
		onMessage: (ws, e) => {
			if (e.data === 'success\n') return
			const newData = JSON.parse(e.data)
			newData['_id'] = getLocalId()
			callback(newData)
			result.value = [newData, ...result.value].slice(0, 20)
		},
		heartbeat: {
			pongTimeout: 5000,
		},
		autoClose: false,
		onDisconnected: () => {
			console.log('disconnected', status.value)
		},
		onError: (res) => {
			console.log('error', res)
		},
	})
	send(JSON.stringify(params))

	// 组件卸载时候关闭
	onBeforeUnmount(() => {
		close()
	})
	// 心跳检测异常关闭则重连;
	useIntervalFn(() => {
		if (status.value === 'CLOSED') {
			open()
			send(JSON.stringify(params))
			console.log('reconnected')
		}
	}, 1000)
	return result
}

// 获取用户信息
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

function hasToken() {
	const token = dsbridge.call('Native.getToken') || useStorage('token').value
	return token !== '' && token !== 'undefined' && token !== 'empty' && token !== undefined
}

function getToken() {
	return dsbridge.call('Native.getToken') || useStorage('token').value
}

function logout() {
	ws.close()
	ws = null
	useStorage('token').value = ''
	useStorage('user').value = ''
	router.push('/index')
}

// 去登录
function login() {
	if (isApp) {
		dsbridge.call('Native.toLogin')
	} else {
		router.push('/')
	}
}

// miner
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
	AppId,
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
	getWebSocket,
	getUser,
	getToken,
	hasToken,
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
}
