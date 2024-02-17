import { Message } from 'element-ui'
import request from '@/api/request'
import { mapGetters } from 'vuex'
import store from '@/api/store'
import moment from 'moment'
import { isArray } from '@/api/utils'
import enums from '@/api/enums'
export default {
	extends: enums,
	data() {
		return {
			filters: {
				seller_id: Number(localStorage.getItem('seller_id') ?? 0) || null,
				channel_id: Number(localStorage.getItem('channel_id') ?? 0) || null,
			},
			seller_id: JSON.parse(sessionStorage.getItem('userinfo') ?? '{}').seller_id,
			zong: JSON.parse(sessionStorage.getItem('userinfo') ?? '{}').seller_id == -1,
			page: 1,
			page_size: 15,
			page_sizes: [15, 50, 100, 200, 500, 1000, 1500],
			total: 0,
			table_data: [],
			itemdata: null,
			Export: 0,
			columnpicker: false,
			column: {},
			dialog0: {
				show: false,
				title: '',
				itemdata: {},
			},
			dialog1: {
				show: false,
				title: '',
				itemdata: {},
			},
			dialog2: {
				show: false,
				title: '',
				itemdata: {},
			},
			dialog3: {
				show: false,
				title: '',
				itemdata: {},
			},
			dialog4: {
				show: false,
				title: '',
				itemdata: {},
			},
			dialog5: {
				show: false,
				title: '',
				itemdata: {},
			},
			dialog6: {
				show: false,
				title: '',
				itemdata: {},
			},
			dialog7: {
				show: false,
				title: '',
				itemdata: {},
			},
			dialog8: {
				show: false,
				title: '',
				itemdata: {},
			},
			dialog9: {
				show: false,
				title: '',
				itemdata: {},
			},
			symbols: [
				{
					id: 'trx',
					name: 'TRX',
				},
				{
					id: 'usdt',
					name: 'USDT',
				},
			],
		}
	},
	computed: {
		...mapGetters(['sellers', 'channels', 'games', 'rooms', 'states']),
	},
	created() {
		this.getSeller()
		if (this.columns) {
			let key = `column_${this.$vnode.key}`
			let data = localStorage.getItem(key)
			if (data) {
				this.column = JSON.parse(data)
			} else {
				for (let i = 0; i < this.columns.length; i++) {
					this.column[this.columns[i]] = true
				}
				localStorage.setItem(key, JSON.stringify(this.column))
			}
		}
	},
	methods: {
		getSeller() {
			request.get('/v1/seller/get_seller', { page_size: 1000 }, { noloading: true }).then((data) => {
				store.dispatch('app/setSellers', data)
			})
			//this.getChannel(this.filters.seller_id)
		},
		getChannel(seller_id) {
			if (!seller_id) {
				this.filters.channel_id = null
				this.channel = []
			} else {
				let data = {
					seller_id: seller_id,
					page: 1,
					page_size: 1000,
				}
				request.get('/v1/channel/get_channel', data, { noloading: true }).then((channels) => {
					store.dispatch('app/setChannels', channels.data)
				})
			}
		},
		sellerChange(seller_id) {
			localStorage.setItem('seller_id', seller_id)
			this.filters.channel_id = null
			localStorage.removeItem('channel_id')
			this.getChannel(seller_id || -1)
		},
		channelChange(channel_id) {
			localStorage.setItem('channel_id', channel_id)
		},
		pageChange(page) {
			this.page = page
			this.getTableData()
		},
		page_sizeChange(page_size) {
			this.page = 1
			this.page_size = page_size
			this.getTableData()
		},
		handleQuery(event) {
			if (event === 1) {
				this.Export = 1
				this.getTableData((result) => {
					this.Export = 0
					window.open('/sapi' + result.filename)
				})
			} else {
				this.Export = 0
				this.page = 1
				this.getTableData()
			}
		},
		handleAdd(index) {
			this[`dialog${index}`].itemdata = {}
			this.AddItem(
				index,
				(title) => {
					this[`dialog${index}`].title = title
					this[`dialog${index}`].show = true
				},
				this[`dialog${index}`].itemdata
			)
		},
		handleEdit(item, index) {
			this[`dialog${index}`].itemdata = JSON.parse(JSON.stringify(item))
			this.ModifyItem(
				index,
				(title) => {
					this[`dialog${index}`].title = title
					this[`dialog${index}`].show = true
				},
				this[`dialog${index}`].itemdata
			)
		},
		handleDelete(item, index) {
			this.DeleteItem(item, index)
		},
		getQueryData() {
			let retdata = {}
			retdata.page = this.page
			retdata.page_size = this.page_size
			let data = JSON.parse(JSON.stringify(this.filters))
			data.channel_id = data.channel_id || 0
			for (let i in data) {
				if (data[i] != null) {
					if (i == 'DateRange' && isArray(data[i])) {
						retdata['start_time'] = moment(data[i][0]).format('YYYY-MM-DD HH:mm:ss')
						retdata['end_time'] = moment(moment(data[i][1]).valueOf() + 86400000).format('YYYY-MM-DD HH:mm:ss')
					} else if (i == 'DateTimeRange' && isArray(data[i])) {
						retdata['start_time'] = moment(data[i][0]).format('YYYY-MM-DD HH:mm:ss')
						retdata['end_time'] = moment(moment(data[i][1]).valueOf() + 1000).format('YYYY-MM-DD HH:mm:ss')
					} else {
						retdata[i] = data[i]
						if (i.toLowerCase() == 'symbol') retdata[i] = data[i].toLowerCase()
					}
				}
			}
			if (this.Export == 1) retdata.Export = this.Export
			return retdata
		},
		copy(data) {
			let oInput = document.createElement('input')
			oInput.value = `${data}`
			document.body.appendChild(oInput)
			oInput.select()
			document.execCommand('copy')
			oInput.remove()
			Message({ type: 'success', message: '复制成功' })
		},
		getSummaries(param) {
			const { columns } = param
			const sums = []
			columns.forEach((column, index) => {
				if (index == 0) {
					if (this.getTotal) {
						let data = this.getTotal(0)
						if (data) {
							sums[index] = data
							return
						}
					}
					sums[index] = '总计'
					return
				}
				if (this.getTotal) {
					sums[index] = this.getTotal(column.property)
				}
			})
			return sums
		},
		getSellerName(item) {
			let map = {}
			for (let i = 0; i < this.sellers.length; i++) {
				map[this.sellers[i].seller_id] = this.sellers[i].seller_name
			}
			return map[item.seller_id]
		},
		getchannel_name(item) {
			let map = {}
			for (let i = 0; i < this.channels.length; i++) {
				map[this.channels[i].channel_id] = this.channels[i].channel_name
			}
			return map[item.channel_id]
		},
		getStateName(item) {
			return this.getMap(this.states)[item.state]
		},
		getMap(arr) {
			let map = {}
			for (let i = 0; i < arr.length; i++) {
				map[arr[i].id] = arr[i].name
			}
			return map
		},
		numDeal(t, n) {
			if (t == 0) return 0
			let f = t / Math.abs(t)
			n = n || 6
			if (n < 0) n = 0
			t = t * Math.pow(10, n)
			t = Math.abs(t)
			t = Math.floor(t)
			return (t / Math.pow(10, n)) * f
		},
		selectColumn() {
			this.columnpicker = true
		},
		setColumn(data) {
			for (let i in data) {
				this.column[i] = data[i]
			}
			let key = `column_${this.$vnode.key}`
			localStorage.setItem(key, JSON.stringify(this.column))
		},
		dealData(data) {
			data.forEach((element) => {
				if (element.CreateTime) element.CreateTime = moment(element.CreateTime).format('YYYY-MM-DD HH:mm:ss')
				if (element.LoginTime) element.LoginTime = moment(element.LoginTime).format('YYYY-MM-DD HH:mm:ss')
				if (element.RegisterTime) element.RegisterTime = moment(element.RegisterTime).format('YYYY-MM-DD HH:mm:ss')
			})
			return data
		},
	},
}
