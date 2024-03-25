<template>
	<div class="container">
		<div style="width: 700px">
			<el-form label-width="120px">
				<el-form-item label="数据库表名:">
					<el-button type="primary" @click="copy_db_table">生成</el-button>
				</el-form-item>
				<el-form-item label="数据库字段名:">
					<el-button type="primary" @click="copy_db_field">生成</el-button>
				</el-form-item>
				<el-form-item label="响应函数:">
					<el-input v-model="funname" style="width: 200px; padding-right: 10px"></el-input>
					<el-button type="primary" @click="copy_method(0)">原始</el-button>
					<el-button type="primary" @click="copy_method(1)">获取</el-button>
					<el-button type="primary" @click="copy_method(2)">创建</el-button>
					<el-button type="primary" @click="copy_method(3)">更新</el-button>
					<el-button type="primary" @click="copy_method(4)">删除</el-button>
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
			if (this.funname == '') return this.$message.error({ message: '请输入函数名', center: true })
			if (idx == 1) {
				//get
				this.copy(`type get_${this.funname}_req struct {
	Page     int \`json:"page"\`      // 页码
	PageSize int \`json:"page_size"\` // 每页数量
}

type get_${this.funname}_res struct {
	Total int64           \`json:"total"\` // 总数
	Data  []*model. \`json:"data"\`  // 数据
}

// @Router /get_${this.funname} [post]
// @Tags a
// @Summary b
// @Param x-token header string true "token"
// @Param body body get_${this.funname}_req true "请求参数"
// @Success 200  {object} get_${this.funname}_res "响应数据"
func get_${this.funname}(ctx *gin.Context) {
	var reqdata get_${this.funname}_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
    if reqdata.Page == 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize == 0 {
		reqdata.PageSize = 15
	}
	response := new(get_${this.funname}_res)
	token := admin.GetToken(ctx)
	tb := xapp.DbQuery().
	itb := tb.WithContext(ctx).Order(tb.ID.Desc())
	itb = itb.Where(tb.SellerID.Eq(token.SellerId))
	{
    }
	var err error
	response.Data, response.Total, err = itb.FindByPage((reqdata.Page-1)*reqdata.PageSize, reqdata.PageSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}`)
			}
			if (idx == 2) {
				//create
				this.copy(`type create_${this.funname}_req struct {
}

// @Router /create_${this.funname} [post]
// @Tags a
// @Summary b
// @Param x-token header string true "token"
// @Param body body create_${this.funname}_req true "请求参数"
// @Success 200 "响应数据"
func create_${this.funname}(ctx *gin.Context) {
	var reqdata create_${this.funname}_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	token := admin.GetToken(ctx)
	item := new(model.)
	item.SellerID = token.SellerId
	{
    }
	tb := xapp.DbQuery().
	itb := tb.WithContext(ctx)
	err := itb.Omit(tb.CreateTime).Create(item)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}`)
			}
			if (idx == 3) {
				//update
				this.copy(`type update_${this.funname}_req struct {
}

// @Router /update_${this.funname} [post]
// @Tags a
// @Summary b
// @Param x-token header string true "token"
// @Param body body update_${this.funname}_req true "请求参数"
// @Success 200 "响应数据"
func update_${this.funname}(ctx *gin.Context) {
	var reqdata update_${this.funname}_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	token := admin.GetToken(ctx)
	tb := xapp.DbQuery().
	itb := tb.WithContext(ctx)
	itb = itb.Where(tb.SellerID.Eq(token.SellerId))
	{
	}
	item := map[string]interface{}{}
	{
	}
	_, err := itb.Updates(item)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}`)
			}
			if (idx == 4) {
				//delete
				this.copy(`type delete_${this.funname}_req struct {
}

// @Router /delete_${this.funname} [post]
// @Tags a
// @Summary b
// @Param x-token header string true "token"
// @Param body body delete_${this.funname}_req true "请求参数"
// @Success 200 "响应数据"
func delete_${this.funname}(ctx *gin.Context) {
	var reqdata delete_${this.funname}_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	token := admin.GetToken(ctx)
	tb := xapp.DbQuery().
	itb := tb.WithContext(ctx)
	itb = itb.Where(tb.SellerID.Eq(token.SellerId))
	{
	}
	_, err := itb.Delete()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}
`)
			}
			if (idx == 0) {
				this.copy(`type ${this.funname}_req struct {
	Page     int \`json:"page"\`      // 页码
	PageSize int \`json:"page_size"\` // 每页数量
}

type ${this.funname}_res struct {
	Total int64           \`json:"total"\` // 总数
	Data  []*model. \`json:"data"\`  // 数据
}

// @Router /${this.funname} [post]
// @Tags a
// @Summary b
// @Param x-token header string true "token"
// @Param body body ${this.funname}_req true "请求参数"
// @Success 200  {object} ${this.funname}_res "响应数据"
func ${this.funname}(ctx *gin.Context) {
	var reqdata ${this.funname}_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
    if reqdata.Page == 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize == 0 {
		reqdata.PageSize = 15
	}
	response := new(${this.funname}_res)
	token := admin.GetToken(ctx)
	tb := xapp.DbQuery().
	itb := tb.WithContext(ctx).Order(tb.ID.Desc())
	itb = itb.Where(tb.SellerID.Eq(token.SellerId))
	{
    }
	var err error
	response.Data, response.Total, err = itb.FindByPage((reqdata.Page-1)*reqdata.PageSize, reqdata.PageSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}`)
			}
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
