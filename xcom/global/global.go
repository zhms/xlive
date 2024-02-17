package global

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
var WorkGroup = sync.WaitGroup{}

// 判断是否是生产环境
func IsEnvPrd() bool {
	return strings.Index(Env, "prd") >= 0
}
