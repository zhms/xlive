package controller

import (
	controller_app "xclientapi/controller/app"
	controller_user "xclientapi/controller/user"
)

var entries ControllerEntries

type ControllerEntries struct {
	controller_app.ControllerApp
	controller_user.ControllerUser
}

func Entries() *ControllerEntries {
	return &entries
}
