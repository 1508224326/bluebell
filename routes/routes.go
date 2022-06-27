package routes

import (
	"github.com/aiyouyo/bluebell/controller"
	"github.com/aiyouyo/bluebell/logger"
	"github.com/aiyouyo/bluebell/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {

	r := gin.New()

	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.LoadHTMLFiles("templates/index.html")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	//r.GET("/", func(c *gin.Context) {
	//	//c.HTML(http.StatusOK, "index.html", nil)
	//	c.String(http.StatusOK, "bluebell ", settingts.Conf.Version)
	//})

	// 注册业务路由
	r.POST("api/v1/signup", controller.SignHandler)

	// 登录业务路由
	r.POST("api/v1/login", controller.LoginHandler)

	v1 := r.Group("/api/v1")

	v1.GET("/posts2", controller.PostListHandlerV2) // 查询帖子 可传递query string参数 按照时间和分数查询
	v1.GET("/topc", controller.TopCommunityHandler)

	v1.Use(middlewares.JwtTokenMiddleware()) // 给后面的路由使用token验证的中间件 看用户是否登录

	{
		v1.GET("/community", controller.CommunityListHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.PostDetailHandler)

		v1.GET("/posts", controller.PostListHandler) // 所有的帖子 从数据库中查询
		v1.POST("/vote", controller.VoteForPostHandler)

		v1.GET("/stars", controller.VotedPostListHandler) // 查询我赞过的

	}

	return r

}
