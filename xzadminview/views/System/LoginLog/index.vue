<template>
	<div class="container">
		<el-form :inline="true" :model="filters">
			<el-form-item label="">
				<el-input v-model="filters.account" placeholder="管理员" style="width: 150px" clearable></el-input>
			</el-form-item>
			<el-form-item label="">
				<el-input v-model="filters.login_ip" placeholder="登录Ip" style="width: 150px" clearable></el-input>
			</el-form-item>
			<el-form-item label="">
				<el-date-picker v-model="filters.DateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" clearable> </el-date-picker>
			</el-form-item>
			<el-form-item>
				<el-button type="primary" icon="el-icon-refresh" v-on:click="handleQuery">查询</el-button>
			</el-form-item>
		</el-form>
		<el-table :data="table_data" style="margin-top: -13px" border class="table" max-height="670px" :cell-style="{ padding: '3px' }" :highlight-current-row="true">
			<el-table-column align="center" prop="id" label="序号" width="80" v-if="column['序号']"></el-table-column>
			<el-table-column align="center" prop="account" label="管理员" width="150" v-if="column['管理员']"></el-table-column>
			<el-table-column align="center" prop="login_ip" label="登录IP" width="200" v-if="column['登录Ip']"></el-table-column>
			<el-table-column align="center" prop="login_ip_location" label="登录地区" width="300"></el-table-column>
			<el-table-column align="center" prop="create_time" label="登录时间" width="170" v-if="column['登录时间']">
				<template slot-scope="scope">
					<span>{{ scope.row.create_time | 北京时间 }}</span>
				</template>
			</el-table-column>
		</el-table>
		<div class="pagination">
			<el-pagination style="margin-right: 10px" background layout="sizes,total, prev, pager, next, jumper" @size-change="page_sizeChange" :current-page="page" @current-change="pageChange" :page-sizes="page_sizes" :total="total" :page-size="page_size"></el-pagination>
		</div>
		<column-picker :show.sync="columnpicker" :columns="columns" :column="column" @setColumn="setColumn" />
	</div>
</template>
<script>
import '@/styles/main.css'
import base from '@/api/base'
import ColumnPicker from '@/layout/columnpicker.vue'
import moment from 'moment'
export default {
	extends: base,
	data() {
		return {
			columns: ['序号', '渠道', '管理员', '登录地区', '登录Ip', '登录时间'],
			filters: {
				DateRange: [moment().format('YYYY-MM-DD'), moment().format('YYYY-MM-DD')],
			},
		}
	},
	created() {
		this.getTableData()
	},
	components: { ColumnPicker },
	methods: {
		getTableData() {
			let data = this.getQueryData()
			this.$post('/v1/admin_get_login_log', data).then((result) => {
				this.table_data = result.data
				this.total = result.total
			})
		},
	},
}
</script>
