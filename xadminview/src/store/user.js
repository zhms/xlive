import request from '../api/request'
import moment from 'moment'
const state = {
	userinfo: JSON.parse(sessionStorage.getItem('userinfo') || '{}'),
}
const mutations = {
	SET_USERINFO: (state, userinfo) => {
		sessionStorage.setItem('userinfo', JSON.stringify(userinfo))
		if (userinfo.seller_id != -1) localStorage.setItem('seller_id', userinfo.seller_id)
		state.userinfo = userinfo
	},
}
const actions = {
	login({ commit }, data) {
		return new Promise((resolve) => {
			request.post('/v1/admin_user/user_login', data, { google: true }).then((userinfo) => {
				userinfo.LoginTime = moment(userinfo.LoginTime).format('YYYY-MM-DD HH:mm:ss')
				sessionStorage.setItem('x-token', userinfo.token)
				commit('SET_USERINFO', userinfo)
				resolve()
			})
		})
	},
}
export default {
	namespaced: true,
	state,
	mutations,
	actions,
}
