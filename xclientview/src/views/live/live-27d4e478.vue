<template>
	<div class="live">
		<div class="player">
			<Icon name="play-circle-o" class="play-icon" @click="play" size="100" :style="{ left: (playerWidth - 100) / 2 + 'px', top: (playerHeight - 60) / 2 + 'px' }" v-if="!isPlay"></Icon>
			<div class="teacher-info flex flex-center">
				<div>{{ liveData?.data.name }} | Current Lecturer: {{ liveData?.data.account }}</div>
			</div>
			<canvas id="canvas" :style="{ width: playerWidth + 'px', height: playerHeight + 'px' }"></canvas>
			<video :id="playerId" x-webkit-airplay="allow" webkit-playsinline playsinline preload="auto" :width="playerWidth" :height="playerHeight" @click="play" :poster="poster" class="video-js"></video>
		</div>
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

const playData = computed(() => ({
	type: 'video/x-flv',
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
.player {
	position: relative;
	#canvas {
		display: none;
		position: absolute;
		top: 0;
		left: 0;
		z-index: 10;
	}
	.teacher-info {
		height: 40px;
		background: #2e4068;
		color: #fff;
		font-weight: bold;
	}
}
.play-icon {
	position: absolute;
	z-index: 10;
	color: #fff;
}
.chat-box {
	position: fixed;
	bottom: 0;
	width: 100%;
	z-index: 100;
	padding: 0 4px;
	box-sizing: border-box;

	:deep(.chat) {
		.message-list {
			height: calc(100vh - 100vw - 42px - 90px - 44px);
			overflow: auto;
		}
	}
	:deep(.users) {
		height: calc(100vh - 100vw - 90px - 42px);
		overflow: auto;
	}
}
.van-notice-bar {
	background: #2e4068 !important;
	color: #fff !important;
}
</style>
