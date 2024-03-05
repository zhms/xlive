<template>
	<div class="container">
		<el-form :inline="true" :model="filters">
			<el-form-item>
				<el-button type="primary" icon="el-icon-plus" class="mr10" @click="handleAdd(0)">添加</el-button>
			</el-form-item>
		</el-form>
		<el-table :data="table_data" border class="table" max-height="670px" :cell-style="{ padding: '0px' }" :highlight-current-row="true">
			<el-table-column align="center" type="index" label="序号" width="100"></el-table-column>
			<el-table-column align="center" prop="role_name" label="角色名" width="200"></el-table-column>
			<el-table-column align="center" prop="parent" label="上级角色" width="200"></el-table-column>
			<el-table-column align="center" prop="memo" label="备注" width="300"></el-table-column>
			<el-table-column label="操作">
				<template slot-scope="scope">
					<el-button type="text" size="small" icon="el-icon-edit" @click="handleEdit(scope.row, 0)">编辑</el-button>
					<el-button type="text" size="small" icon="el-icon-delete" class="red" @click="handleDelete(scope.row, 0)">删除</el-button>
				</template>
			</el-table-column>
		</el-table>
		<div class="pagination">
			<el-pagination style="margin-right: 10px" background layout="sizes,total, prev, pager, next, jumper" @size-change="page_sizeChange" :current-page="page" @current-change="pageChange" :page-sizes="page_sizes" :total="total" :page-size="page_size"></el-pagination>
		</div>
		<edit-view :show.sync="dialog0.show" :title="dialog0.title" :itemdata="dialog0.itemdata" @getTableData="getTableData" />
	</div>
</template>
<script>
import '@/styles/main.css'
import base from '@/api/base'
import EditView from './EditView.vue'
export default {
	extends: base,
	components: { EditView },
	created() {
		this.getTableData()
	},
	methods: {
		getTableData() {
			let data = this.getQueryData()
			this.$post('/v1/admin_role/get_admin_role', data).then((result) => {
				console.log(result)
				this.table_data = result.data
				this.total = result.total
			})
		},
		AddItem(index, next) {
			if (index == 0) return next(`添加角色`, {})
		},
		ModifyItem(index, next, item) {
			if (index == 0 && item.Parent == 'god') return this.$message.error('该角色不可编辑')
			if (index == 0) next(`编辑角色`)
		},
		DeleteItem(item) {
			if (item.Parent == 'god') return this.$message.error('该角色不可删除')
			this.$delete('/v1/admin_role/delete_admin_role', item, { google: true }).then(() => {
				this.$message.success('删除成功')
				this.getTableData()
			})
		},
	},
}
</script>
