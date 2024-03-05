<template>
	<div class="container">
		<el-form :inline="true" :model="filters">
			<el-form-item>
				<el-button type="primary" icon="el-icon-plus" class="mr10" @click="handleAdd(0)">添加</el-button>
			</el-form-item>
		</el-form>
		<el-table :data="table_data" border class="table" max-height="670px" :cell-style="{ padding: '0px' }" :highlight-current-row="true">
			<el-table-column align="center" prop="id" label="房间Id" width="150"></el-table-column>
			<el-table-column align="center" prop="name" label="房间名称" width="150">
				<template slot-scope="scope">
					<span style="cursor: pointer" class="blue" @click="handleEdit(scope.row, 0)">{{ scope.row.name }}</span>
				</template>
			</el-table-column>
			<el-table-column align="center" prop="title" label="聊天标题" width="150"></el-table-column>
			<el-table-column align="center" prop="account" label="讲师名称" width="150"> </el-table-column>
			<el-table-column align="center" label="状态" width="100">
				<template slot-scope="scope">
					<span :class="scope.row.state == 1 ? 'blue' : 'red'">{{ scope.row.state == 1 ? '直播中' : '未开启' }}</span>
				</template>
			</el-table-column>
			<el-table-column align="center" prop="push_url" label="推流域名" width="100">
				<template slot-scope="scope">
					<span style="cursor: pointer" class="blue" @click="copy_push_url(scope.row)">{{ '复制' }}</span>
				</template>
			</el-table-column>
			<el-table-column align="center" prop="live_url" label="直播间域名" width="100">
				<template slot-scope="scope">
					<span style="cursor: pointer" class="blue" @click="copy_live_url(scope.row)">{{ '复制' }}</span>
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
		ModifyItem(index, next) {
			if (index == 0) return next('编辑直播间')
		},
		AddItem(index, next) {
			if (index == 0) {
				next('添加直播间', {
					name: '',
					title: '',
					account: '',
					state: 2,
				})
			}
		},
		DeleteItem(item) {
			let data = {
				id: item.id,
			}
			this.$delete('/v1/live_room/delete_live_room', data, { google: true }).then(() => {
				this.$message.success('删除成功')
				this.getTableData()
			})
		},
		copy_live_url(data) {
			if (data.state != 1) {
				this.$message.error('复制失败,直播间尚未开启')
				return
			}
			let url = data.live_url + '&s=' + JSON.parse(sessionStorage.getItem('userinfo') ?? '{}').account
			this.copy(url)
		},
		copy_push_url(data) {
			if (data.state != 1) {
				this.$message.error('复制失败,直播间尚未开启')
				return
			}
			this.copy(data.push_url)
		},
	},
}
</script>
