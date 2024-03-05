package controller

import (
	controller_admin "xadminapi/controller/admin"
	controller_live "xadminapi/controller/live"
	controller_user "xadminapi/controller/user"
)

var entries ControllerEntries

type ControllerEntries struct {
	controller_admin.ControllerAdminSeller
	controller_admin.ControllerAdminUser
	controller_admin.ControllerAdminRole
	controller_admin.ControllerAdminLog
	controller_admin.ControllerConfig
	controller_user.ControllerUser
	controller_live.ControllerLiveRoom
	controller_live.ControllerLiveChat
	controller_live.ControllerLiveIpBan
}

func Entries() *ControllerEntries {
	return &entries
}
