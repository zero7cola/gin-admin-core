// Package middlewares 存放系统中间件
package middlewares

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zero7cola/gin-admin-core/model/adminOperationLog"
	"github.com/zero7cola/gin-admin-core/pkg/auth"
)

func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 获取请求数据
		var requestBody []byte
		if c.Request.Body != nil {
			// c.Request.Body 是一个 buffer 对象，只能读取一次
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 读取后，重新赋值 c.Request.Body ，以供后续的其他操作
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 设置开始时间
		c.Next()

		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			adminLog := adminOperationLog.AdminOperationLog{}
			adminLog.UserId = cast.ToUint64(auth.CurrentAdminUID(c))
			adminLog.Path = c.Request.URL.Path
			adminLog.Method = c.Request.Method
			adminLog.Input = string(requestBody)

			go func() {
				adminLog.Create()
			}()
		}
	}
}
