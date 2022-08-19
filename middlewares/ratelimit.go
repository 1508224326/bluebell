package middlewares

import (
	"github.com/aiyouyo/bluebell/controller"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"time"
)

// RateLimitMiddlewares 令牌桶限流中间件
func RateLimitMiddlewares(fillInterval time.Duration, cap int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		bucket := ratelimit.NewBucket(fillInterval, cap)

		if bucket.TakeAvailable(1) <= 0 {
			controller.ResponseErrorWithMsg(c, controller.CodeServeBusy, "超时等待")
			c.Abort()
		}

		c.Next()
	}
}
