package router

import (
	"fmt"
	"net/http"
	"strings"
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

	r.Use(Cors())
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

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//				允许跨域设置																										可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //	跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //	处理请求
	}
}
