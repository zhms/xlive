<template>
	<div class="dialogBox">
		<el-dialog :title="title" :visible.sync="visable" width="600px" center>
			<el-checkbox v-for="col in columns" :key="col" v-model="data[col]" style="margin-bottom: 10px">{{ col }}</el-checkbox>
			<span slot="footer" class="dialog-footer">
				<el-button type="primary" @click="handleCommit">确定</el-button>
			</span>
		</el-dialog>
	</div>
</template>
<script>
export default {
	data() {
		return {
			title: '编辑列',
		}
	},
	props: {
		show: {
			type: Boolean,
			default: false,
		},
		columns: {
			type: Array,
			default: () => [],
		},
		column: {
			type: Object,
			default: () => {},
		},
	},
	created() {},
	computed: {
		visable: {
			get: function () {
				return this.show
			},
			set: function (val) {
				this.$emit('update:show', val)
			},
		},
		data: {
			get: function () {
				return this.column
			},
			set: function (v) {
				this.column = v
			},
		},
	},
	methods: {
		handleCommit() {
			this.visable = false
			this.$emit('setColumn', this.data)
		},
	},
}
</script>
