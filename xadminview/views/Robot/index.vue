<template>
	<div class="container">
		<el-form :inline="true" :model="filters">
			<el-form-item label="">
				<el-input v-model="filters.account" placeholder="账号" style="width: 150px" :clearable="true"></el-input>
			</el-form-item>
			<el-form-item>
				<el-button type="primary" icon="el-icon-refresh" v-on:click="handleQuery">查询</el-button>
				<el-button type="primary" icon="el-icon-plus" class="mr10" @click="handleAdd(0)">添加</el-button>
				<el-button type="primary" icon="el-icon-upload2" class="mr10" @click="handleAdd(1)">批量导入</el-button>
				<el-button type="primary" icon="el-icon-download" class="mr10" @click="DownLoadExcelTemplate()">下载导入模板</el-button>
			</el-form-item>
		</el-form>
		<el-row>
			<el-col :span="24" style="margin-top: -13px">
				<div class="grid-content">
					<span style="margin-right: 30px">当前上线机器人 : {{ robot_count }}</span>
					<el-button type="primary" icon="el-icon-arrow-up" @click="setRobotCount">上线机器人</el-button>
				</div>
			</el-col>
		</el-row>
		<el-table :data="table_data" style="margin-top: -17px" border class="table" max-height="670px" :cell-style="{ padding: '0px' }" :highlight-current-row="true">
			<el-table-column align="center" prop="id" label="Id" width="100"></el-table-column>
			<el-table-column align="center" prop="account" label="账号" width="200">
				<template slot-scope="scope">
					<span style="cursor: pointer" class="blue" @click="handleEdit(scope.row, 0)">{{ scope.row.account }}</span>
				</template>
			</el-table-column>
			<el-table-column label="操作" align="left" width="300">
				<template slot-scope="scope">
					<el-button type="text" icon="el-icon-edit" @click="handleEdit(scope.row, 0)">编辑</el-button>
					<el-button type="text" class="red" icon="el-icon-edit" @click="handleDelete(scope.row)">删除</el-button>
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
			robot_count: 0,
		}
	},
	created() {
		this.getTableData()
	},
	methods: {
		getTableData() {
			let data = this.getQueryData()
			this.$post('/v1/get_robot', data).then((result) => {
				this.table_data = result.data
				this.total = result.total
			})
			this.$post('/v1/get_robot_count', data, { noloading: true }).then((result) => {
				this.robot_count = result.data.robot_count
			})
		},
		setRobotCount() {},
		ModifyItem(index, next, item) {
			if (index == 0) return next('编辑机器人')
		},
		AddItem(index, next) {
			if (index == 0) return next('添加机器人')
			if (index == 1) return next('批量导入机器人')
		},
		DeleteItem(item) {
			this.$post('/v1/delete_robot', item, { google: true }).then(() => {
				this.$message.success('删除成功')
				this.getTableData()
			})
		},
		DownLoadExcelTemplate() {
			const data = XLSX.utils.aoa_to_sheet([['账号']])
			const wb = XLSX.utils.book_new()
			XLSX.utils.book_append_sheet(wb, data, 'kalacloud-data')
			XLSX.writeFile(wb, `导入机器人模板_${moment().format('YYYY-MM-DD HH-mm-ss')}.xlsx`)
		},
	},
}
</script>

<style>
.el-row {
	margin-bottom: 20px;
}
.el-col {
	border-radius: 4px;
}
.bg-purple-dark {
	background: #99a9bf;
}
.bg-purple {
	background: #d3dce6;
}
.bg-purple-light {
	background: #e5e9f2;
}
.grid-content {
	border-radius: 4px;
	min-height: 36px;
}
.row-bg {
	padding: 10px 0;
	background-color: #f9fafc;
}
</style>
