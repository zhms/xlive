<template>
	<div class="container">
		<el-form :inline="true" :model="filters">
			<!-- <el-form-item label="">
				<el-input v-model.number="filters.room_id" placeholder="房间Id" style="width: 150px" :clearable="true"></el-input>
			</el-form-item> -->
			<el-form-item>
				<el-button type="primary" icon="el-icon-refresh" v-on:click="handleQuery">查询</el-button>
			</el-form-item>
		</el-form>
		<el-table :data="table_data" style="margin-top: -13px" border class="table" max-height="670px" :cell-style="{ padding: '0px' }" :highlight-current-row="true">
			<el-table-column align="center" prop="id" label="id" width="100"></el-table-column>
			<el-table-column align="center" prop="ip" label="Ip" width="100"></el-table-column>
			<el-table-column align="center" prop="admin_account" label="封禁人" width="200"> </el-table-column>
			<el-table-column align="center" prop="create_time" label="封禁时间" width="200">
				<template slot-scope="scope">
					<span>{{ scope.row.create_time | 北京时间 }}</span>
				</template>
			</el-table-column>
			<el-table-column label="操作" align="left" width="300">
				<template slot-scope="scope">
					<el-button type="text" icon="el-icon-edit" @click="handleDelete(scope.row)">解封</el-button>
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
export default {
	extends: base,
	data() {
		return {}
	},
	created() {
		this.getTableData()
	},
	methods: {
		getTableData() {
			let data = this.getQueryData()
			this.$post('/v1/get_ip_ban', data).then((result) => {
				this.table_data = result.data
				this.total = result.total
			})
		},
		DeleteItem(item) {
			this.$post('/v1/delete_ip_ban', item, { google: true }).then((result) => {
				this.getTableData()
			})
		},
	},
}
</script>
