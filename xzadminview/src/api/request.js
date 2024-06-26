import axios from 'axios'
import { MessageBox, Message, Loading } from 'element-ui'
import { saveAs } from 'file-saver'
const service = axios.create({
	baseURL: '/api',
	timeout: 120000,
})
let loading
service.interceptors.request.use(
	(config) => {
		config.headers['x-token'] = sessionStorage.getItem('x-token') || ''
		return config
	},
	(error) => {
		return Promise.reject(error)
	}
)
service.interceptors.response.use(
	(response) => {
		if (loading) loading.close()
		if (response.headers['content-type'] == 'application/octet-stream') {
			let filename = response.headers['content-disposition'].split('filename=')[1]
			filename = decodeURIComponent(filename)
			saveAs(new Blob([response.data]), filename)
			return { code: 0, data: [] }
		}
		const res = response.data
		if (res.code != 0) {
			Message({ type: 'error', message: res.data, center: true })
			return Promise.reject(new Error(res.data))
		} else {
			return res.data
		}
	},
	(error, b) => {
		if (loading) loading.close()
		if (error.response.status == 500) {
			Message({ type: 'error', message: 'http status 500', center: true })
		} else if (error.response.status == 404) {
			Message({ type: 'error', message: 'http status 404', center: true })
		} else {
			if (error.response.data.code == 100102) {
				error.response.data.data = error.response.data.data || ' 请退出从新登录'
				Message({ type: 'error', message: error.response.data.msg + ' ' + error.response.data.data, center: true })
			} else {
				error.response.data.data = error.response.data.data || ''
				Message({ type: 'error', message: error.response.data.msg + ' ' + error.response.data.data, center: true })
			}
		}
		return Promise.reject(error)
	}
)

async function getGoogleCode() {
	try {
		let ret = await MessageBox.confirm('确认操作', '', {
			confirmButtonText: '确定',
			cancelButtonText: '取消',
		})
		return '123'
	} catch (e) {}
	return null
}

export default {
	async download(url, data = {}, options = {}) {
		data.export = 1
		options.loading = !options.noloading
		if (options.loading) {
			if (loading) loading.close()
			loading = Loading.service({ lock: true, spinner: 'el-icon-loading', background: 'rgba(0, 0, 0, 0.7)' })
		}
		return new Promise((resolve, reject) => {
			let options = {}
			options.url = url
			options.data = data
			options.method = 'POST'
			options.responseType = 'arraybuffer'
			service(options)
				.then((data) => {
					if (options.loading) if (loading) loading.close()
					resolve(data)
				})
				.catch((e) => {
					if (options.loading) if (loading) loading.close()
				})
		})
	},
	async post(url, data = {}, options = {}) {
		let googlecode
		options.loading = !options.noloading
		if (options.google) {
			googlecode = await getGoogleCode()
			if (!googlecode) {
				return new Promise((resolve, reject) => {
					reject()
				})
			}
		}
		if (options.loading) {
			if (loading) loading.close()
			loading = Loading.service({ lock: true, spinner: 'el-icon-loading', background: 'rgba(0, 0, 0, 0.7)' })
		}
		return new Promise((resolve, reject) => {
			let options = {}
			options.url = url
			options.data = data
			options.method = 'POST'
			options.headers = { VerifyCode: googlecode }
			service(options)
				.then((data) => {
					if (options.loading) if (loading) loading.close()
					resolve(data)
				})
				.catch((e) => {
					if (options.loading) if (loading) loading.close()
				})
		})
	},
}
