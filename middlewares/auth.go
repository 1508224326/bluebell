package middlewares

import (
	"github.com/aiyouyo/bluebell/controller"
	"github.com/aiyouyo/bluebell/pkg/jwt"
	"github.com/gin-gonic/gin"
	"strings"
)

// JwtTokenMiddleware token验证的中间件
func JwtTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if len(token) == 0 { // 没有token
			controller.ResponseError(c, controller.CodeWithoutToken)
			c.Abort()
			return
		}

		// Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
		// eyJ1c2VyTmFtZSI6IuW8oOS4iSIsImV4cCI6MTY0ODUzMTI0OSwiaXNzIjoiYmx1ZWJlbGwuY29tIn0.
		// 15XVaBOhwKqrjhQpI1RKAVL0vqdWKHHnYNtAIBPt-RM
		// fmt.Println(token)

		// 这里对客户端上传的token进行分割
		parts := strings.SplitN(token, " ", 2)
		if !(len(parts) == 2 || parts[0] == "Bearer") { //说明token格式不正确
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}

		// 解析token
		parseToken, err := jwt.ParseToken(parts[1])
		if err != nil {
			controller.ResponseErrorWithMsg(c, controller.CodeInvalidToken, err.Error())
			c.Abort()
			return
		}

		// 中间件判断用户登录的话 将用户id保存在 context里面 保证后续可以获取到
		controller.SetCtxValue(c, controller.CtxUserID, parseToken.UserID)
		//c.Set(controller.CtxUserID, parseToken.UserID)
		c.Next()
	}
}
