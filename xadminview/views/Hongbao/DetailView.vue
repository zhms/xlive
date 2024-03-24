<template>
	<div class="dialogBox">
		<el-dialog style="margin-top: -100px" :title="title" :visible.sync="visable" width="660px" center @open="handleOpen">
			<el-table :data="table_data" style="margin-top: -13px" border class="table" :cell-style="{ padding: '0px' }" :highlight-current-row="true">
				<el-table-column align="center" prop="id" label="id" width="100"></el-table-column>
				<el-table-column align="center" prop="amount" label="金额" width="100">
					<template slot-scope="scope">
						<span>{{ scope.row.amount | amount2 }}</span>
					</template>
				</el-table-column>
				<el-table-column align="center" prop="account" label="收款人" width="200"></el-table-column>
				<el-table-column align="center" prop="account" label="时间" width="200">
					<template slot-scope="scope">
						<span>{{ scope.row.create_time | 北京时间 }}</span>
					</template>
				</el-table-column>
			</el-table>
		</el-dialog>
	</div>
</template>
<script>
import dlgbase from '@/api/dlgbase'
import { amount2 } from '@/api/filters'
export default {
	extends: dlgbase,
	data() {
		return {
			table_data: [],
		}
	},
	methods: {
		onOpen() {
			this.title = `红包详情 | 个数:${this.itemdata.total_count}(${this.itemdata.used_count}) 金额:${this.itemdata.total_amount}(${amount2(this.itemdata.used_amount)})`
			this.$post('/v1/get_hongbao_detail', this.itemdata).then((result) => {
				// let test = result.data[0]
				// let f = JSON.stringify(test)
				// for (let i = 0; i < 200; i++) {
				// 	this.table_data.push(JSON.parse(f))
				// }
				this.table_data = result.data
			})
		},
	},
}
</script>

<style scoped></style>
