<template>
	<div class="dialogBox">
		<el-dialog style="margin-top: -100px" :title="title" :visible.sync="visable" width="630px" center @open="handleOpen">
			<div v-if="title != '批量导入机器人'">
				<el-form :inline="true" label-width="120px">
					<el-form-item label="账号:">
						<el-input v-model="itemdata.account" style="width: 400px"></el-input>
					</el-form-item>
				</el-form>
				<div class="dlg-footer">
					<el-button type="primary" @click="handleCommit">确定</el-button>
				</div>
			</div>
			<div v-if="title == '批量导入机器人'" class="upload-container">
				<el-upload ref="uploadRef" :show-file-list="false" action="/upload" :on-change="onUpload" :auto-upload="false" :accept="'.xlsx,.xls'">
					<i class="el-icon-upload"></i>
					<div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
					<div class="el-upload__tip" slot="tip">只能上传 xls/xlsx 文件，且不超过 2M.</div>
				</el-upload>
			</div>
		</el-dialog>
	</div>
</template>
<script>
import dlgbase from '@/api/dlgbase'
import XLSX from 'xlsx'
export default {
	extends: dlgbase,
	data() {
		return {}
	},
	methods: {
		commitData(next) {
			if (this.title == '编辑机器人') {
				if (!this.itemdata.account) return this.$message.error('请填写账号')
				let data = JSON.parse(JSON.stringify(this.itemdata))
				this.$post('/v1/update_robot', data, { google: true }).then(() => {
					this.$message.success('修改成功')
					next(true)
				})
			}
			if (this.title == '添加机器人') {
				if (!this.itemdata.account) return this.$message.error('请填写账号')
				let data = JSON.parse(JSON.stringify(this.itemdata))
				this.$post('/v1/create_robot', data, { google: true }).then(() => {
					this.$message.success('添加成功')
					next(true)
				})
			}
		},
		onOpen() {},
		onUpload(file) {
			const reader = new FileReader()
			reader.onload = async (event) => {
				const data = new Uint8Array(event.target.result)
				const workbook = XLSX.read(data, { type: 'array' })
				const sheetName = workbook.SheetNames[0]
				const sheet = workbook.Sheets[sheetName]
				const jsonData = XLSX.utils.sheet_to_json(sheet, { header: 1 })
				for (let i = 1; i < jsonData.length; i++) {
					let data = {
						account: `${jsonData[i][0]}`,
					}
					if (data.account.length > 0) {
						await this.$post('/v1/create_robot', data, { google: false })
					}
				}
				this.$message.success('上传成功')
				this.visable = false
				this.$emit('getTableData')
			}
			reader.readAsArrayBuffer(file.raw)
			return false
		},
	},
}
</script>

<style scoped>
.upload-container {
	display: flex;
	justify-content: center;
	align-items: center;
}
.dlg-footer {
	display: flex;
	justify-content: center;
	align-items: center;
}
</style>
