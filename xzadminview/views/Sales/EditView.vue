<template>
	<div class="dialogBox">
		<el-dialog style="margin-top: -100px" :title="title" :visible.sync="visable" width="630px" center @open="handleOpen">
			<el-form :inline="true" label-width="120px">
				<el-form-item label="账号:">
					<el-input v-model="itemdata.account" :disabled="title == '编辑账号'" style="width: 400px"></el-input>
				</el-form-item>
				<el-form-item label="密码:">
					<el-input v-model="itemdata.password" show-password style="width: 400px"></el-input>
				</el-form-item>
				<!-- <el-form-item label="角色:">
					<el-select v-model="itemdata.role_name" placeholder="请选择" style="width: 400px">
						<el-option v-for="item in dlgroles" :key="item.role_name" :label="item.role_name" :value="item.role_name"></el-option>
					</el-select>
				</el-form-item> -->
				<!-- <el-form-item label="状态:">
					<el-checkbox border label="禁用" v-model="itemdata.state" :true-label="2" :false-label="1"></el-checkbox>
				</el-form-item> -->
				<el-form-item label="房间号:">
					<el-select v-model="itemdata.room_id" placeholder="请选择" style="width: 300px">
						<el-option v-for="item in room_id" :key="item.id" :label="item.value" :value="item.id"></el-option>
					</el-select>
				</el-form-item>
				<el-form-item label="备注:">
					<el-input type="textarea" v-model="itemdata.memo" :rows="4" style="width: 400px"></el-input>
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
		return {
			dlgroles: [],
			dlgchannels: [],
			room_id: [],
		}
	},
	methods: {
		commitData(next) {
			if (this.title == '编辑账号') {
				if (!this.itemdata.room_id) return this.$message.error('请选择房间号')
				let data = JSON.parse(JSON.stringify(this.itemdata))
				if (data.password && data.password.length > 0) {
					data.password = this.$md5(data.password)
				}
				this.$post('/v1/update_sales', data, { google: true }).then(() => {
					this.$message.success('修改成功')
					next(true)
				})
			}
			if (this.title == '添加账号') {
				if (!this.itemdata.account) return this.$message.error('请填写账号')
				if (!this.itemdata.password) return this.$message.error('请填写密码')
				if (!this.itemdata.room_id) return this.$message.error('请选择房间号')
				let data = JSON.parse(JSON.stringify(this.itemdata))
				data.password = this.$md5(data.password)
				this.$post('/v1/create_sales', data, { google: true }).then(() => {
					this.$message.success('添加成功')
					next(true)
				})
			}
		},
		onOpen() {
			this.dlgroles = []
			this.dlgchannels = []
			if (this.title == '编辑账号') {
				delete this.itemdata.password
			}
			if (this.title == '添加账号') {
				this.itemdata.state = 1
				this.itemdata.role_name = '业务员'
			}
			this.$post('/v1/get_live_room_id', {}).then((result) => {
				this.room_id = []
				for (let i = 0; i < result.ids.length; i++) {
					this.room_id.push({ id: result.ids[i], value: result.ids[i] })
				}
			})
		},
	},
}
</script>
