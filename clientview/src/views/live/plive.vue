<template lang="pug">
.live
  .user-box
    .logo
      img(src="https://static.lotterybox.com/game/live/logo111.jpg")
    .user-head.flex.flex-center
      div Online({{ onlineData?.data.online_count || 0 }})
    User

  .player
    Icon.play-icon(
      name="play-circle-o",
      @click="play",
      size="100",
      :style="{ left: (playerWidth - 100) / 2 + 'px', top: (playerHeight - 60) / 2 + 'px' }",
      v-if="!isPlay"
    )
    .teacher-info.flex
      div {{ liveData?.data.name }} | Current Lecturer: {{ liveData?.data.account }}
      .logout(v-if="isVisitor", @click="$router.push('/')") Login
      .logout(v-else, @click="logout") Logout
    canvas#canvas(
      :style="{ width: playerWidth + 'px', height: playerHeight + 'px' }"
    )
    video.video-js(
      :id="playerId",
      x-webkit-airplay="allow",
      webkit-playsinline,
      playsinline,
      preload="auto",
      :width="playerWidth",
      :height="playerHeight",
      @click="play",
      :poster="poster"
    )
    .course
      img(
        src="https://static.lotterybox.com/game/live/2023-12-30 10.54.47.jpg",
        :style="{ width: playerWidth + 'px', height: bodyHeight * 0.3 + 'px' }"
      )

  .chat-box
    NoticeBar(:text="liveData?.data.title", left-icon="volume-o")
    Chat
</template>
<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import useMyFetch from '@/fetch.js'
import { rootScale, bodyWidth, bodyHeight, sleep, logout } from '@/base'
import { Button, Icon, NoticeBar } from 'vant'
import { useStorage, useIntervalFn } from '@vueuse/core'
import Chat from './chat.vue'
import User from './user.vue'
import qs from 'qs'
const urlQuery = qs.parse(location.search.slice(1))

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

let player
const liveUrl = ref('')
const playData = computed(() => ({
	type: 'video/x-flv',
	src: liveUrl.value,
	isLive: true,
}))

const { data: liveData } = useMyFetch('/api/v1/app/get_live_info', {
	afterFetch: async (res) => {
		liveUrl.value = res.data.data.live_url
	},
}).get()

function play() {
	if (liveUrl.value == '') {
		showToast('The live is not ready yet, please try again later')
		return
	}
	player.setDataSource(playData.value)
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

// 在线人数
const { data: onlineData } = useMyFetch('/api/v1/app/get_online_info').get()

onMounted(() => {
	initPlayer()
})
</script>

<style lang="scss" scoped>
.player {
	flex-shrink: 0;
	position: relative;
	height: 100vh;
	overflow: hidden;
	border-right: 6px solid #3a4f7f;
	border-left: 6px solid #3a4f7f;

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
		padding: 0 10px;

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
	width: 300px;
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
	width: 400px;
	overflow: auto;
	height: 100%;

	:deep(.chat) {
		.message-list {
			height: calc(100vh - 50px - 40px);
			overflow: auto;
		}
	}
}
.course {
	img {
		width: 100%;
		height: 400px;
	}
}

.van-notice-bar {
	background: #2e4068 !important;
	color: #fff !important;
}
</style>
