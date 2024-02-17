package controller

import (
	controller_app "xclientapi/controller/app"
	controller_hash "xclientapi/controller/hash"
	controller_mail "xclientapi/controller/mail"
	controller_marquee "xclientapi/controller/marquee"
	controller_slide "xclientapi/controller/slide"
	controller_user "xclientapi/controller/user"
)

var entries ControllerEntries

type ControllerEntries struct {
	controller_app.ControllerApp
	controller_user.ControllerUser
	controller_mail.ControllerMail
	controller_marquee.ControllerMarquee
	controller_slide.ControllerSlide
	controller_hash.ControllerHash
}

func Entries() *ControllerEntries {
	return &entries
}
