package service

import (
	service_app "xclientapi/service/app"
	service_hash "xclientapi/service/hash"
	service_mail "xclientapi/service/mail"
	service_marquee "xclientapi/service/marquee"
	service_slide "xclientapi/service/slide"
	service_user "xclientapi/service/user"
)

var entries *ServiceEntries = &ServiceEntries{}

type ServiceEntries struct {
	service_user.ServiceUser
	service_app.ServiceApp
	service_slide.ServiceSlide
	service_marquee.ServiceMarquee
	service_mail.ServiceMail
	service_hash.ServiceHash
}

func Init() {
	entries.ServiceUser.Init()
	entries.ServiceApp.Init()
	entries.ServiceSlide.Init()
	entries.ServiceMarquee.Init()
	entries.ServiceMail.Init()
	entries.ServiceHash.Init()
}

func Entries() *ServiceEntries {
	return entries
}
