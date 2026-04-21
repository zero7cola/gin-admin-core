package core

import "github.com/gin-gonic/gin"

type Module interface {
	Name() string
	Register(c *gin.Context, rg *gin.RouterGroup)
	Prefix() string
}
