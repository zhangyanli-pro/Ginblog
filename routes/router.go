package routes

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/wejectchen/ginblog/api/v1"
	"github.com/wejectchen/ginblog/middleware"
	"github.com/wejectchen/ginblog/utils"
)

func createMyRender() multitemplate.Renderer {
	p := multitemplate.NewRenderer()
	p.AddFromFiles("admin", "web/admin/dist/index.html")
	p.AddFromFiles("front", "web/front/dist/index.html")
	return p
}

func InitRouter() {
	// zyl备注：设置gin框架的运行模式debug/release/test
	gin.SetMode(utils.AppMode)
	// zyl备注：区分gin.New()和gin.Default()的区别:
	// gin.New()方法创建一个空的 Gin Engine 实例，该实例不包含任何中间件（middleware）和路由（router）
	//gin.Default()方法创建一个带有默认中间件和默认路由的 Gin Engine 实例。 默认中间件包括日志记录（Logger）、恢复（Recovery）和静态文件服务（StaticFile）等;默认路由包含了常见的HTTP方法（GET、POST、PUT、DELETE等）的处理方法。
	r := gin.New()
	// 设置信任网络 []string
	// nil 为不计算，避免性能消耗，上线应当设置
	_ = r.SetTrustedProxies(nil)

	r.HTMLRender = createMyRender()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	r.Static("/static", "./web/front/dist/static")
	r.Static("/admin", "./web/admin/dist")
	r.StaticFile("/favicon.ico", "/web/front/dist/favicon.ico")

	// zyl备注：为什么不走此处的路由？？？
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "front", nil)
	})

	r.GET("/admin", func(c *gin.Context) {
		c.HTML(200, "admin", nil)
	})

	/*
		后台管理路由接口
	*/
	auth := r.Group("api/v1")
	// zyl备注：自定义中间件，这里表示该中间件仅作用于/api/v1开头的路由
	auth.Use(middleware.JwtToken())
	{
		// 用户模块的路由接口
		auth.GET("admin/users", v1.GetUsers)
		auth.PUT("user/:id", v1.EditUser)
		auth.DELETE("user/:id", v1.DeleteUser)
		//修改密码
		auth.PUT("admin/changepw/:id", v1.ChangeUserPassword)
		// 分类模块的路由接口
		auth.GET("admin/category", v1.GetCate)
		auth.POST("category/add", v1.AddCategory)
		auth.PUT("category/:id", v1.EditCate)
		auth.DELETE("category/:id", v1.DeleteCate)
		// 文章模块的路由接口
		auth.GET("admin/article/info/:id", v1.GetArtInfo)
		auth.GET("admin/article", v1.GetArt)
		auth.POST("article/add", v1.AddArticle)
		auth.PUT("article/:id", v1.EditArt)
		auth.DELETE("article/:id", v1.DeleteArt)
		// 上传文件
		auth.POST("upload", v1.UpLoad)
		// 更新个人设置
		auth.GET("admin/profile/:id", v1.GetProfile)
		auth.PUT("profile/:id", v1.UpdateProfile)
		// 评论模块
		auth.GET("comment/list", v1.GetCommentList)
		auth.DELETE("delcomment/:id", v1.DeleteComment)
		auth.PUT("checkcomment/:id", v1.CheckComment)
		auth.PUT("uncheckcomment/:id", v1.UncheckComment)
	}

	/*
		前端展示页面接口
	*/
	router := r.Group("api/v1")
	{
		// 用户信息模块
		router.POST("user/add", v1.AddUser)
		router.GET("user/:id", v1.GetUserInfo)
		router.GET("users", v1.GetUsers)

		// 文章分类信息模块
		router.GET("category", v1.GetCate)
		router.GET("category/:id", v1.GetCateInfo)

		// 文章模块
		router.GET("article", v1.GetArt)
		router.GET("article/list/:id", v1.GetCateArt)
		router.GET("article/info/:id", v1.GetArtInfo)

		// 登录控制模块
		router.POST("login", v1.Login)
		router.POST("loginfront", v1.LoginFront)

		// 获取个人设置信息
		router.GET("profile/:id", v1.GetProfile)

		// 评论模块
		router.POST("addcomment", v1.AddComment)
		router.GET("comment/info/:id", v1.GetComment)
		router.GET("commentfront/:id", v1.GetCommentListFront)
		router.GET("commentcount/:id", v1.GetCommentCount)
	}

	_ = r.Run(utils.HttpPort)

}
