package initialize

import (
	"asynchub/flowsvr/src/cache"
	"asynchub/flowsvr/src/db"
	"fmt"

	"github.com/gin-gonic/gin"
)

// InitResource 初始化资源, 包括数据库和缓存
func InitResource() {
	// 初始化mysql
	if err := db.InitDB(); err != nil {
		panic(fmt.Sprintf("初始化mysql数据库失败: %s", err.Error()))
	}

	// 初始化redis
	if err := cache.InitCacche(); err != nil {
		panic(fmt.Sprintf("初始化redis缓存失败: %s", err.Error()))
	}
}

func RegisterRouter(router *gin.Engine) {
	// 分组路由，v1版本
	v1 := router.Group("/v1")
	{
		// 创建任务接口，前面是路径，后面是执行的函数，跳进去
		v1.POST("/creater_task", task.CreateTask")

	}
}
