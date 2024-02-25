<template>
	<div class="container">
		<el-form :inline="true" :model="filters">
			<el-form-item label="">
				<el-input v-model="filters.user_id" placeholder="会员Id" style="width: 150px" :clearable="true"></el-input>
			</el-form-item>
			<el-form-item label="">
				<el-input v-model="filters.account" placeholder="账号" style="width: 150px" :clearable="true"></el-input>
			</el-form-item>
			<el-form-item label="">
				<el-input v-model="filters.agent" placeholder="业务员" style="width: 150px" :clearable="true"></el-input>
			</el-form-item>
			<el-form-item label="">
				<el-input v-model="filters.login_ip" placeholder="登录Ip" style="width: 150px" :clearable="true"></el-input>
			</el-form-item>
			<el-form-item>
				<el-button type="primary" icon="el-icon-refresh" v-on:click="handleQuery">查询</el-button>
				<el-button type="primary" icon="el-icon-plus" class="mr10" @click="handleAdd(0)">添加</el-button>
				<el-button type="primary" icon="el-icon-upload2" class="mr10" @click="handleAdd(0)">批量导入</el-button>
				<el-button type="primary" icon="el-icon-bottom" class="mr10" @click="handleAdd(0)">导出</el-button>
				<el-button type="primary" icon="el-icon-download" class="mr10" @click="DownLoadExcelTemplate()">下载导入模板</el-button>
			</el-form-item>
		</el-form>
		<el-table :data="table_data" border class="table" max-height="670px" :cell-style="{ padding: '0px' }" :highlight-current-row="true">
			<el-table-column align="center" prop="id" label="id" width="100"></el-table-column>
			<el-table-column align="center" prop="user_id" label="会员Id" width="100"></el-table-column>
			<el-table-column align="center" prop="account" label="账号" width="200">
				<template slot-scope="scope">
					<span style="cursor: pointer" class="blue" @click="handleEdit(scope.row, 0)">{{ scope.row.account }}</span>
				</template>
			</el-table-column>
			<el-table-column align="center" prop="agent" label="业务员" width="120"></el-table-column>
			<el-table-column align="center" label="状态" width="100">
				<template slot-scope="scope">
					<span :class="scope.row.state == 1 ? 'blue' : 'red'">{{ getStateName(scope.row) }}</span>
				</template>
			</el-table-column>
			<el-table-column align="center" prop="login_count" label="登录次数" width="100"></el-table-column>
			<el-table-column align="center" prop="login_ip" label="登录Ip" width="150"></el-table-column>
			<el-table-column align="center" prop="login_location" label="登录地区" width="120"></el-table-column>
			<el-table-column align="center" prop="create_time" label="注册时间" width="160"></el-table-column>
			<el-table-column align="center" prop="memo" label="备注" width="200"></el-table-column>
			<!--<el-table-column label="操作" align="left" width="300">
				<template slot-scope="scope">
					<el-button type="text" size="small" icon="el-icon-edit" @click="handleEdit(scope.row, 0)">编辑</el-button>
					<el-button type="text" size="small" icon="el-icon-edit" @click="handleEdit(scope.row, 1)">登录验证码</el-button>
					<el-button type="text" size="small" icon="el-icon-edit" @click="handleEdit(scope.row, 2)">操作验证码</el-button>
					<el-button type="text" size="small" icon="el-icon-delete" class="red" @click="handleDelete(scope.row)">删除</el-button>
				</template>
			</el-table-column> -->
		</el-table>
		<div class="pagination">
			<el-pagination style="margin-right: 10px" background layout="sizes,total, prev, pager, next, jumper" @size-change="page_sizeChange" :current-page="page" @current-change="pageChange" :page-sizes="page_sizes" :total="total" :page-size="page_size"></el-pagination>
		</div>
		<edit-view :show.sync="dialog0.show" :title="dialog0.title" :itemdata="dialog0.itemdata" :filters="filters" @getTableData="getTableData" />
		<div>
			<el-dialog :title="dialog1.title" :visible.sync="dialog1.show" width="350px" center>
				<vueqr :text="dialog1.url" :size="300"> </vueqr>
			</el-dialog>
		</div>
	</div>
</template>
<script>
import '@/styles/main.css'
import base from '@/api/base'
import vueqr from 'vue-qr'
import EditView from './EditView.vue'
import XLSX from 'xlsx'
import moment from 'moment'
export default {
	extends: base,
	components: { EditView, vueqr },
	data() {
		return {}
	},
	created() {
		this.getTableData()
	},
	methods: {
		getTableData() {
			let data = this.getQueryData()
			this.$get('/v1/user/get_user', data).then((result) => {
				this.table_data = this.dealData(result.data)
				this.total = result.total
			})
		},
		ModifyItem(index, next, item) {
			if (index == 0) next('编辑账号')
			if (index == 1 || index == 2) {
				let data = {
					Account: item.account,
				}
				if (index == 1) {
					this.$post('/v1/admin_user/set_login_googlesecret', data, { google: true }).then((result) => {
						this.dialog1.show = true
						this.dialog1.url = result.url
						this.dialog1.title = '登录验证码'
					})
				} else if (index == 2) {
					this.$post('/v1/admin_user/set_opt_googlesecret', data, { google: true }).then((result) => {
						this.dialog1.show = true
						this.dialog1.url = result.url
						this.dialog1.title = '操作验证码'
					})
				}
			}
		},
		AddItem(index, next) {
			if (index == 0) next('添加账号')
		},
		DeleteItem(item) {
			this.$delete('/v1/admin_user/delete_admin_user', item, { google: true }).then(() => {
				this.$message.success('删除成功')
				this.getTableData()
			})
		},
		DownLoadExcelTemplate() {
			console.log('下载导入模板')
			const data = XLSX.utils.aoa_to_sheet([['账号', '密码']])
			const wb = XLSX.utils.book_new()
			XLSX.utils.book_append_sheet(wb, data, 'kalacloud-data')
			XLSX.writeFile(wb, `导入会员模板_${moment().format('YYYY-MM-DD HH-mm-ss')}.xlsx`)
		},
	},
}
</script>
