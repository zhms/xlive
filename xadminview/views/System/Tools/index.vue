<template>
	<div class="container">
		<div style="width: 700px">
			<el-form label-width="120px">
				<el-form-item label="数据库表名:">
					<el-button type="primary" size="small" @click="copy_db_table">生成</el-button>
				</el-form-item>
				<el-form-item label="数据库字段名:">
					<el-button type="primary" size="small" @click="copy_db_field">生成</el-button>
				</el-form-item>
				<el-form-item label="生成响应函数:">
					<el-input v-model="funname" size="small" style="width: 200px; padding-right: 10px"></el-input>
					<el-button type="primary" size="small" @click="copy_method(0)">原始</el-button>
					<el-button type="primary" size="small" @click="copy_method(1)">获取</el-button>
					<el-button type="primary" size="small" @click="copy_method(2)">创建</el-button>
					<el-button type="primary" size="small" @click="copy_method(3)">更新</el-button>
					<el-button type="primary" size="small" @click="copy_method(4)">删除</el-button>
				</el-form-item>
			</el-form>
		</div>
	</div>
</template>
<script>
import '@/styles/main.css'
import base from '@/api/base'
export default {
	extends: base,
	data() {
		return {
			funname: '',
		}
	},
	created() {
		let userinfo = JSON.parse(sessionStorage.getItem('userinfo'))
		if (userinfo.env.indexOf('prd') >= 0) {
			this.$router.push('/404')
		}
	},
	methods: {
		copy_method(idx) {
			if (funname == '') return this.$message.error({ message: '请输入函数名', center: true })
			let res = ''
			let resex = ''
			let funname = this.funname
			if (idx == 1) {
				funname = 'get_' + funname
				res = `
type ${funname}_res struct {
}`
				resex = ` {object} ${funname}_res `
			}
			if (idx == 2) {
				funname = 'create_' + funname
			}
			if (idx == 3) {
				funname = 'update_' + funname
			}
			if (idx == 4) {
				funname = 'delete_' + funname
			}
			if (idx == 0) {
				funname = this.funname
				funname = funname
				res = `
type ${funname}_res struct {
}`
			}
			let methodstr = `type ${funname}_req struct {
}
${res}
// @Router /${funname} [post]
// @Tags a
// @Summary b
// @Param x-token header string true "token"
// @Param body body ${funname}_req true "请求参数"
// @Success 200 ${resex}"响应数据"
func ${funname}(ctx *gin.Context) {
    var reqdata ${funname}_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
    //response := new(${funname}_res)
	//ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
    //ctx.JSON(http.StatusOK, xenum.Success)
}`
			this.copy(methodstr)
		},

		copy_db_table() {
			this.$post('/v1/admin_tools', { query_type: 'db_table' }).then((result) => {
				this.copy(result.data)
			})
		},
		copy_db_field() {
			this.$post('/v1/admin_tools', { query_type: 'db_field' }).then((result) => {
				this.copy(result.data)
			})
		},
	},
}
</script>
