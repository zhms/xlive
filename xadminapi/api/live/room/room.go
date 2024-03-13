package live_room

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"xadminapi/api/admin"
	"xapp/xapp"
	"xapp/xdb"
	"xapp/xenum"
	"xapp/xglobal"

	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
)

func Init() {
	xglobal.ApiV1.POST("/get_live_room", admin.Auth("系统管理", "角色管理", "查", ""), get_live_room)
	xglobal.ApiV1.POST("/create_live_room", admin.Auth("系统管理", "角色管理", "增", "创建直播间"), create_live_room)
	xglobal.ApiV1.POST("/update_live_room", admin.Auth("系统管理", "角色管理", "改", "更新直播间"), update_live_room)
	xglobal.ApiV1.POST("/delete_live_room", admin.Auth("系统管理", "角色管理", "删", "删除直播间"), delete_live_room)
}

func push_url(pushDomain, pushKey, appName, streamName string, expireTime int) string {
	pushURL := ""
	if pushKey == "" {
		pushURL = fmt.Sprintf("rtmp://%s/%s/%s", pushDomain, appName, streamName)
		return pushURL
	}

	timeStamp := time.Now().Unix() + int64(expireTime)
	sstring := fmt.Sprintf("/%s/%s-%d-0-0-%s", appName, streamName, timeStamp, pushKey)
	md5hash := fmt.Sprintf("%x", md5.Sum([]byte(sstring)))
	pushURL = fmt.Sprintf("rtmp://%s/%s/%s?auth_key=%d-0-0-%s", pushDomain, appName, streamName, timeStamp, md5hash)
	return pushURL
}

func pull_url(playDomain, playKey, appName, streamName string, expireTime int) (flvPlayURL string) {
	if playKey == "" {
		flvPlayURL = fmt.Sprintf("https://%s/%s/%s.flv", playDomain, appName, streamName)
	} else {
		timeStamp := time.Now().Unix() + int64(expireTime)
		flvSString := fmt.Sprintf("/%s/%s.flv-%d-0-0-%s", appName, streamName, timeStamp, playKey)
		flvMD5Hash := fmt.Sprintf("%x", md5.Sum([]byte(flvSString)))
		flvPlayURL = fmt.Sprintf("https://%s/%s/%s.flv?auth_key=%d-0-0-%s", playDomain, appName, streamName, timeStamp, flvMD5Hash)
	}
	return flvPlayURL
}

func get_stream_url(app string, name string, streamurl string) (string, string) {
	pushurl := push_url("push."+streamurl, "gb905x9LUXR6BPPu", app, name, 1440*60)
	pullurl := pull_url("pull."+streamurl, "aPgxS12GHIamE8lV", app, name, 1440*60)
	return pushurl, pullurl
}

type get_live_room_req struct {
	Page     int `json:"page"`      // 页码
	PageSize int `json:"page_size"` // 每页数量
}

type get_live_room_res struct {
	Total int64           `json:"total"` // 总数
	Data  []xdb.XLiveRoom `json:"data"`  // 数据
}

// @Router /get_live_room [post]
// @Tags 直播间 - 直播间
// @Summary 获取直播间
// @Param x-token header string true "token"
// @Param body body get_live_room_req true "请求参数"
// @Success 200 {object} get_live_room_res "响应数据"
func get_live_room(ctx *gin.Context) {
	var reqdata get_live_room_req
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
	if token == nil {
		return
	}
	response := new(get_live_room_res)
	db := xapp.Db().Model(&xdb.XLiveRoom{})
	db = db.Where(xdb.SellerId+xdb.EQ, token.SellerId)
	db = db.Count(&response.Total)
	db = db.Offset((reqdata.Page - 1) * reqdata.PageSize)
	db = db.Limit(reqdata.PageSize)
	db = db.Order(xdb.Id + xdb.DESC)
	err := db.Find(&response.Data).Error
	if err != nil {
		ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.MakeSucess(response))
}

type create_live_room_req struct {
	Name    string `json:"name"`    // 直播间名称
	Account string `json:"account"` // 直播间账号
	State   int    `json:"state"`   // 直播间状态
	Title   string `json:"title"`   // 直播间标题
}

// @Router /create_live_room [post]
// @Tags 直播间 - 直播间
// @Summary 创建直播间
// @Param x-token header string true "token"
// @Param body body create_live_room_req true "请求参数"
// @Success 200 "响应数据"
func create_live_room(ctx *gin.Context) {
	var reqdata create_live_room_req
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
	if token == nil {
		return
	}
	createdata := xdb.XLiveRoom{
		SellerId: token.SellerId,
		Name:     reqdata.Name,
		Account:  reqdata.Account,
		State:    reqdata.State,
		Title:    reqdata.Title,
	}
	err := xapp.Db().Model(&xdb.XLiveRoom{}).Omit(xdb.CreateTime).Create(&createdata).Error
	if err != nil {
		ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}

type update_live_room_req struct {
	Id      int    `json:"id"`      // 直播间Id
	Name    string `json:"name"`    // 直播间名称
	Account string `json:"account"` // 直播间账号
	State   int    `json:"state"`   // 直播间状态
	Title   string `json:"title"`   // 直播间标题
}

// @Router /update_live_room [post]
// @Tags 直播间 - 直播间
// @Summary 更新直播间
// @Param x-token header string true "token"
// @Param body body update_live_room_req true "请求参数"
// @Success 200 "响应数据"
func update_live_room(ctx *gin.Context) {
	var reqdata update_live_room_req
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
	if token == nil {
		return
	}
	roomdata := &xdb.XLiveRoom{}
	db := xapp.Db().Model(roomdata).Where(xdb.SellerId+xdb.EQ, token.SellerId).Where(xdb.Id+xdb.EQ, reqdata.Id).First(roomdata)
	if db.Error != nil {
		ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, db.Error.Error()))
		return
	}
	updatedata := map[string]interface{}{}
	if reqdata.Name != "" {
		updatedata["name"] = reqdata.Name
	}
	if reqdata.Account != "" {
		updatedata["account"] = reqdata.Account
	}
	if reqdata.Title != "" {
		updatedata["title"] = reqdata.Title
	}
	if reqdata.State == xdb.StateYes || reqdata.State == xdb.StateNo && roomdata.State != reqdata.State {
		updatedata["state"] = reqdata.State
		if reqdata.State == 1 {
			streamurl := ""
			clienturl := ""
			appname := ""
			xapp.Db().Table(xdb.TableXKv).Select(xdb.V).Where(xdb.K+xdb.EQ, "client_url").Row().Scan(&clienturl)
			xapp.Db().Table(xdb.TableXKv).Select(xdb.V).Where(xdb.K+xdb.EQ, "app_name").Row().Scan(&appname)
			xapp.Db().Table(xdb.TableXKv).Select(xdb.V).Where(xdb.K+xdb.EQ, "stream_url").Row().Scan(&streamurl)
			pushurl, pullurl := get_stream_url(appname, "r"+fmt.Sprint(reqdata.Id), streamurl)
			updatedata["push_url"] = pushurl
			updatedata["pull_url"] = pullurl
			updatedata["live_url"] = clienturl + "?r=" + fmt.Sprint(reqdata.Id)
		}
	}
	db = xapp.Db().Model(roomdata).Where(xdb.SellerId+xdb.EQ, token.SellerId).Where(xdb.Id+xdb.EQ, reqdata.Id).Updates(updatedata)
	if db.Error != nil {
		ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, db.Error.Error()))
		return
	}
	if reqdata.State == 1 {
		db = xapp.Db().Model(roomdata).Where(xdb.Id+xdb.EQ, reqdata.Id).First(roomdata)
		if db.Error != nil {
			ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, db.Error.Error()))
			return
		}
		bytes, _ := json.Marshal(roomdata)
		_, err := xapp.Redis().Client().HSet(context.Background(), "living", fmt.Sprintf("%v_%v", token.SellerId, reqdata.Id), bytes).Result()
		if err != nil {
			ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, db.Error.Error()))
			return
		}
	} else {
		_, err := xapp.Redis().Client().HDel(context.Background(), "living", fmt.Sprintf("%v_%v", token.SellerId, reqdata.Id)).Result()
		if err != nil {
			ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, db.Error.Error()))
			return
		}
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}

type delete_live_room_req struct {
	Id int `json:"id" validate:"required"` // 直播间Id
}

// @Router /delete_live_room [post]
// @Tags 直播间 - 直播间
// @Summary 删除直播间
// @Param x-token header string true "token"
// @Param body body delete_live_room_req true "请求参数"
// @Success 200 "响应数据"
func delete_live_room(ctx *gin.Context) {
	var reqdata delete_live_room_req
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
	if token == nil {
		return
	}
	roomdata := &xdb.XLiveRoom{}
	db := xapp.Db().Model(roomdata).Where(xdb.SellerId+xdb.EQ, token.SellerId).Where(xdb.Id+xdb.EQ, reqdata.Id).Delete(roomdata)
	if db.Error != nil {
		ctx.JSON(http.StatusOK, xenum.MakeError(xenum.InternalError, db.Error.Error()))
		return
	}
	ctx.JSON(http.StatusOK, xenum.Success)
}
