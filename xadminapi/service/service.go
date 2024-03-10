package service

import (
	service_admin "xadminapi/service/admin"
	service_hongbao "xadminapi/service/hongbao"
	service_live "xadminapi/service/live"
	service_user "xadminapi/service/user"
)

var entries *ServiceEntries = &ServiceEntries{}

type ServiceEntries struct {
	service_admin.ServiceAdmin
	service_user.ServiceUser
	service_live.ServiceLiveRoom
	service_live.ServiceLiveChat
	service_live.ServiceLiveIpBan
	service_hongbao.ServiceHongbao
}

func Entries() *ServiceEntries {
	return entries
}

func Init() {
	entries.ServiceAdmin.Init()
	entries.ServiceUser.Init()
	entries.ServiceLiveRoom.Init()
	entries.ServiceLiveChat.Init()
	entries.ServiceLiveIpBan.Init()
	entries.ServiceHongbao.Init()
}
