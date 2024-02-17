package middleware

import (
	"xclientapi/server"

	"github.com/gin-gonic/gin"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokendata := server.GetToken(c)
		if tokendata == nil {
			return
		}
	}
}
