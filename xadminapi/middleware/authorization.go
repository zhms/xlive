package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"xadminapi/model"
	"xadminapi/server"

	"xcom/enum"
	"xcom/global"
	"xcom/utils"

	"github.com/beego/beego/logs"
	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Authorization(mainmenu string, submenu string, opt string, optname string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if global.IsEnvPrd() && optname != "" {
			VerifyCode := c.Request.Header.Get("VerifyCode")
			if VerifyCode == "" {
				c.JSON(http.StatusBadRequest, enum.VerifyNotFoundCode)
				c.Abort()
				return
			}
			if len(VerifyCode) != 6 {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": enum.VerifyInCorrectCode,
					"msg":  "谷歌验证码不正确",
				})
				c.Abort()
				return
			}
			tokendata := server.GetToken(c)
			if tokendata == nil {
				return
			}
			if tokendata.GoogleSecret == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": enum.VerifyNotFoundSecret,
					"msg":  "未绑定谷歌验证码",
				})
				c.Abort()
				return
			}
			if !utils.VerifyGoogleCode(tokendata.GoogleSecret, VerifyCode) {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": enum.VerifyInCorrectCode,
					"msg":  "谷歌验证码不正确",
				})
				c.Abort()
				return
			}
		}
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		tokendata := server.GetToken(c)
		if tokendata == nil {
			return
		}
		mapdata := map[string]interface{}{}
		json.Unmarshal([]byte(tokendata.AuthData), &mapdata)
		m, ok := mapdata[mainmenu]
		if !ok {
			c.JSON(http.StatusUnauthorized, enum.AuthNotFoundMainMenu)
			c.Abort()
			return
		}
		mm, ok := m.(map[string]interface{})
		if !ok {
			c.JSON(http.StatusUnauthorized, enum.AuthNotFoundMainMenu)
			c.Abort()
			return
		}
		s, ok := mm[submenu]
		if !ok {
			c.JSON(http.StatusUnauthorized, enum.AuthNotFoundSubMenu)
			c.Abort()
			return
		}
		ms, ok := s.(map[string]interface{})
		if !ok {
			c.JSON(http.StatusUnauthorized, enum.AuthNotFoundSubMenu)
			c.Abort()
			return
		}
		o, ok := ms[opt]
		if !ok {
			c.JSON(http.StatusUnauthorized, enum.AuthNotFoundOpt)
			c.Abort()
			return
		}
		fo, ok := o.(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, enum.AuthNotFoundOpt)
			c.Abort()
			return
		}
		if fo != enum.StateYes {
			c.JSON(http.StatusUnauthorized, enum.AuthNotAllow)
			c.Abort()
			return
		}
		//获取请求数据
		reqbytes, _ := io.ReadAll(c.Request.Body)
		req := string(reqbytes)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqbytes))

		c.Next()
		if optname == "" {
			return
		}

		//响应数据
		//res := bodyLogWriter.body.String()
		// //构造日志对象
		// logdata := map[string]any{}
		// logdata["path"] = c.Request.URL.Path
		// logdata["method"] = c.Request.Method
		// logdata["ip"] = c.ClientIP()
		// logdata["time"] = startTime.Format("2006-01-02 15:04:05")

		// logdata["res"] = res
		// logdata["status"] = c.Writer.Status()
		// logdata["req"] = req

		// logdata["token"] = gin.H{
		// 	"SellerId":  tokendata.SellerId,
		// 	"Account":   tokendata.Account,
		// 	"UserId":    tokendata.UserId,
		// }
		// //写日志
		// bytes, _ := json.Marshal(logdata)
		// server.Redis(enum.Redis_Token).LPush("opt:log", string(bytes))

		optlog := model.XAdminOptLog{}
		optlog.SellerId = tokendata.SellerId
		optlog.Account = tokendata.Account
		optlog.ReqData = req
		optlog.ReqPath = c.Request.URL.Path
		optlog.OptName = optname
		optlog.ReqIp = c.ClientIP()
		optlog.CreateTime = utils.Now()
		err := server.Db().Model(&optlog).Create(&optlog).Error
		if err != nil {
			logs.Error("写入操作日志失败", err.Error())
		}
	}
}
