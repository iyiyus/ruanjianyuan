package router

import (
	"go-source/internal/handler"
	"go-source/internal/middleware"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// 安装 API（前端调用，无需认证）
	r.GET("/api/install/status", handler.InstallStatus)
	r.POST("/api/install/check-db", handler.InstallCheckDB)
	r.POST("/api/install/run", handler.InstallRun)

	// 静态文件（必须在中间件之前，否则未安装时assets被拦截）
	r.Static("/assets", "./dist/assets")
	r.StaticFile("/favicon.ico", "./dist/favicon.ico")
	r.Static("/uploads", "./uploads")

	// 未安装拦截（只拦截API，不拦截静态文件）
	r.Use(middleware.RequireInstalled())

	// iOS软件源公开接口
	r.GET("/appstore", handler.AppStore)
	r.POST("/appstore", handler.AppStore)

	// 根路径跳转到软件源
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/appstore")
	})

	// 公开黑名单接口
	r.GET("/Blacklist", handler.GetBlacklistPublic)
	r.GET("/Blacklist/check", handler.CheckBlacklist)

	// 认证
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", handler.Login)
	}

	// 需要JWT的接口
	api := r.Group("/api")
	api.Use(middleware.JWTAuth())
	{
		// 用户信息
		api.GET("/user/info", handler.GetUserInfo)
		api.GET("/user/list", handler.GetUserList)
		api.POST("/user/change-password", handler.ChangePassword)
		api.POST("/user/profile", handler.UpdateProfile)
		api.GET("/role/list", handler.GetRoleList)

		// 菜单
		api.GET("/v3/system/menus", handler.GetMenus)

		// 仪表盘
		api.GET("/dashboard", handler.GetDashboard)

		// App管理
		api.GET("/apps", handler.GetCategoryList)
		api.GET("/apps/:id", handler.GetCategory)
		api.POST("/apps", handler.CreateCategory)
		api.PUT("/apps/:id", handler.UpdateCategory)
		api.DELETE("/apps/:id", handler.DeleteCategory)

		// 卡密管理
		api.GET("/kami", handler.GetKamiList)
		api.POST("/kami/generate", handler.GenerateKami)
		api.DELETE("/kami/:id", handler.DeleteKami)
		api.DELETE("/kami", handler.BatchDeleteKami)

		// 黑名单
		api.GET("/black", handler.GetBlackList)
		api.POST("/black", handler.CreateBlack)
		api.DELETE("/black/:id", handler.DeleteBlack)

		// 监控
		api.GET("/monitor", handler.GetMonitorList)
		api.DELETE("/monitor/:id", handler.DeleteMonitor)
		api.DELETE("/monitor", handler.ClearMonitor)

		// 系统配置
		api.GET("/config", handler.GetConfigList)
		api.GET("/config/map", handler.GetConfigMap)
		api.PUT("/config", handler.UpdateConfig)

		// 文件上传
		api.POST("/upload", handler.UploadFile)
		api.POST("/upload/ipa", handler.UploadIPA)

		// 管理员日志
		api.GET("/admin/log", handler.GetAdminLog)

		// 软件源搬运
		api.POST("/source/fetch", handler.FetchSource)
		api.POST("/source/import", handler.ImportApps)
	}

	// 托管前端静态文件（dist目录）- SPA支持
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") ||
			path == "/appstore" ||
			path == "/Blacklist" ||
			strings.HasPrefix(path, "/Blacklist/") ||
			strings.HasPrefix(path, "/uploads/") ||
			strings.HasPrefix(path, "/assets/") {
			c.JSON(404, gin.H{"code": 404, "msg": "not found"})
			return
		}
		c.File("./dist/index.html")
	})

	return r
}
