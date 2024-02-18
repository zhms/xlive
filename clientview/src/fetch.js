import { createFetch, useStorage } from '@vueuse/core'
import { getToken, login, appEnv } from './base'
import { showToast } from 'vant'

const useMyFetch = createFetch({
	baseUrl: '',
	options: {
		immediate: false,
		beforeFetch({ options }) {
			const token = getToken()
			options.headers = {
				...options.headers,
				Token: token,
				...appEnv,
			}
		},
		afterFetch(ctx) {
			if (typeof ctx.data === 'string') {
				ctx.data = JSON.parse(ctx.data)
			}

			if (ctx.data.code === 401) {
				login()
			}

			return ctx
		},
		onFetchError(ctx) {
			// showToast(JSON.parse(ctx.data).message)
			// console.error(ctx);
			return ctx
		},
	},
})

export default useMyFetch
