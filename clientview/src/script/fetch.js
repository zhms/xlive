import { createFetch, useStorage } from '@vueuse/core'
import { getToken, login, appEnv } from './base'
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
				...appEnv,
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
			}
			return ctx
		},
	},
})

export default useMyFetch
