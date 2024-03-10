<template>
	<div class="container">
		<el-form :inline="true" :model="filters">
			<el-form-item>
				<el-button type="primary" icon="el-icon-refresh" size="small" v-on:click="handleQuery">刷新</el-button>
				<el-button type="primary" icon="el-icon-plus" size="small" v-on:click="handleAdd(0)">发红包</el-button>
			</el-form-item>
		</el-form>
		<el-table :data="table_data" style="margin-top: -15px" border class="table" max-height="670px" :cell-style="{ padding: '0px' }" :highlight-current-row="true">
			<el-table-column align="center" prop="id" label="id" width="100"></el-table-column>
			<el-table-column align="center" prop="room_id" label="房间Id" width="100"></el-table-column>
			<el-table-column align="center" prop="total_amount" label="红包金额" width="150"></el-table-column>
			<el-table-column align="center" prop="total_count" label="红包个数" width="150"> </el-table-column>
			<el-table-column align="center" prop="used_amount" label="已领金额" width="150"></el-table-column>
			<el-table-column align="center" prop="used_count" label="已领个数" width="150"> </el-table-column>
			<el-table-column align="center" prop="sender" label="发布人" width="150"> </el-table-column>
			<el-table-column align="center" prop="memo" label="备注" width="200"> </el-table-column>
			<el-table-column align="center" prop="create_time" label="发送时间" width="150"> </el-table-column>
			<el-table-column label="操作">
				<template slot-scope="scope">
					<el-button type="text" size="small" icon="el-icon-search" @click="handleEdit(scope.row, 0)">详情</el-button>
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
			this.$post('/v1/hongbao/get_hongbao_list', {}).then((result) => {
				this.table_data = this.dealData(result.data)
				this.total = result.total
			})
		},
		AddItem(index, next) {
			if (index == 0) return next('发红包')
		},
	},
}
</script>
