package service

import (
	service_admin "xadminapi/service/admin"
	service_user "xadminapi/service/user"
)

var entries *ServiceEntries = &ServiceEntries{}

type ServiceEntries struct {
	service_admin.ServiceAdmin
	service_user.ServiceUser
}

func Entries() *ServiceEntries {
	return entries
}

func Init() {
	entries.ServiceAdmin.Init()
	entries.ServiceUser.Init()

}
