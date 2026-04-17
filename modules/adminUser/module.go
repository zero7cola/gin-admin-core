package adminUser

import (
	"github.com/gin-gonic/gin"
)

type Module struct{}

func (m *Module) Register(rg *gin.RouterGroup) {
	g := rg.Group("/user")

	uc := new(AdminUserController)

	g.GET("/users", uc.Index)
}
