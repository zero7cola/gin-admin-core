## 使用

### 初始化
```go
func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	db := getDB()
	redis := getRedis()
    logger := getLogger()

	// 传入配置,并传入必要的实例
	path, _ := filepath.Abs("config.yaml")
	err = core.InitWithFile(path, core.WithDB(db), core.WithRedis(redis), core.WithLogger(logger))

    if err != nil {
		panic(err)
	}

	// 初始化路由，注册外部module
	admin.Register(r, "/admin", &topic.Module{})

	err = r.Run(":8789")
        if err != nil {
		panic(err)
	}

}
```

### module示例
```go
package topic

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zero7cola/gin-admin-core/pkg/auth"
	"github.com/zero7cola/gin-admin-core/setting"
)

type Module struct{}

func (m *Module) Name() string {
	return "topic"
}

func (m *Module) Prefix() string {
	return "/topic"
}

func (m *Module) Register(rg *gin.RouterGroup) {

	data, err := json.Marshal(setting.GlobalSetting)

	if err != nil {
		panic(err)
	}

	rg.GET("/index", func(c *gin.Context) {
		currentAdmin := auth.CurrentAdminUser(c)
		c.JSON(200, gin.H{
			"data": fmt.Sprintf("%v", string(data)),
			"msg":  currentAdmin,
		})
	})
}

```