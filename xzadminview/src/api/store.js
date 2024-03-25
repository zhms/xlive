import Vue from 'vue'
import Vuex from 'vuex'
Vue.use(Vuex)
const modulesFiles = require.context('../store', true, /\.js$/)
const modules = modulesFiles.keys().reduce((modules, modulePath) => {
	const moduleName = modulePath.replace(/^\.\/(.*)\.\w+$/, '$1')
	const value = modulesFiles(modulePath)
	modules[moduleName] = value.default
	return modules
}, {})
const getters = {
	sidebar: (state) => state.app.sidebar,
	routes: (state) => state.route.router,
	userinfo: (state) => state.user.userinfo,
	sellers: (state) => state.app.sellers,
	channels: (state) => state.app.channels,
	games: (state) => state.app.games,
	rooms: (state) => state.app.rooms,
	states: (state) => state.app.states,
}
export const store = new Vuex.Store({
	modules,
	getters,
})
export default store
//
