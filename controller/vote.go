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

// VoteForPostHandler  给指定帖子投票的接口
// @Summary 给指定帖子投票的接口 可投反对票-1 赞成票1 不投票0
// @Description 帖子投票的接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamVotePost true "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} []*models.ApiPostDetail.
// @Router /postsV2 [get]
func VoteForPostHandler(c *gin.Context) {
	vote := new(models.ParamVotePost)
	err := c.ShouldBindJSON(vote)
	if err != nil {
		errors, ok := err.(validator.ValidationErrors)
		if !ok { // json解析错误
			ResponseError(c, CodeInvalidParams)
			return
		}
		// 参数校验错误
		ResponseErrorWithMsg(c, CodeInvalidParams, ClearTopErr(errors.Translate(trans)))
		return
	}

	id, err := GetCtxUserID(c) // 获得投票的作者id
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	vote.AuthorID = id // 赋值id

	err = logic.VoteForPost(vote)
	if err != nil {
		zap.L().Error("VoteForPostHandler", zap.Error(err))
		fmt.Println(err)
		if errors.Is(err, redis.ErrorPostExpired) {
			ResponseError(c, CodeVotedExpired)
			return
		}
		if errors.Is(err, redis.ErrorRepeatVote) {
			ResponseError(c, CodeVoteRepeat)
			return
		}
		ResponseError(c, CodeServeBusy)
		return

	}

	// 返回响应
	ResponseSuccess(c, nil)
}

// VotedPostListHandler  查询投过票的帖子的接口
// @Summary 查询投过票的帖子 投过赞成票的帖子会被查询出来并分页展示
// @Description 查询投过票的帖子
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200 {object} []*models.ApiPostDetail.
// @Router /postsV2 [get]
func VotedPostListHandler(c *gin.Context) {

	// 获得发起请求的作者id
	id, err := GetCtxUserID(c)
	if err != nil { // 需要登录
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 将query string参数绑定到结构体上
	queryStr := &models.ParamPosts{ // 默认参数
		Page:        1,
		Size:        10,
		OrderBy:     models.OderByTime, // 按发布时间排序
		CommunityID: -1,                // 默认不按照社区分组
	}
	err = c.ShouldBindQuery(queryStr) // 绑定请求参数
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	postList, err := logic.GetVotedPostList(id, queryStr)

	if err != nil {
		zap.L().Error("GetVotedPostList(authorID string) failed, err:", zap.Int64("uid", id), zap.Error(err))
		if errors.Is(err, redis.ErrorEmptyIds) {
			ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
			return
		}
		ResponseError(c, CodeServeBusy)
		return
	}

	ResponseSuccess(c, postList)
}
