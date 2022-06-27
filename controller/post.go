package controller

import (
	"errors"
	"fmt"
	"github.com/aiyouyo/bluebell/dao/redis"
	"github.com/aiyouyo/bluebell/logic"
	"github.com/aiyouyo/bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// CreatePostHandler  创建帖子的的接口
// @Summary 创建帖子的的接口
// @Description 输入标题 内容 选择社区创建帖子
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object json models.PostBase true "帖子内容参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /post [post]
func CreatePostHandler(c *gin.Context) {

	// 1. 获取参数
	//c.ShouldBindJSON()  // 将前端请求的json参数绑定到模板上
	post := new(models.PostBase)
	err := c.ShouldBindJSON(post)
	if err != nil { //
		errors, ok := err.(validator.ValidationErrors)
		if !ok { // 这一步说明请求参数序列化出错了, 提交的参数有误
			ResponseError(c, CodeInvalidParams)
			return
		}
		// 说明请求参数不完整 validator 校验出错
		ResponseErrorWithMsg(c, CodeInvalidParams, ClearTopErr(errors.Translate(trans)))
		return
	}

	// 获取用户id
	uid, err := GetCtxUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
	}
	post.AuthorID = uid

	// 2. 调用logic 存储帖子
	err = logic.CreatePost(post)
	if err != nil {
		zap.L().Error("CreatePost err : ", zap.Error(err))
		fmt.Println(err, "出错啦")
		ResponseError(c, CodeServeBusy) // 后台出错 不需要暴露错误给用户
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// PostDetailHandler  查询帖子详细的的接口
// @Summary 查询帖子详细的的接口
// @Description 根据路径参数查询帖子的具体内容
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param int true "帖子内容参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /post/:id [get]
func PostDetailHandler(c *gin.Context) {

	// 获取post id的参数
	pid := c.Param("id")

	// //2. 查数据库
	data, err := logic.GetPostDetailByID(pid)
	if err != nil {
		fmt.Println(err)
		zap.L().Error("logic.GetPostDetailByID failed: ", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}

	// 3.返回响应
	ResponseSuccess(c, data)

}

// PostListHandler  查询帖子分页展示的接口
// @Summary 查询帖子详细的的接口 在数据库中查询按照时间顺序
// @Description 在数据库中查询按照时间顺序
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200 {object} []*models.ApiPostDetail.
// @Router /posts [get]
func PostListHandler(c *gin.Context) {

	// 获取请求参数
	page, size := getPageSize(c) // 页面 和 size

	// 查询所有帖子
	postList, err := logic.GetPostList(page, size)
	if err != nil {
		ResponseError(c, CodeServeBusy)
		return
	}

	// 返回响应
	ResponseSuccess(c, postList)

}

// PostListHandlerV2  查询帖子分页展示的接口升级版
// @Summary 查询帖子详细的的接口 可查询投票信息 按照分数或者时间排序 同时可以按照社区分组
// @Description 查询帖子详细的的接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param {object} *models.ParamPosts json true "1 10 time 1"
// @Security ApiKeyAuth
// @Success 200 {object} []*models.ApiPostDetail.
// @Router /postsV2 [get]
func PostListHandlerV2(c *gin.Context) {

	// 将query string参数绑定到结构体上
	queryStr := &models.ParamPosts{ // 默认参数
		Page:        1,
		Size:        10,
		OrderBy:     models.OderByTime, // 按发布时间排序
		CommunityID: -1,                // 默认不按照社区分组
	}
	err := c.ShouldBindQuery(queryStr)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	fmt.Println(queryStr)
	postList, err := logic.GetPostListV2(queryStr)

	if err != nil {
		if errors.Is(err, redis.ErrorEmptyIds) {
			ResponseError(c, CodeInvalidParams)
			return
		}
		ResponseError(c, CodeServeBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, postList)

}
