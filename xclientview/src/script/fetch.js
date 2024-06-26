import { createFetch, useStorage } from '@vueuse/core'
import { getToken, login, RoomId } from './base'
import { showToast } from 'vant'

const useMyFetch = createFetch({
	baseUrl: '',
	options: {
		immediate: true,
		beforeFetch({ options }) {
			const token = getToken()
			options.headers = {
				...options.headers,
				'x-token': token,
				RoomId: RoomId,
			}
		},
		afterFetch(ctx) {
			if (typeof ctx.data === 'string') {
				ctx.data = JSON.parse(ctx.data)
			}
			return ctx
		},
		onFetchError(ctx) {
			let data = JSON.parse(ctx.data)
			if (data.code === 100201) {
				showToast('account not exist')
			} else if (data.code === 100202) {
				showToast('password not correct')
			} else if (data.code === 10) {
				login()
			} else if (data.code === 4) {
				showToast('living not available now')
			} else if (data.code === 100203) {
				showToast('account has been baned')
			}
			return ctx
		},
	},
})

export default useMyFetch
