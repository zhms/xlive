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
		}
	},
	methods: {
		commitData(next) {
			if (this.title == '编辑账号') {
				let data = JSON.parse(JSON.stringify(this.itemdata))
				if (data.password && data.password.length > 0) {
					data.password = this.$md5(data.password)
				}
				this.$post('/v1/admin_update_user', data, { google: true }).then(() => {
					this.$message.success('修改成功')
					next(true)
				})
			}
			if (this.title == '添加账号') {
				if (!this.itemdata.account) return this.$message.error('请填写账号')
				if (!this.itemdata.password) return this.$message.error('请填写密码')
				if (!this.itemdata.role_name) return this.$message.error('请选择角色')
				let data = JSON.parse(JSON.stringify(this.itemdata))
				data.password = this.$md5(data.password)
				this.$post('/v1/admin_create_user', data, { google: true }).then(() => {
					this.$message.success('添加成功')
					next(true)
				})
			}
		},
		onOpen() {
			this.dlgroles = []
			this.dlgchannels = []
			// this.getRoles()
			if (this.title == '编辑账号') {
				delete this.itemdata.password
			}
			if (this.title == '添加账号') {
				this.itemdata.state = 1
				this.itemdata.role_name = '超级管理员'
			}
			console.log(this.itemdata)
		},
		sellerChange() {
			this.getRoles()
		},
		getRoles() {
			let data = {}
			this.$post('/v1/admin_get_role', data, { noloading: true }).then((roledata) => {
				this.dlgroles = roledata.data
			})
		},
	},
}
</script>
