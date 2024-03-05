package service_live

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"
	"xadminapi/model"
	"xadminapi/server"
	"xcom/edb"
	"xcom/xutils"

	"github.com/gin-gonic/gin"
)

type ServiceLiveRoom struct {
}

func (this *ServiceLiveRoom) Init() {
}

type GetLiveRoomListReq struct {
	Page     int `json:"page"`      // 页码
	PageSize int `json:"page_size"` // 每页数量
}

type GetLiveRoomListRes struct {
	Total int64             `json:"total"` // 总数
	Data  []model.XLiveRoom `json:"data"`  // 数据
}

func (this *ServiceLiveRoom) GetLiveRoomList(ctx *gin.Context, idata interface{}) (rdata interface{}, merr map[string]interface{}, err error) {
	reqdata := idata.(GetLiveRoomListReq)
	if reqdata.Page <= 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize <= 0 {
		reqdata.PageSize = 15
	}
	token := server.GetToken(ctx)
	data := GetLiveRoomListRes{}
	db := server.Db().Model(&model.XLiveRoom{})
	db = xutils.DbWhere(db, edb.SellerId, token.SellerId, int(0))
	db = db.Count(&data.Total)
	db = db.Offset((reqdata.Page - 1) * reqdata.PageSize)
	db = db.Limit(reqdata.PageSize)
	db = db.Order(edb.Id + edb.DESC)
	err = db.Find(&data.Data).Error
	if err != nil {
		return nil, nil, err
	}
	return data, nil, nil
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

func (this *ServiceLiveRoom) get_stream_url(app string, name string, streamurl string) (string, string) {
	pushurl := push_url("push."+streamurl, "gb905x9LUXR6BPPu", app, name, 1440*60)
	pullurl := pull_url("pull."+streamurl, "aPgxS12GHIamE8lV", app, name, 1440*60)
	return pushurl, pullurl
}

type CreateLiveRoomReq struct {
	Name    string `json:"name"`    // 直播间名称
	Account string `json:"account"` // 直播间账号
	State   int    `json:"state"`   // 直播间状态
	Title   string `json:"title"`   // 直播间标题
}

func (this *ServiceLiveRoom) CreateLiveRoom(ctx *gin.Context, idata interface{}) (merr map[string]interface{}, err error) {
	reqdata := idata.(CreateLiveRoomReq)
	token := server.GetToken(ctx)
	if token == nil {
		return nil, nil
	}
	createdata := model.XLiveRoom{
		SellerId: token.SellerId,
		Name:     reqdata.Name,
		Account:  reqdata.Account,
		State:    reqdata.State,
		Title:    reqdata.Title,
	}
	db := server.Db().Model(&model.XLiveRoom{}).Omit(edb.CreateTime).Create(&createdata)
	if db.Error != nil {
		return nil, db.Error
	}
	return nil, db.Error
}

type UpdateLiveRoomReq struct {
	Id      int    `json:"id"`      // 直播间Id
	Name    string `json:"name"`    // 直播间名称
	Account string `json:"account"` // 直播间账号
	State   int    `json:"state"`   // 直播间状态
	Title   string `json:"title"`   // 直播间标题
}

func (this *ServiceLiveRoom) UpdateLiveRoom(ctx *gin.Context, idata interface{}) (merr map[string]interface{}, err error) {
	reqdata := idata.(UpdateLiveRoomReq)
	token := server.GetToken(ctx)
	if token == nil {
		return nil, nil
	}
	rows, _ := server.Db().Table(edb.TableLiveRoom).Where(edb.Id+edb.EQ, reqdata.Id).Rows()
	roomdata := xutils.DbFirst(rows)
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
	if reqdata.State == edb.StateYes || reqdata.State == edb.StateNo && roomdata.Int(edb.State) != reqdata.State {
		updatedata["state"] = reqdata.State
		if reqdata.State == 1 {
			streamurl := ""
			clienturl := ""
			appname := ""
			server.Db().Table(edb.TableKv).Select(edb.V).Where(edb.K+edb.EQ, "client_url").Row().Scan(&clienturl)
			server.Db().Table(edb.TableKv).Select(edb.V).Where(edb.K+edb.EQ, "app_name").Row().Scan(&appname)
			server.Db().Table(edb.TableKv).Select(edb.V).Where(edb.K+edb.EQ, "stream_url").Row().Scan(&streamurl)
			pushurl, pullurl := this.get_stream_url(appname, "r"+fmt.Sprint(reqdata.Id), streamurl)
			updatedata["push_url"] = pushurl
			updatedata["pull_url"] = pullurl
			updatedata["live_url"] = clienturl + "?r=" + fmt.Sprint(reqdata.Id)
		}
	}
	db := server.Db().Model(&model.XLiveRoom{}).Where("id = ? and seller_id = ?", reqdata.Id, token.SellerId).Updates(updatedata)
	if db.Error != nil {
		return nil, db.Error
	}
	if reqdata.State == 1 {
		rows, err := server.Db().Table(edb.TableLiveRoom).Where(edb.Id+edb.EQ, reqdata.Id).Rows()
		if err != nil {
			return nil, err
		}
		roomdata := xutils.DbFirst(rows)
		_, err = server.Redis().Client().HSet(context.Background(), "living", fmt.Sprintf("%v_%v", token.SellerId, reqdata.Id), roomdata.ToString()).Result()
		if err != nil {
			return nil, err
		}
	} else {
		_, err := server.Redis().Client().HDel(context.Background(), "living", fmt.Sprintf("%v_%v", token.SellerId, reqdata.Id)).Result()
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

type DeleteLiveRoomReq struct {
	Id int `json:"id"` // 直播间Id
}

func (this *ServiceLiveRoom) DeleteLiveRoom(ctx *gin.Context, idata interface{}) (merr map[string]interface{}, err error) {
	reqdata := idata.(DeleteLiveRoomReq)
	token := server.GetToken(ctx)
	if token == nil {
		return nil, nil
	}
	db := server.Db().Model(&model.XLiveRoom{}).Where(edb.SellerId+edb.EQ, token.SellerId).Where(edb.Id+edb.EQ, reqdata.Id).Delete(&model.XLiveRoom{})
	if db.Error != nil {
		return nil, db.Error
	}
	server.Redis().Client().HDel(context.Background(), "living", fmt.Sprintf("%v_%v", token.SellerId, reqdata.Id)).Result()
	return nil, db.Error
}
