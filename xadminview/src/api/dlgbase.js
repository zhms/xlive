import { mapGetters } from 'vuex'
import enums from '@/api/enums'
export default {
	extends: enums,
	data() {
		return {
			seller_id: JSON.parse(sessionStorage.getItem('userinfo') ?? '{}').seller_id,
			zong: JSON.parse(sessionStorage.getItem('userinfo') || '{}').seller_id == -1,
			table_data: [],
		}
	},
	props: {
		show: {
			type: Boolean,
			default: false,
		},
		title: {
			type: String,
			default: null,
		},
		itemdata: {
			type: Object,
			default: () => {
				return {
					Password: '',
				}
			},
		},
		filters: {
			type: Object,
			default: () => {
				return {}
			},
		},
	},
	computed: {
		...mapGetters(['sellers', 'channels', 'games', 'rooms', 'states']),
		visable: {
			get: function () {
				return this.show
			},
			set: function (val) {
				this.$emit('update:show', val)
			},
		},
	},
	methods: {
		handleCommit() {
			this.commitData((closedialog) => {
				this.visable = !closedialog
				this.$emit('getTableData')
			})
		},
		handleOpen() {
			if (this.onOpen) this.onOpen()
		},
		getMap(arr) {
			let map = {}
			for (let i = 0; i < arr.length; i++) {
				map[arr[i].id] = arr[i].name
			}
			return map
		},
		getStateName(item) {
			return this.getMap(this.states)[item.State]
		},
		getchannel_name(channel_id) {
			let map = {}
			for (let i = 0; i < this.channels.length; i++) {
				map[this.channels[i].channel_id] = this.channels[i].channel_name
			}
			return map[channel_id]
		},
	},
}
