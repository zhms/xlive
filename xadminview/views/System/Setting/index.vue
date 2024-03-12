<template>
	<div class="container">
		<el-form :inline="true" :model="filters">
			<el-form-item label="">
				<el-input v-model="filters.ConfigName" style="width: 200px" clearable placeholder="配置名"></el-input>
			</el-form-item>
			<el-form-item label="">
				<el-input v-model="filters.Memo" style="width: 200px" clearable placeholder="说明模糊查询"></el-input>
			</el-form-item>
			<el-form-item>
				<el-button type="primary" icon="el-icon-refresh" v-on:click="handleQuery">查询</el-button>
				<el-button type="primary" icon="el-icon-plus" v-on:click="handleAdd(0)">添加</el-button>
			</el-form-item>
		</el-form>
		<el-table :data="table_data" border class="table" max-height="670px" :cell-style="{ padding: '0px' }" :highlight-current-row="true">
			<el-table-column align="center" prop="Id" label="序号" width="50" v-if="column['序号']"> </el-table-column>
			<el-table-column align="center" label="前端" width="200" v-if="column['前端']">
				<template slot-scope="scope">
					<span>{{ MapYesNo[scope.row.ForClient] }}</span>
				</template>
			</el-table-column>
			<el-table-column align="center" prop="ConfigName" label="配置名" width="200" v-if="column['配置名']"></el-table-column>
			<el-table-column align="center" prop="ConfigValue" label="配置值" width="400" v-if="column['配置值']"></el-table-column>
			<el-table-column align="center" prop="Memo" label="说明" width="400" v-if="column['说明']"></el-table-column>
			<el-table-column label="操作" v-if="column['操作']">
				<template slot-scope="scope">
					<el-button type="text" size="small" icon="el-icon-edit" @click="handleEdit(scope.row, 0)">编辑</el-button>
				</template>
			</el-table-column>
		</el-table>
		<div class="pagination">
			<el-pagination style="margin-right: 10px" background layout="sizes,total, prev, pager, next, jumper" @size-change="page_sizeChange" :current-page="page" @current-change="pageChange" :page-sizes="page_sizes" :total="total" :page-size="page_size"></el-pagination>
		</div>
		<column-picker :show.sync="columnpicker" :columns="columns" :column="column" @setColumn="setColumn" />
		<edit-view :show.sync="dialog0.show" :title="dialog0.title" :itemdata="dialog0.itemdata" :filters="filters" @getTableData="getTableData" />
	</div>
</template>
<script>
import '@/styles/main.css'
import base from '@/api/base'
import EditView from './EditView.vue'
import ColumnPicker from '@/layout/columnpicker.vue'
export default {
	extends: base,
	components: { EditView, ColumnPicker },
	data() {
		return {
			columns: ['序号', '渠道', '前端', '配置名', '配置值', '说明', '操作'],
		}
	},
	created() {
		this.getTableData()
	},
	methods: {
		getTableData() {
			// let data = this.getQueryData()
			// data.channel_id = Number(data.channel_id)
			// if (data.ConfigName != null || data.ConfigName != undefined) {
			// 	if (data.ConfigName != '') {
			// 		data.ConfigName = [data.ConfigName]
			// 	} else {
			// 		delete data.ConfigName
			// 	}
			// }
			// this.$post('/get_system_config', data).then((result) => {
			// 	this.table_data = result.data
			// 	this.total = result.total
			// })
		},
		ModifyItem(index, next, item) {
			if (index == 0) return next('编辑配置')
		},
		AddItem(index, next) {
			if (index == 0) return next('添加配置', { ForClient: 1 })
		},
	},
}
</script>
