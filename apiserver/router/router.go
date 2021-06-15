package router

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jzsg/fabric-explorer/apiserver/handler"
)

// Router 全局路由
var router *gin.Engine
var onceCreateRouter sync.Once

//GetRouter 获取路由
func GetRouter() *gin.Engine {
	onceCreateRouter.Do(func() {
		router = createRouter()
	})

	return router
}

func createRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/login", handler.Login)       //登录
	r.POST("/register", handler.Register) //机构注册
	//r.Use(handler.TokenAuthMiddleware())
	jzsg := r.Group("/v1/jzsg")
	{

		jzsg.POST("/invoke/policy", handler.UploadInvokePolicy)
		jzsg.GET("/query/policy/:id", handler.QueryPolicy)
		//jzsg.POST("/invoke/policy/:id", handler.InvokePolicy)
		//jzsg.POST("/invoke/service/:id", handler.InvokeService)
		jzsg.POST("/invoke/service", handler.UploadInvokeService)
		jzsg.GET("/query/service/:id", handler.QueryService)
		jzsg.POST("/invoke/company", handler.UploadCompany)
		jzsg.GET("/query/company/:id", handler.QueryCompany)
		//jzsg.GET("/query/:id", handler.QueryPolicyByNumber) //查询链上信息
		////template := jzsg.Group("/template")
		//jzsg.StaticFS("/template", http.Dir("./template"))

		jzsg.GET("/block", handler.GetBlocks)
		jzsg.GET("/block/:id", handler.GetBlockFullInfo)
	}
	return r
}
