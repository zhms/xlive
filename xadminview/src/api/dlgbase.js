import { mapGetters } from 'vuex'
import enums from '@/api/enums'
export default {
	extends: enums,
	data() {
		return {
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
	},
}
