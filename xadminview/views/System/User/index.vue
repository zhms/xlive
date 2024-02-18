<template>
	<div class="container">
		<el-form :inline="true" :model="filters">
			<!-- <el-form-item label="">
				<el-select v-model="filters.channel_id" placeholder="渠道" style="width: 150px" clearable @change="channelChange">
					<el-option v-for="item in channels" :key="item.channel_id" :label="item.channel_name" :value="item.channel_id"> </el-option>
				</el-select>
			</el-form-item> -->
			<el-form-item label="">
				<el-input v-model="filters.account" placeholder="账号" style="width: 150px" :clearable="true"></el-input>
			</el-form-item>
			<el-form-item>
				<el-button type="primary" icon="el-icon-refresh" v-on:click="handleQuery">查询</el-button>
				<el-button type="primary" icon="el-icon-plus" class="mr10" @click="handleAdd(0)">添加</el-button>
			</el-form-item>
		</el-form>
		<el-table :data="table_data" border class="table" max-height="670px" :cell-style="{ padding: '0px' }" :highlight-current-row="true">
			<el-table-column align="center" prop="id" label="id" width="80"></el-table-column>
			<el-table-column align="center" prop="account" label="账号" width="100">
				<template slot-scope="scope">
					<span style="cursor: pointer" class="blue" @click="handleEdit(scope.row, 0)">{{ scope.row.account }}</span>
				</template>
			</el-table-column>
			<el-table-column align="center" prop="SellerName" label="运营商" width="130" v-if="zong">
				<template slot-scope="scope">
					<span>{{ getSellerName(scope.row) }}</span>
				</template>
			</el-table-column>
			<!-- <el-table-column align="center" label="渠道" width="100">
				<template slot-scope="scope">
					<span>{{ getSellerName(scope.row) }}</span>
				</template>
			</el-table-column> -->
			<el-table-column align="center" prop="role_name" label="角色" width="150"></el-table-column>
			<el-table-column align="center" label="状态" width="100">
				<template slot-scope="scope">
					<span :class="scope.row.state == 1 ? 'blue' : 'red'">{{ getStateName(scope.row) }}</span>
				</template>
			</el-table-column>
			<el-table-column align="center" prop="login_time" label="登录时间" width="160"></el-table-column>
			<el-table-column align="center" prop="login_ip" label="登录ip" width="120"></el-table-column>
			<el-table-column align="center" prop="login_count" label="登录次数" width="100"></el-table-column>
			<el-table-column align="center" prop="memo" label="备注" width="200"></el-table-column>
			<el-table-column label="操作" align="left" width="300">
				<template slot-scope="scope">
					<el-button type="text" size="small" icon="el-icon-edit" @click="handleEdit(scope.row, 0)">编辑</el-button>
					<el-button type="text" size="small" icon="el-icon-edit" @click="handleEdit(scope.row, 1)">登录验证码</el-button>
					<el-button type="text" size="small" icon="el-icon-edit" @click="handleEdit(scope.row, 2)">操作验证码</el-button>
					<el-button type="text" size="small" icon="el-icon-delete" class="red" @click="handleDelete(scope.row)">删除</el-button>
				</template>
			</el-table-column>
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
			this.$get('/v1/admin_user/get_admin_user', data).then((result) => {
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
	},
}
</script>
