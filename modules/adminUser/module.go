package adminUser

import (
	"github.com/gin-gonic/gin"
)

type Module struct{}

func (m *Module) Name() string {
	return "adminUser"
}

func (m *Module) Prefix() string {
	return "/user"
}

func (m *Module) Register(rg *gin.RouterGroup) {

	uc := new(AdminUserController)

	rg.GET("/users", uc.Index)
}
