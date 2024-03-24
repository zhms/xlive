<template>
	<div class="container">
		<div class="query">
			<el-form :inline="true" :model="filters">
				<el-form-item label="">
					<el-date-picker v-model="filters.DateTimeRange" type="datetimerange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" clearable> </el-date-picker>
				</el-form-item>
				<el-form-item>
					<el-button type="primary" v-on:click="handleQuery">查询</el-button>
				</el-form-item>
			</el-form>
		</div>
		<div class="chart">
			<line-chart :data="data" />
		</div>
	</div>
</template>

<script>
import base from '@/api/base'
import * as echarts from 'echarts'
import moment from 'moment'
require('echarts/theme/shine')

export default {
	extends: base,
	data() {
		return {
			filters: {
				DateTimeRange: [moment().format('YYYY-MM-DD 00:00:00'), moment().format('YYYY-MM-DD 23:59:59')],
			},
			chartLine: null,
			data: [
				{
					name: '峰值人数',
					data: [
						// {
						// 	label: '2023-10-11 00:00:01',
						// 	value: 100,
						// },
						// {
						// 	label: '2023-10-11 00:00:02',
						// 	value: 100,
						// },
						// {
						// 	label: '2023-10-11 00:00:03',
						// 	value: 100,
						// },
						// {
						// 	label: '2023-10-11 00:00:04',
						// 	value: 100,
						// },
						// {
						// 	label: '2023-10-11 00:00:05',
						// 	value: 100,
						// },
						// {
						// 	label: '2023-10-11 00:00:06',
						// 	value: 100,
						// },
						// {
						// 	label: '2023-10-11 00:00:07',
						// 	value: 100,
						// },
					],
				},
			],
		}
	},
	mounted() {},
	methods: {},
	created() {
		this.getTableData()
	},
	methods: {
		getTableData() {
			let data = this.getQueryData()
			this.$post('/v1/get_online_data', data).then((result) => {
				this.data[0].data = []
				for (let i = result.data.length - 1; i >= 0; i--) {
					this.data[0].data.push({
						label: moment(result.data[i].create_time).format('YYYY-MM-DD HH:mm:ss'),
						value: result.data[i].v1,
					})
				}
			})
		},
	},
}
</script>

<style scope>
.query {
	margin-left: 20px;
	margin-top: 20px;
}
.chart {
	width: 100%;
	height: 600px;
}
</style>
