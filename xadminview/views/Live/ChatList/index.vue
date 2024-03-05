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
		<el-table :data="table_data" border class="table" max-height="670px" :cell-style="{ padding: '0px' }" :highlight-current-row="true">
			<el-table-column align="center" prop="id" label="id" width="100"></el-table-column>
			<el-table-column align="center" prop="user_id" label="会员Id" width="100"></el-table-column>
			<el-table-column align="center" prop="account" label="账号" width="200"> </el-table-column>
			<el-table-column align="center" prop="content" label="内容" width="300"></el-table-column>
			<el-table-column align="center" prop="ip" label="ip" width="120"></el-table-column>
			<el-table-column align="center" prop="ip_location" label="ip地区" width="150"></el-table-column>
			<el-table-column label="操作" align="left" width="300">
				<template slot-scope="scope">
					<el-button type="text" size="small" icon="el-icon-edit" @click="handleEdit(scope.row, 0)" v-if="scope.row.state == 1">通过</el-button>
					<el-button type="text" size="small" class="red" icon="el-icon-edit" @click="handleEdit(scope.row, 1)" v-if="scope.row.state == 1">拒绝</el-button>
					<el-button type="text" size="small" class="red" icon="el-icon-edit" @click="handleEdit(scope.row, 3)" v-if="scope.row.state == 1">封ip</el-button>
					<el-button type="text" size="small" class="red" icon="el-icon-edit" @click="handleEdit(scope.row, 2)" v-if="scope.row.state == 1">封号</el-button>
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
			data.room_id = Number(this.filters.room_id) || 0
			this.$post('/v1/live_chat/get_live_chat', data).then((result) => {
				this.table_data = this.dealData(result.data)
				this.total = result.total
			})
		},
		ModifyItem(index, next, item) {
			if (index == 0) {
				let data = JSON.parse(JSON.stringify(item))
				data.state = 2
				this.$patch('/v1/live_chat/audit_live_chat', data).then((result) => {})
			}
		},
	},
}
</script>
