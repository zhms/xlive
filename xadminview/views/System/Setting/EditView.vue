<template>
	<div class="dialogBox">
		<el-dialog :title="title" :visible.sync="visable" width="600px" center @open="handleOpen">
			<el-form :inline="true" label-width="100px">
				<!-- <el-form-item label="渠道:" v-if="itemdata.channel_id != 0">
					<el-select v-model="itemdata.channel_id" placeholder="渠道" style="width: 200px" :disabled="title == '编辑配置'" clearable>
						<el-option v-for="item in channels" :key="item.channel_id" :label="item.channel_name" :value="item.channel_id"> </el-option>
					</el-select>
				</el-form-item> -->
				<el-form-item label="前端:">
					<el-select v-model="itemdata.ForClient" placeholder="前端" style="width: 200px">
						<el-option v-for="item in ListYesNo" :key="item.id" :label="item.name" :value="item.id"> </el-option>
					</el-select>
				</el-form-item>
				<el-form-item label="配置名:">
					<el-input v-model="itemdata.ConfigName" style="width: 400px" :disabled="title == '编辑配置'"></el-input>
				</el-form-item>
				<el-form-item label="配置值:">
					<el-input type="textarea" v-model="itemdata.ConfigValue" style="width: 400px" :rows="6"></el-input>
				</el-form-item>
				<el-form-item label="说明:">
					<el-input type="textarea" v-model="itemdata.Memo" style="width: 400px" :rows="8"></el-input>
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
		return {}
	},
	methods: {
		onOpen() {},
		commitData(next) {
			if (this.title == '编辑配置') {
				this.itemdata.channel_id = Number(this.itemdata.channel_id)
				this.itemdata.ForClient = Number(this.itemdata.ForClient)
				let reqdata = {
					Config: [
						{
							channel_id: this.itemdata.channel_id,
							ConfigName: this.itemdata.ConfigName,
							ConfigValue: this.itemdata.ConfigValue,
							ForClient: this.itemdata.ForClient,
							Memo: this.itemdata.Memo,
						},
					],
				}
				this.$post('/modify_system_config', reqdata, { google: true }).then(() => {
					next(true)
				})
			}
			if (this.title == '添加配置') {
				this.itemdata.channel_id = Number(this.itemdata.channel_id)
				this.itemdata.ForClient = Number(this.itemdata.ForClient)
				this.$post('/add_system_config', this.itemdata, { google: true }).then(() => {
					next(true)
				})
			}
		},
	},
}
</script>
