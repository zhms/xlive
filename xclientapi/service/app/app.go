package service_app

import (
	"context"
	"errors"
	"fmt"
	"xclientapi/server"
	"xcom/edb"
	"xcom/enum"
	"xcom/utils"
	"xcom/xcom"

	"github.com/redis/go-redis/v9"
)

type ServiceApp struct {
}

func (this *ServiceApp) Init() {
}

type AppGetLiveInfoReq struct {
	SellerId int `json:"-"`
}

type AppGetLiveInfoRes struct {
	Name    string `json:"name"`
	Account string `json:"account"`
	LiveUrl string `json:"live_url"`
	Title   string `json:"title"`
}

func (this *ServiceApp) GetLiveInfo(host string, reqdata *AppGetLiveInfoReq) (response *AppGetLiveInfoRes, merr map[string]interface{}, err error) {
	reqdata.SellerId = xcom.GetSellerId(host)
	if reqdata.SellerId == 0 {
		return nil, enum.SellerNotFound, nil
	}
	response = &AppGetLiveInfoRes{}
	result, err := server.Redis().Client().HGet(context.Background(), "living", fmt.Sprint(reqdata.SellerId)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, nil, err
	}
	if result == "" {
		return response, nil, nil
	}
	mdata := &utils.XMap{}
	mdata.FromBytes([]byte(result))
	response.Name = mdata.String(edb.Name)
	response.Account = mdata.String(edb.Account)
	response.LiveUrl = mdata.String(edb.LiveUrl)
	response.Title = mdata.String(edb.Title)
	return response, nil, nil
}

type AppGetOnlineInfoReq struct {
	SellerId int `json:"-"`
}

type AppGetOnlineInfoRes struct {
	OnlineCount int `json:"online_count"`
}

func (this *ServiceApp) GetOnlineInfo(host string, reqdata *AppGetOnlineInfoReq) (response *AppGetOnlineInfoRes, merr map[string]interface{}, err error) {
	reqdata.SellerId = xcom.GetSellerId(host)
	if reqdata.SellerId == 0 {
		return nil, enum.SellerNotFound, nil
	}

	response = &AppGetOnlineInfoRes{}
	response.OnlineCount = 100
	return response, nil, nil
}
