<template>
	<div class="container">
		<el-form :inline="true" :model="filters">
			<el-form-item label="">
				<el-input v-model.number="filters.room_id" placeholder="房间Id" style="width: 150px" :clearable="true"></el-input>
			</el-form-item>
			<el-form-item>
				<el-button type="primary" icon="el-icon-refresh" v-on:click="handleQuery">查询</el-button>
			</el-form-item>
		</el-form>
		<el-table :data="table_data" style="margin-top: -13px" border class="table" max-height="670px" :cell-style="{ padding: '0px' }" :highlight-current-row="true">
			<el-table-column align="center" prop="id" label="id" width="100"></el-table-column>
			<el-table-column align="center" prop="room_id" label="房间Id" width="100"></el-table-column>
			<el-table-column align="center" prop="account" label="账号" width="200"> </el-table-column>
			<el-table-column align="center" prop="content" label="内容" width="500" show-overflow-tooltip></el-table-column>
			<el-table-column align="center" label="状态" width="100">
				<template slot-scope="scope">
					<span>{{ scope.row.state == 1 ? '待审核' : scope.row.state == 2 ? '通过' : '拒绝' }}</span>
				</template>
			</el-table-column>
			<el-table-column align="center" prop="ip" label="Ip地址" width="120"></el-table-column>
			<el-table-column align="center" prop="ip_location" label="ip地区" width="150"></el-table-column>
			<el-table-column label="操作" align="left" width="300">
				<template slot-scope="scope">
					<el-button type="text" size="mini" icon="el-icon-edit" @click="handleEdit(scope.row, 2)" v-if="scope.row.state == 1">通过</el-button>
					<el-button type="text" size="mini" class="red" icon="el-icon-edit" @click="handleEdit(scope.row, 3)" v-if="scope.row.state == 1">拒绝</el-button>
					<el-button type="text" size="mini" class="red" icon="el-icon-edit" @click="handleEdit(scope.row, 4)" v-if="scope.row.state == 1">封ip</el-button>
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
			data.room_id = Number(this.filters.room_id) || 0
			this.$post('/v1/get_chat_data', data).then((result) => {
				this.table_data = result.data
				this.total = result.total
			})
		},
		ModifyItem(index, next, item) {
			let data = JSON.parse(JSON.stringify(item))
			data.state = index
			this.$post('/v1/update_chat_data', data).then((result) => {
				for (let i = 0; i < this.table_data.length; i++) {
					if (this.table_data[i].id == item.id) {
						this.table_data[i].state = index
						break
					}
				}
			})
		},
	},
}
</script>
