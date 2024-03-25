import Cookies from 'js-cookie'
const state = {
	sidebar: {
		opened: Cookies.get('sidebarStatus') ? !!+Cookies.get('sidebarStatus') : true,
		withoutAnimation: false,
	},
	sellers: [],
	channels: [],
	games: [
		{ id: 1, name: '哈希大小' },
		{ id: 2, name: '哈希单双' },
		{ id: 3, name: '幸运哈希' },
		{ id: 4, name: '幸运庄闲' },
		{ id: 5, name: '哈希牛牛' },
	],
	rooms: [
		{ id: 1, name: '初级场' },
		{ id: 2, name: '中级场' },
		{ id: 3, name: '高级场' },
	],
	states: [
		{
			id: 1,
			name: '启用',
		},
		{
			id: 2,
			name: '禁用',
		},
	],
}
const mutations = {
	TOGGLE_SIDEBAR: (state) => {
		state.sidebar.opened = !state.sidebar.opened
		state.sidebar.withoutAnimation = false
		if (state.sidebar.opened) {
			Cookies.set('sidebarStatus', 1)
		} else {
			Cookies.set('sidebarStatus', 0)
		}
	},
	SET_SELLERS: (state, sellers) => {
		state.sellers = sellers
	},
	SET_CHANNELS: (state, channels) => {
		state.channels = channels
	},
}
const actions = {
	toggleSideBar({ commit }) {
		commit('TOGGLE_SIDEBAR')
	},
	setSellers({ commit }, sellers) {
		commit('SET_SELLERS', sellers)
	},
	setChannels({ commit }, channels) {
		commit('SET_CHANNELS', channels)
	},
}
export default {
	namespaced: true,
	state,
	mutations,
	actions,
}
