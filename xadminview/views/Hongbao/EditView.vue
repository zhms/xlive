<template>
	<div class="dialogBox">
		<el-dialog style="margin-top: -100px" :title="title" :visible.sync="visable" width="530px" center @open="handleOpen">
			<el-form :inline="true" label-width="120px">
				<el-form-item label="红包金额:">
					<el-input v-model.number="itemdata.total_amount" style="width: 300px"></el-input>
				</el-form-item>
				<el-form-item label="红包个数:">
					<el-input v-model.number="itemdata.total_count" style="width: 300px"></el-input>
				</el-form-item>
				<el-form-item label="房间号:">
					<el-select v-model="itemdata.room_id" placeholder="请选择" style="width: 300px">
						<el-option v-for="item in room_id" :key="item.id" :label="item.value" :value="item.id"></el-option>
					</el-select>
				</el-form-item>
				<el-form-item label="备注:">
					<el-input type="textarea" v-model="itemdata.memo" :rows="4" style="width: 300px"></el-input>
				</el-form-item>
			</el-form>
			<span slot="footer" class="dialog-footer">
				<el-button type="primary" @click="handleCommit">确定</el-button>
			</span>
		</el-dialog>
	</div>
</template>
<script>
import dlgbase from '@/api/dlgbase'
export default {
	extends: dlgbase,
	data() {
		return {
			room_id: [],
		}
	},
	methods: {
		commitData(next) {
			this.$post('/v1/hongbao/send_hongbao', this.itemdata, { google: true }).then(() => {
				next(true)
			})
		},
		onOpen() {
			this.$post('/v1/live_room/get_live_room_id', {}).then((result) => {
				this.room_id = []
				for (let i = 0; i < result.ids.length; i++) {
					this.room_id.push({ id: result.ids[i], value: result.ids[i] })
				}
			})
		},
	},
}
</script>

<style scoped></style>
