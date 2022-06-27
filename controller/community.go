package controller

import (
	"errors"
	"github.com/aiyouyo/bluebell/dao/mysql"
	"github.com/aiyouyo/bluebell/logic"
	"github.com/aiyouyo/bluebell/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CommunityListHandler 获取社区列表信息的接口
// @Summary 获取所有列表的接口
// @Description
// @Tags 社区相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200 {object} *models.CommunityBase
// @Router /community [get]
func CommunityListHandler(c *gin.Context) {
	// 调用逻辑层的处理
	communityList, err := logic.GetCommunityList()
	if err != nil {
		ResponseError(c, CodeServeBusy)
		return
	}

	ResponseSuccess(c, communityList)
}

// CommunityDetailHandler 获取社区的详情
// @Summary 获取社区的详情 携带有id路径参数
// @Description 根据社区ID获得社区的详细信息 创建时间描述等
// @Tags 社区相关接口
// @Accept application/json
// @Produce application/json
// @Param int query true "社区ID查询参数"
// @Param Authorization header string false "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200 {object} *models.CommunityDetail
// @Router /community/:id [get]
func CommunityDetailHandler(c *gin.Context) {

	// 1. 获取路径参数id

	idStr := c.Param("id")
	cid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil { // 请求参数有误
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 2. 根据参数id调用逻辑层查询详情
	communityList, err := logic.GetCommunityDetail(cid)
	if err != nil { // 出错了
		if errors.Is(err, mysql.ErrorInvalidID) { // 错误的ID
			ResponseError(c, CodeInvalidParams)
			return
		} else { // 其他数据库错误
			ResponseError(c, CodeServeBusy) // 统一返回服务繁忙 不要把详细信息返回给用户
			return
		}
	}

	ResponseSuccess(c, communityList)

}

// TopCommunityHandler 查询当前最热社区
func TopCommunityHandler(c *gin.Context) {

	topK := &models.TopKCommunity{
		TopK: 3,
	}
	err := c.ShouldBindQuery(topK)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	topCommunity, err := logic.TopCommunity(topK)
	if err != nil {
		ResponseError(c, CodeServeBusy) // 统一返回服务繁忙 不要把详细信息返回给用户
		return
	}
	ResponseSuccess(c, topCommunity)
}
