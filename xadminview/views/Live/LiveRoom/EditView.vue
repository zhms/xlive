<template>
	<div class="dialogBox">
		<el-dialog style="margin-top: -100px" :title="title" :visible.sync="visable" width="630px" center @open="handleOpen">
			<el-form :inline="true" label-width="120px">
				<el-form-item label="直播间名称:">
					<el-input v-model="itemdata.name" style="width: 400px"></el-input>
				</el-form-item>
				<el-form-item label="讲师名称:">
					<el-input v-model="itemdata.account" style="width: 400px"></el-input>
				</el-form-item>
				<el-form-item label="聊天标题:">
					<el-input v-model="itemdata.title" style="width: 400px"></el-input>
				</el-form-item>
				<el-form-item label="状态:" v-if="title == '编辑直播间'">
					<el-radio v-model="itemdata.state" :label="1">开启直播</el-radio>
					<el-radio v-model="itemdata.state" :label="2">关闭直播</el-radio>
				</el-form-item>
			</el-form>
			<span slot="footer" class="dialog-footer">
				<el-button type="primary" @click="handleCommit">确定</el-button>
			</span>
		</el-dialog>
	</div>
</template>
<script>
import dlgbase from '@/api/dlgbase'
export default {
	extends: dlgbase,
	data() {
		return {}
	},
	methods: {
		onOpen() {},
		commitData(next) {
			if (this.title == '编辑直播间') {
				let data = JSON.parse(JSON.stringify(this.itemdata))
				this.$post('/v1/live_room/update_live_room', data, { google: true }).then(() => {
					this.$message.success('修改成功')
					next(true)
				})
			}
			if (this.title == '添加直播间') {
				if (!this.itemdata.name) return this.$message.error('请填写直播间名称')
				if (!this.itemdata.account) return this.$message.error('请填写讲师名称')
				if (!this.itemdata.title) return this.$message.error('请填写聊天标题')
				this.itemdata.state = this.itemdata.state ?? 2
				let data = JSON.parse(JSON.stringify(this.itemdata))
				this.$post('/v1/live_room/create_live_room', data, { google: true }).then(() => {
					this.$message.success('添加成功')
					next(true)
				})
			}
		},
	},
}
</script>
