package api

import (
	api_admin "xadminapi/api/admin"
	api_hongbao "xadminapi/api/hongbao"
	api_live "xadminapi/api/live"
	api_user "xadminapi/api/user"
)

var entries ApiEntries

type ApiEntries struct {
	api_admin.ApiAdminSeller
	api_admin.ApiAdminUser
	api_admin.ApiAdminRole
	api_admin.ApiAdminLog
	api_admin.ApiConfig
	api_user.ApiUser
	api_live.ApiLiveRoom
	api_live.ApiLiveChat
	api_live.ApiLiveIpBan
	api_hongbao.ApiHongbao
}

func Entries() *ApiEntries {
	return &entries
}
