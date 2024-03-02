<template>
	<div class="container">
		<el-form :inline="true" :model="filters">
			<!-- <el-form-item label="">
				<el-input v-model="filters.user_id" placeholder="会员Id" style="width: 150px" :clearable="true"></el-input>
			</el-form-item>
			<el-form-item label="">
				<el-input v-model="filters.account" placeholder="账号" style="width: 150px" :clearable="true"></el-input>
			</el-form-item>
			<el-form-item label="">
				<el-input v-model="filters.agent" placeholder="业务员" style="width: 150px" :clearable="true"></el-input>
			</el-form-item>
			<el-form-item label="">
				<el-input v-model="filters.login_ip" placeholder="登录Ip" style="width: 150px" :clearable="true"></el-input>
			</el-form-item>
			<el-form-item>
				<el-button type="primary" icon="el-icon-refresh" v-on:click="handleQuery">查询</el-button>
				<el-button type="primary" icon="el-icon-plus" class="mr10" @click="handleAdd(0)">添加</el-button>
				<el-button type="primary" icon="el-icon-upload2" class="mr10" @click="handleAdd(0)">批量导入</el-button>
				<el-button type="primary" icon="el-icon-bottom" class="mr10" @click="handleAdd(0)">导出</el-button>
				<el-button type="primary" icon="el-icon-download" class="mr10" @click="DownLoadExcelTemplate()">下载导入模板</el-button>
			</el-form-item> -->
		</el-form>
		<el-table :data="table_data" border class="table" max-height="670px" :cell-style="{ padding: '0px' }" :highlight-current-row="true">
			<el-table-column align="center" prop="id" label="房间Id" width="100"></el-table-column>
			<el-table-column align="center" prop="name" label="房间名称" width="200">
				<template slot-scope="scope">
					<span style="cursor: pointer" class="blue" @click="handleEdit(scope.row, 0)">{{ scope.row.name }}</span>
				</template>
			</el-table-column>
			<el-table-column align="center" prop="account" label="讲师名称" width="200"> </el-table-column>
			<el-table-column align="center" label="状态" width="100">
				<template slot-scope="scope">
					<span :class="scope.row.state == 1 ? 'blue' : 'red'">{{ scope.row.state == 1 ? '直播中' : '未开启' }}</span>
				</template>
			</el-table-column>
			<el-table-column align="center" prop="push_url" label="推流域名" width="100">
				<template slot-scope="scope">
					<span style="cursor: pointer" class="blue" @click="copy(scope.row.push_url)">{{ '复制' }}</span>
				</template>
			</el-table-column>
			<el-table-column align="center" prop="pull_url" label="拉流域名" width="100">
				<template slot-scope="scope">
					<span style="cursor: pointer" class="blue" @click="copy(scope.row.pull_url)">{{ '复制' }}</span>
				</template>
			</el-table-column>
			<el-table-column align="center" prop="live_url" label="直播间域名" width="100">
				<template slot-scope="scope">
					<span style="cursor: pointer" class="blue" @click="copy(scope.row.live_url)">{{ '复制' }}</span>
				</template>
			</el-table-column>
			<el-table-column align="center" prop="create_time" label="创建时间" width="200"> </el-table-column>

			<el-table-column label="操作" align="left" width="300">
				<template slot-scope="scope">
					<el-button type="text" size="small" icon="el-icon-edit" @click="handleEdit(scope.row, 0)">编辑</el-button>
					<el-button type="text" size="small" class="red" icon="el-icon-edit" @click="handleDelete(scope.row)">删除</el-button>
				</template>
			</el-table-column>
		</el-table>
		<div class="pagination">
			<el-pagination style="margin-right: 10px" background layout="sizes,total, prev, pager, next, jumper" @size-change="page_sizeChange" :current-page="page" @current-change="pageChange" :page-sizes="page_sizes" :total="total" :page-size="page_size"></el-pagination>
		</div>
		<edit-view :show.sync="dialog0.show" :title="dialog0.title" :itemdata="dialog0.itemdata" :filters="filters" @getTableData="getTableData" />
		<div>
			<el-dialog :title="dialog1.title" :visible.sync="dialog1.show" width="350px" center>
				<vueqr :text="dialog1.url" :size="300"> </vueqr>
			</el-dialog>
		</div>
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
			this.$post('/v1/live_room/get_live_room', data).then((result) => {
				this.table_data = this.dealData(result.data)
				this.total = result.total
			})
		},
		ModifyItem(index, next, item) {
			if (index == 0) next('编辑直播间')
		},
		AddItem(index, next) {
			if (index == 0) next('添加直播间')
		},
		DeleteItem(item) {
			let data = {
				id: item.id,
			}
			this.$delete('/v1/live_room/get_live_room', data, { google: true }).then(() => {
				this.$message.success('删除成功')
				this.getTableData()
			})
		},
	},
}
</script>
