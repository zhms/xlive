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
			filters: {},
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
			this.AddItem(index, (title, item) => {
				item = item || {}
				this[`dialog${index}`].itemdata = item
				this[`dialog${index}`].title = title
				this[`dialog${index}`].show = true
			})
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
			let oInput = document.createElement('textarea')
			oInput.value = `${data}`
			document.body.appendChild(oInput)
			oInput.select()
			document.execCommand('copy')
			oInput.remove()
			Message({ type: 'success', message: '复制成功', center: true })
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
	},
}
