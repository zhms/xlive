<template>
	<div class="container">
		<el-form :inline="true" :model="filters">
			<el-form-item label="">
				<el-input v-model="filters.account" placeholder="账号" style="width: 150px" :clearable="true"></el-input>
			</el-form-item>
			<el-form-item label="">
				<el-input v-model="filters.agent" placeholder="业务员" style="width: 150px" :clearable="true"></el-input>
			</el-form-item>
			<el-form-item label="">
				<el-input v-model="filters.login_ip" placeholder="登录Ip" style="width: 150px" :clearable="true"></el-input>
			</el-form-item>
			<el-form-item label="">
				<el-select v-model="filters.is_online" placeholder="在线状态" style="width: 100px" clearable>
					<el-option v-for="item in online_options" :key="item.value" :label="item.label" :value="item.value"> </el-option>
				</el-select>
			</el-form-item>
			<el-form-item>
				<el-button type="primary" icon="el-icon-refresh" v-on:click="handleQuery">查询</el-button>
				<el-button type="primary" icon="el-icon-plus" class="mr10" @click="handleAdd(0)">添加</el-button>
				<el-button type="primary" icon="el-icon-upload2" class="mr10" @click="handleAdd(1)">批量导入</el-button>
				<el-button type="primary" icon="el-icon-bottom" class="mr10" @click="ExportUser()">导出</el-button>
				<el-button type="primary" icon="el-icon-download" class="mr10" @click="DownLoadExcelTemplate()">下载导入模板</el-button>
			</el-form-item>
		</el-form>
		<el-row>
			<el-col :span="24" style="margin-top: -10px; padding-bottom: 10px">
				<div class="grid-content">
					<span style="margin-right: 30px; color: #909399">总人数 : {{ total }}</span>
					<span style="margin-right: 30px; color: #909399">在线人数 : {{ online_count }}</span>
					<span style="margin-right: 30px; color: #909399">合并Ip数 : {{ ip_total }}</span>
				</div>
			</el-col>
		</el-row>
		<el-table :data="table_data" style="margin-top: -10px" border class="table" max-height="670px" :cell-style="{ padding: '0px' }" :highlight-current-row="true">
			<el-table-column align="center" prop="id" label="id" width="100"></el-table-column>
			<el-table-column align="center" prop="account" label="账号" width="200">
				<template slot-scope="scope">
					<span style="cursor: pointer" class="blue" @click="handleEdit(scope.row, 0)">{{ scope.row.account }}</span>
				</template>
			</el-table-column>
			<el-table-column align="center" prop="agent" label="业务员" width="120"></el-table-column>
			<el-table-column align="center" label="状态" width="100">
				<template slot-scope="scope">
					<span :class="scope.row.state != 1 ? 'red' : ''">{{ scope.row.state | 启用禁用 }}</span>
				</template>
			</el-table-column>
			<el-table-column align="center" label="在线状态" width="100">
				<template slot-scope="scope">
					<span>{{ scope.row.is_online == 1 ? '在线' : '离线' }}</span>
				</template>
			</el-table-column>
			<el-table-column align="center" prop="login_count" label="登录次数" width="100"></el-table-column>
			<el-table-column align="center" prop="login_ip" label="登录Ip" width="150"></el-table-column>
			<el-table-column align="center" prop="login_ip_location" label="登录地区" width="200"></el-table-column>
			<el-table-column align="center" prop="create_time" label="注册时间" width="160">
				<template slot-scope="scope">
					<span>{{ scope.row.create_time | 北京时间 }}</span>
				</template>
			</el-table-column>

			<el-table-column label="操作" align="left" width="300">
				<template slot-scope="scope">
					<el-button type="text" class="red" icon="el-icon-edit" @click="handleEdit(scope.row, 3)" v-if="scope.row.state == 1">封号</el-button>
					<el-button type="text" icon="el-icon-edit" @click="handleEdit(scope.row, 4)" v-if="scope.row.state != 1">解封</el-button>

					<el-button type="text" class="red" icon="el-icon-edit" @click="handleEdit(scope.row, 2)" v-if="scope.row.chat_state == 1">禁言</el-button>
					<el-button type="text" icon="el-icon-edit" @click="handleEdit(scope.row, 1)" v-if="scope.row.chat_state != 1">解除禁言</el-button>
				</template>
			</el-table-column>
		</el-table>
		<div class="pagination">
			<el-pagination style="margin-right: 10px" background layout="sizes,total, prev, pager, next, jumper" @size-change="page_sizeChange" :current-page="page" @current-change="pageChange" :page-sizes="page_sizes" :total="total" :page-size="page_size"></el-pagination>
		</div>
		<edit-view :show.sync="dialog0.show" :title="dialog0.title" :itemdata="dialog0.itemdata" :filters="filters" @getTableData="getTableData" />
		<edit-view :show.sync="dialog1.show" :title="dialog1.title" :itemdata="dialog1.itemdata" :filters="filters" @getTableData="getTableData" />
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
		return {
			online_options: [
				{ value: 1, label: '在线' },
				{ value: 2, label: '离线' },
			],
			online_count: 0,
			ip_total: 0,
		}
	},
	created() {
		this.getTableData()
	},
	methods: {
		getTableData() {
			let data = this.getQueryData()
			data.is_online = data.is_online == '' ? 0 : data.is_online
			this.$post('/v1/get_user', data).then((result) => {
				this.table_data = result.data
				this.total = result.total
				this.online_count = result.online_count
				this.ip_total = result.ip_total
			})
		},
		ExportUser() {
			let data = this.getQueryData()
			this.$download('/v1/get_user', data).then((result) => {
				this.$message.success('导出成功')
			})
		},
		ModifyItem(index, next, item) {
			if (index == 0) return next('编辑会员')
			let reqdata = {
				id: item.id,
			}
			if (index == 1) {
				reqdata.chat_state = 1
			} else if (index == 2) {
				reqdata.chat_state = 2
			} else if (index == 3) {
				reqdata.state = 2
			} else if (index == 4) {
				reqdata.state = 1
			}
			this.$post('/v1/update_user', reqdata, { google: true }).then(() => {
				this.$message.success('修改成功')
				this.getTableData()
			})
		},
		AddItem(index, next) {
			if (index == 0) return next('添加会员')
			if (index == 1) return next('批量导入')
		},
		DownLoadExcelTemplate() {
			const data = XLSX.utils.aoa_to_sheet([['账号', '密码']])
			const wb = XLSX.utils.book_new()
			XLSX.utils.book_append_sheet(wb, data, 'kalacloud-data')
			XLSX.writeFile(wb, `导入会员模板_${moment().format('YYYY-MM-DD HH-mm-ss')}.xlsx`)
		},
	},
}
</script>
