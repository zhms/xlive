<template lang="pug">
.lottery
  .block-box
    .blocks.flex
      template(v-for="item in renderRewards")
        .start-btn.flex.flex-center(
          v-if="item == 'start'",
          @click="getRewards"
        )
          span START
        .block.flex.flex-center(
          v-else,
          :class="{ 'active-block': activeReward == item.name }"
        )
          img(:src="item.icon")
          span {{ item.name }}
</template>
<script setup>
import { computed, ref } from 'vue'
import qs from 'qs'
import { useIntervalFn } from '@vueuse/core'
import useMyFetch from '@/script/fetch'
import { checkFetchError, sleep, getWebSocket } from '@/script/base.js'
import { Popup } from 'vant'

const period = ref('')

const users = ref([])
const isPlay = ref(false)
const activeReward = ref('')
const rewards = ref([])
const walkRewards = computed(() => [rewards.value[0], rewards.value[1], rewards.value[2], rewards.value[4], rewards.value[7], rewards.value[6], rewards.value[5], rewards.value[3]])

const renderRewards = computed(() => {
	const arr = [...rewards.value]
	if (!arr.length) return []
	arr.splice(4, 0, 'start')
	return arr
})

const winPopupShow = ref(false)
const losePopupShow = ref(false)

useMyFetch('/api/choujiang/list', {
	immediate: true,
	afterFetch: (res) => {
		rewards.value = res.data.data.lottery_data
		period.value = res.data.data.period
	},
}).post(() => ({ period: period.value }))

function getRewards() {
	if (isPlay.value) return
	isPlay.value = true

	useMyFetch('/api/choujiang/getRewards', {
		immediate: true,
		afterFetch: async (res) => {
			if (checkFetchError(res)) {
				isPlay.value = false
				return
			}
			const reward = res.data.data.info

			// let walkLength = walkRewards.value.indexOf(activeReward.value) + 5 * 8;

			let walkLength = walkRewards.value.findIndex((item) => item.name == reward) + 5 * 8
			let i = 0

			const { pause } = useIntervalFn(async () => {
				if (walkLength-- == -1) {
					pause()
					await sleep(1000)
					isPlay.value = false
				} else {
					activeReward.value = walkRewards.value[i++ % 8].name
				}
			}, 50)
		},
	}).post(() => ({ period: period.value }))
}
</script>
<style lang="scss" scoped>
.lottery {
	width: 360px;
	height: 594px;
	margin: 0 auto;
	background: url('https://static.lotterybox.com/game/live/hongbao/3.png');
	background-size: cover;
	padding: 0 31px;
	padding-top: 230px;
}
.blocks {
	flex-wrap: wrap;
	.block,
	.start-btn {
		margin-top: 12px;
	}
	.block {
		width: 94px;
		height: 94px;
		background: #feeeee;
		border-radius: 8px;
		flex-direction: column;

		color: #b03f00;
		font-size: 14px;
		font-weight: bold;

		img {
			width: 60px;
			height: 60px;
		}

		span {
			margin-top: 2px;
		}
	}
	.active-block {
		background: #ffe000;
	}
	.start-btn {
		width: 80px;
		height: 80px;
		background: linear-gradient(180deg, #ff523e 0%, #ea170f 100%);
		box-shadow: inset 0px 3px 2px 0px #ff8065;
		border-radius: 40px;
		border: 2px solid #ff9d76;
		span {
			font-size: 20px;
			font-weight: bold;
			color: #ffffff;
		}
	}
}
</style>
