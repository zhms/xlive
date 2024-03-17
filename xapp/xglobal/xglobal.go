package xglobal

import (
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var Id string
var Project string
var Env string
var Router *gin.Engine
var Running bool = true
var Working = new(sync.WaitGroup)
var ApiV1 *gin.RouterGroup
var ApiV2 *gin.RouterGroup
var ApiV3 *gin.RouterGroup
var ApiV4 *gin.RouterGroup
var ApiV5 *gin.RouterGroup
var ApiV6 *gin.RouterGroup
var ApiV7 *gin.RouterGroup
var ApiV8 *gin.RouterGroup
var ApiV9 *gin.RouterGroup

func IsEnvPrd() bool {
	return strings.Index(Env, "prd") >= 0
}

func IsEnvDev() bool {
	return strings.Index(Env, "dev") >= 0
}
