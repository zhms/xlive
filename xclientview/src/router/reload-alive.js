import { h, reactive } from 'vue'

const reload = reactive({
	value: false,
	timer: null,
})
const setReload = (value) => {
	reload.value = value
	clearTimeout(reload.timer)
	reload.timer = setTimeout(() => {
		reload.value = false
	}, 50)
}
const box = (page) => {
	return {
		name: 'ReloadAlive',
		components: {
			page,
		},
		data() {
			return {
				hook: 0,
			}
		},
		activated() {
			if (this.state$.inited && reload.value) {
				this.hook++
			}
		},
		render() {
			return h(page, { key: 'reload-alive-' + this.hook })
		},
	}
}
const init = (router) => {
	const rewriteRouter = (key) => {
		const f = router[key]
		router[key] = function (to) {
			if (typeof to === 'string') {
				to = { path: to }
			}
			setReload(to.reload === false ? false : true)
			return f.call(this, to).catch((err) => err)
		}
	}
	rewriteRouter('push')
	rewriteRouter('replace')

	// 重写 router.go
	const gof = router.go
	router.go = function (n, _reload = false) {
		setReload(_reload)
		return gof.call(this, n)
	}
}
const mixin = {
	data() {
		const state$ = {
			inited: false,
		}
		setTimeout(() => (state$.inited = true))
		return { state$ }
	},
	methods: {
		$checkSafeActivated() {
			return this.state$.inited && !reload.value
		},
	},
}

export { box, init, mixin }
