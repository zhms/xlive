<template>
	<div class="container">
		<el-form :inline="true" :model="filters">
			<el-form-item label="运营商:" v-show="zong">
				<el-select v-model="filters.seller_id" placeholder="运营商" style="width: 150px" clearable @change="sellerChange">
					<el-option v-for="item in sellers" :key="item.seller_id" :label="item.SellerName" :value="item.seller_id"> </el-option>
				</el-select>
			</el-form-item>
			<el-form-item label="">
				<el-input v-model="filters.account" placeholder="管理员" style="width: 150px" size="small" :clearable="true"></el-input>
			</el-form-item>
			<el-form-item label="">
				<el-input v-model="filters.opt_name" placeholder="操作" style="width: 150px" size="small" :clearable="true"></el-input>
			</el-form-item>
			<el-form-item label="">
				<el-date-picker v-model="filters.DateRange" size="small" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" clearable> </el-date-picker>
			</el-form-item>
			<el-form-item>
				<el-button type="primary" icon="el-icon-refresh" size="small" @click="handleQuery">查询</el-button>
			</el-form-item>
		</el-form>
		<el-table style="margin-top: -15px" :data="table_data" border class="table" max-height="670px" :cell-style="{ padding: '3px' }" :highlight-current-row="true">
			<el-table-column align="center" prop="id" label="序号" width="80"></el-table-column>
			<el-table-column align="center" prop="account" label="管理员" width="100"></el-table-column>
			<el-table-column align="center" prop="SellerName" label="运营商" width="100" v-if="zong">
				<template slot-scope="scope">
					<span>{{ getSellerName(scope.row) }}</span>
				</template>
			</el-table-column>
			<el-table-column align="center" prop="opt_name" label="操作" width="200"></el-table-column>
			<el-table-column align="center" prop="req_ip" label="ip" width="130"></el-table-column>
			<el-table-column align="center" prop="create_time" label="时间" width="200"></el-table-column>
			<el-table-column label="内容">
				<template slot-scope="scope">
					<span type="text" style="cursor: pointer" @click="copy(scope.row.req_data)"> {{ scope.row.req_data }}</span>
				</template>
			</el-table-column>
		</el-table>
		<div class="pagination">
			<el-pagination style="margin-right: 10px" background layout="sizes,total, prev, pager, next, jumper" @size-change="page_sizeChange" :current-page="page" @current-change="pageChange" :page-sizes="page_sizes" :total="total" :page-size="page_size"></el-pagination>
		</div>
	</div>
</template>
<script>
import '@/styles/main.css'
import base from '@/api/base'
import moment from 'moment'
export default {
	extends: base,
	data() {
		return {
			filters: {
				DateRange: [moment().format('YYYY-MM-DD'), moment().format('YYYY-MM-DD')],
			},
		}
	},
	created() {
		this.getTableData()
	},
	methods: {
		getTableData() {
			let data = this.getQueryData()
			data.channel_id = Number(data.channel_id)
			this.$post('/v1/admin_log/get_opt_log', data).then((result) => {
				this.table_data = this.dealData(result.data)
				this.total = result.total
			})
		},
	},
}
</script>
