<template>
	<div class="live">
		<div class="user-box">
			<div class="logo">
				<img src="https://static.lotterybox.com/game/live/logo111.jpg" />
			</div>
			<div class="user-head flex flex-center">
				<div>Online({{ OnlineCount || 0 }})</div>
			</div>
			<User />
		</div>
		<div class="player">
			<Icon name="play-circle-o" class="play-icon" @click="play" size="100" :style="{ left: playerWidth / 2 + 'px', top: (playerHeight - 60) / 2 + 'px' }" v-if="!isPlay"></Icon>
			<div class="teacher-info flex">
				<div>{{ liveData?.data.name }} | Current Lecturer: {{ liveData?.data.account }}</div>
				<div class="logout" @click="logout">Logout</div>
			</div>
			<canvas id="canvas"></canvas>
			<video :id="playerId" webkit-playsinline="true" playsinline="true" preload="auto" @click="play" :poster="poster" class="video-js"></video>

			<div class="course">
				<img src="https://static.lotterybox.com/game/live/2023-12-30 10.54.47.jpg" />
			</div>
		</div>

		<div class="chat-box">
			<NoticeBar :text="liveData?.data.title" left-icon="volume-o"></NoticeBar>
			<Chat />
		</div>
	</div>
</template>
<script setup>
import { Button, Icon, NoticeBar, showToast } from 'vant'
import { ref, computed, onMounted, watch } from 'vue'
import useMyFetch from '@/script/fetch.js'
import { rootScale, bodyWidth, bodyHeight, sleep, logout, wsconn, getLiveData, OnlineCount } from '@/script/base'
import { useStorage, useIntervalFn } from '@vueuse/core'
import Chat from './chat.vue'
import User from './user.vue'

const posterList = ref(['https://static.lotterybox.com/game/live/2023-12-31 20.21.45.jpg', 'https://static.lotterybox.com/game/live/2023-12-31 20.21.53.jpg', 'https://static.lotterybox.com/game/live/2023-12-31 20.21.57.jpg', 'https://static.lotterybox.com/game/live/2023-12-31 20.22.02.jpg'])
const currentPosterIndex = ref(0)
const poster = computed(() => posterList.value[currentPosterIndex.value])
useIntervalFn(() => {
	currentPosterIndex.value++
	if (currentPosterIndex.value >= posterList.value.length) {
		currentPosterIndex.value = 0
	}
}, 3000)

const user = JSON.parse(useStorage('user').value)
const isVisitor = computed(() => user.is_visitor == 1)

const playerWidth = computed(() => {
	return bodyWidth.value - 400 - 300 - 20
})

const playerHeight = computed(() => {
	return bodyHeight.value * 0.7 - 40
})

const playerId = ref('e' + +new Date())
const isPlay = ref(false)
const liveData = ref(getLiveData())
let player
const liveUrl = ref(liveData.value.data.pull_url)
console.log('liveUrl', liveUrl.value)
const playData = computed(() => ({
	type: 'video/x-flv',
	src: liveUrl.value,
	isLive: true,
}))

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
		() => {}
	)

	player.on('play', () => {
		isPlay.value = true
	})
}

watch(bodyWidth, () => {
	player.resize(playerWidth.value, playerHeight.value)
})

let onlineData = ref({ data: { online_count: 0 } })
// 在线人数
// const { data: onlineData } = useMyFetch('/api/v1/app/get_online_info').get()

onMounted(() => {
	initPlayer()
})

setTimeout(() => wsconn(), 1000)
</script>

<style lang="scss" scoped>
.player {
	background-color: #3a4f7f;
	flex-shrink: 0;
	position: relative;
	height: 100vh;
	overflow: hidden;
	border-right: 6px solid #3a4f7f;
	border-left: 6px solid #3a4f7f;
	flex: 1;
	.play-icon {
		position: absolute;
		z-index: 10;
		color: #fff;
	}
	#canvas {
		display: none;
		position: absolute;
		top: 0;
		left: 0;
		z-index: 10;
		flex: 1;
	}
	.teacher-info {
		height: 40px;
		background: #2e4068;
		color: #fff;
		font-weight: bold;
		padding: 0 10px;
		.flex {
			display: flex;
			align-items: center;
			justify-content: space-between;
		}

		.logout {
			color: #ccc;
			cursor: pointer;
		}
	}
}
.play-icon {
	position: absolute;
	z-index: 10;
	color: #fff;
}

.user-box {
	width: 250px;
	background: #2e4068;
	.logo {
		background: #fff;
		padding: 4px 20px;
		img {
			height: 80px;
		}
	}
	.user-head {
		font-size: 18px;
		font-weight: bold;
		color: #fff;
		height: 40px;
		text-align: center;
	}
	:deep(.users) {
		height: calc(100vh - 132px);
		overflow: auto;
	}
}
.live {
	display: flex;
}
.chat-box {
	flex-shrink: 0;
	position: relative;
	width: 350px;
	overflow: auto;
	height: 100%;
	margin-left: auto;

	:deep(.chat) {
		.message-list {
			height: calc(100vh - 50px - 40px);
			overflow: auto;
		}
	}
}
.course {
	height: 30%;
	display: flex;
	justify-content: center;
	align-items: center;
	img {
		max-height: 100%;
	}
}

.van-notice-bar {
	background: #2e4068 !important;
	color: #fff !important;
}
.video-js {
	background: #000;
	flex: 1;
	width: 100%;
	height: 70%;
}
</style>
