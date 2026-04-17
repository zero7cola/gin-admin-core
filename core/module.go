package core

import "github.com/gin-gonic/gin"

type Module interface {
	Name() string
	Register(rg *gin.RouterGroup)
	Prefix() string
}
