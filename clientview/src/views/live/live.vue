<template>
	<div class="box" :style="{ width: bodyWidth + 'px', height: bodyHeight + 'px' }">
		<div class="teacher-info flex flex-center">
			<div>{{ liveData?.data.name }} | Current Lecturer: {{ liveData?.data.account }}</div>
		</div>
		<Icon name="play-circle-o" class="play-icon" @click="play" size="100" v-if="!isPlay" :style="{ left: (playerWidth - 100) / 2 + 'px' }"></Icon>
		<video :id="playerId" x-webkit-airplay="allow" webkit-playsinline playsinline @click="play" :poster="poster" class="video-js video"></video>
		<div class="chat-box">
			<NoticeBar :text="liveData?.data.title" left-icon="volume-o"></NoticeBar>
			<Tabs>
				<Tab title="Chat">
					<Chat></Chat>
				</Tab>
				<Tab :title="'Online User(' + (OnlineCount || 0) + ')'">
					<User></User>
				</Tab>
			</Tabs>
		</div>
	</div>
</template>
<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import useMyFetch from '@/script/fetch.js'
import { rootScale, bodyWidth, bodyHeight, sleep, wsconn, getLiveData, OnlineCount } from '@/script/base'
import { Button, Icon, NoticeBar, Tab, Tabs, showToast } from 'vant'
import Chat from './chat.vue'
import User from './user.vue'
import { useStorage, useIntervalFn } from '@vueuse/core'

const posterList = ref(['https://static.lotterybox.com/game/live/2023-12-31 20.21.45.jpg', 'https://static.lotterybox.com/game/live/2023-12-31 20.21.53.jpg', 'https://static.lotterybox.com/game/live/2023-12-31 20.21.57.jpg', 'https://static.lotterybox.com/game/live/2023-12-31 20.22.02.jpg'])
const currentPosterIndex = ref(0)
const poster = computed(() => posterList.value[currentPosterIndex.value])
useIntervalFn(() => {
	currentPosterIndex.value++
	if (currentPosterIndex.value >= posterList.value.length) {
		currentPosterIndex.value = 0
	}
}, 3000)

const mediaWidth = ref(640)
const isMobile = computed(() => mediaWidth.value > bodyWidth.value)
const playerWidth = computed(() => {
	return bodyWidth.value
})
const playerHeight = computed(() => {
	return bodyWidth.value
})
const playerId = ref('e' + +new Date())
const isPlay = ref(false)
const liveData = ref(getLiveData())
let player
const liveUrl = ref(liveData.value.data.pull_url)
console.log('liveUrl', liveUrl.value)
const playData = computed(() => ({
	//type: 'video/x-flv',
	src: liveUrl.value,
	isLive: true,
}))

const user = JSON.parse(useStorage('user').value)
const isVisitor = computed(() => user.is_visitor == 1)

function play() {
	if (liveUrl.value == '') {
		showToast('The live is not ready yet, please try again later')
		return
	}
	console.log('playData', playData.value)
	player.setDataSource(playData.value)
	isPlay.value = true
}

function initPlayer(data) {
	player = window.neplayer(
		playerId.value,
		{
			controls: true,
			autoplay: true,
			loop: false,
			errMsg7: 'Live is over',
			streamTimeoutTime: 30 * 1000,
			techOrder: ['html5', 'flvjs'],
			bigPlayButton: false,
			decoderPath: 'https://yx-web-nosdn.netease.im/sdk-release/webplayer/',
			canvasId: 'canvas',
			enableStashBuffer: true,
		},
		() => {
			console.log('inited!')
		}
	)

	player.on('play', () => {
		console.log('play')
		isPlay.value = true
	})
}

watch(bodyWidth, () => {
	player.resize(playerWidth.value, playerHeight.value)
})

onMounted(() => {
	initPlayer()
})

setTimeout(() => wsconn(), 1000)
</script>

<style lang="scss" scoped>
.box {
	position: absolute;
	background-color: #2e4068;
	widows: 100%;
	//height: 100%;
}
.video {
	width: 100%;
	height: 250px;
}

.teacher-info {
	height: 40px;
	background: #2e4068;
	color: #fff;
	font-weight: bold;
}

.play-icon {
	position: absolute;
	z-index: 10;
	color: #fff;
	top: 100px;
}

.van-notice-bar {
	background: #2e4068 !important;
	color: #fff !important;
}

.chat-box {
	position: fixed;
	width: 100%;
	z-index: 100;
	box-sizing: border-box;
	bottom: 0;
	:deep(.chat) {
		.message-list {
			height: calc(100vh - 100vw - 100px);
			overflow: auto;
		}
	}
	:deep(.users) {
		height: calc(100vh - 100vw - 80px);
		overflow: auto;
	}
}
</style>
