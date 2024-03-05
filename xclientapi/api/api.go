package api

import (
	api_app "xclientapi/api/app"
	api_user "xclientapi/api/user"
)

var entries ApiEntries

type ApiEntries struct {
	api_app.ApiApp
	api_user.ApiUser
}

func Entries() *ApiEntries {
	return &entries
}
