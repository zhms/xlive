package service

import (
	service_app "xclientapi/service/app"
	service_user "xclientapi/service/user"
)

var entries *ServiceEntries = &ServiceEntries{}

type ServiceEntries struct {
	service_user.ServiceUser
	service_app.ServiceApp
}

func Init() {
	entries.ServiceApp.Init()
	entries.ServiceUser.Init()
}

func Entries() *ServiceEntries {
	return entries
}
