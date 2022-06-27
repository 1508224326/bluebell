package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	CtxUserID = "user_id" // 设置上下文存储用id的Key
)

var (
	ErrorUserNotLogin = errors.New("用户未登录")
)

// GetCtxUserID 获取用户id
func GetCtxUserID(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(CtxUserID) // 获取

	if !ok { // 不存在
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64) // 获取的不对
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

func SetCtxValue(c *gin.Context, key string, value interface{}) {
	func(c *gin.Context) {
		c.Set(key, value)
	}(c)
}

func getPageSize(c *gin.Context) (int64, int64) {
	var (
		page int64 = 1
		size int64 = 10
		err  error = nil
	)
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}

	if size <= 0 || page <= 0 { // 保证大于0
		page = 1
		size = 0
	}

	return page, size
}
