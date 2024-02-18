<template lang="pug">
.live
  .player
    Icon.play-icon(
      name="play-circle-o",
      @click="play",
      size="100",
      :style="{ left: (playerWidth - 100) / 2 + 'px', top: (playerHeight - 60) / 2 + 'px' }",
      v-if="!isPlay"
    )
    .teacher-info.flex.flex-center
      div {{ liveData?.data.name }} | Current Lecturer: {{ liveData?.data.anchor.name }}
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

  .chat-box
    NoticeBar(:text="liveData?.data.meta_title", left-icon="volume-o")
    Tabs
      Tab(title="Chat")
        Chat
      Tab(:title="'Online User(' + (onlineData?.data.on_line_num || '') + ')'")
        User
</template>
<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import useMyFetch from '@/fetch.js'
import { rootScale, bodyWidth, bodyHeight, sleep } from '@/base'
import { Button, Icon, NoticeBar, Tab, Tabs } from 'vant'
import Chat from './chat.vue'
import User from './user.vue'
import { useStorage, useIntervalFn } from '@vueuse/core'
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

let player
const liveUrl = ref('')
const playData = computed(() => ({
	//type: 'application/x-mpegURL',
	type: 'video/flv',
	src: liveUrl.value,
	isLive: true,
}))

const { data: liveData } = useMyFetch('/api/yunxin/liveChannel', {
	immediate: true,
	afterFetch: async (res) => {
		liveUrl.value = 'https://pull.dbxapp.xyz/abc/abc.flv?auth_key=1708276065-0-0-ed7a1ac409877b492ec6052c07d6aa7a' //res.data.data.url;
		await sleep(500)
	},
}).post(() => ({ roomid: urlQuery.roomid }))

function play() {
	console.log('-----------: ', playData.value)
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

// 在线人数
const { data: onlineData } = useMyFetch('/api/yunxin/charRoom/info', {
	immediate: true,
}).post(() => ({ roomid: urlQuery.roomid }))

// 登录检测
useMyFetch('/api/user/islogin', {
	immediate: true,
	afterFetch: (res) => {
		// console.log(res);
	},
})

onMounted(() => {
	initPlayer()
})
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
