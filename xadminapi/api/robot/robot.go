package robot

import (
	"net/http"
	"xadminapi/api/admin"
	"xapp/xapp"
	"xapp/xdb/model"
	"xapp/xenum"
	"xapp/xglobal"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

func Init() {
	xglobal.ApiV1.POST("/get_robot", admin.Auth("机器人管理", "机器人管理", "查", ""), get_robot)
	xglobal.ApiV1.POST("/create_robot", admin.Auth("机器人管理", "机器人管理", "增", "创建机器人"), create_robot)
	xglobal.ApiV1.POST("/update_robot", admin.Auth("机器人管理", "机器人管理", "改", "更新机器人"), update_robot)
	xglobal.ApiV1.POST("/delete_robot", admin.Auth("机器人管理", "机器人管理", "删", "删除机器人"), delete_robot)
}

type get_robot_req struct {
	Page     int    `json:"page"`      // 页码
	PageSize int    `json:"page_size"` // 每页数量
	Account  string `json:"account"`   // 账号
}

type get_robot_res struct {
	Total int64           `json:"total"` // 总数
	Data  []*model.XRobot `json:"data"`  // 数据
}

// @Router /get_robot [post]
// @Tags a
// @Summary b
// @Param x-token header string true "token"
// @Param body body get_robot_req true "请求参数"
// @Success 200  {object} get_robot_res "响应数据"
func get_robot(ctx *gin.Context) {
	var reqdata get_robot_req
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
	response := new(get_robot_res)
	token := admin.GetToken(ctx)
	tb := xapp.DbQuery().XRobot
	itb := tb.WithContext(ctx).Order(tb.ID.Desc())
	itb = itb.Where(tb.SellerID.Eq(token.SellerId))
	if reqdata.Account != "" {
		itb = itb.Where(tb.Account.Eq(reqdata.Account))
	}
	var err error
	response.Data, response.Total, err = itb.FindByPage((reqdata.Page-1)*reqdata.PageSize, reqdata.PageSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
	ctx.JSON(http.StatusOK, xenum.Success)
}

type create_robot_req struct {
	Account string `json:"account" validate:"required"` // 账号
}

// @Router /create_robot [post]
// @Tags a
// @Summary b
// @Param x-token header string true "token"
// @Param body body create_robot_req true "请求参数"
// @Success 200 "响应数据"
func create_robot(ctx *gin.Context) {
	var reqdata create_robot_req
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
	robot := new(model.XRobot)
	robot.SellerID = token.SellerId
	robot.Account = reqdata.Account
	tb := xapp.DbQuery().XRobot
	itb := tb.WithContext(ctx)
	err := itb.Omit(tb.CreateTime).Create(robot)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}

type update_robot_req struct {
}

// @Router /update_robot [post]
// @Tags a
// @Summary b
// @Param x-token header string true "token"
// @Param body body update_robot_req true "请求参数"
// @Success 200 "响应数据"
func update_robot(ctx *gin.Context) {
	var reqdata update_robot_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	//response := new(update_robot_res)
	//ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
	//ctx.JSON(http.StatusOK, xenum.Success)
}

type delete_robot_req struct {
}

// @Router /delete_robot [post]
// @Tags a
// @Summary b
// @Param x-token header string true "token"
// @Param body body delete_robot_req true "请求参数"
// @Success 200 "响应数据"
func delete_robot(ctx *gin.Context) {
	var reqdata delete_robot_req
	if err := ctx.ShouldBindJSON(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	validator := val.New()
	if err := validator.Struct(&reqdata); err != nil {
		ctx.JSON(http.StatusBadRequest, xenum.MakeError(xenum.BadParams, err.Error()))
		return
	}
	//response := new(delete_robot_res)
	//ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
	//ctx.JSON(http.StatusOK, xenum.Success)
}
