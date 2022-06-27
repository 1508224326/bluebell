package controller

import (
	"errors"
	"github.com/aiyouyo/bluebell/dao/mysql"
	"github.com/aiyouyo/bluebell/logic"
	"github.com/aiyouyo/bluebell/models"
	"github.com/aiyouyo/bluebell/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignHandler 用户注册的接口
// @Summary 用户注册的接口
// @Description 输入用户名和密码即可完成注册
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param object json models.ParamSignUp true "注册参数"
// @Security ApiKeyAuth
// @Success 200 {object}
// @Router /sign [post]
func SignHandler(c *gin.Context) {
	// 1. 获取参数参数校验
	var p = new(models.ParamSignUp) // 注册参数的结构体
	// 将参数绑定到这个结构体中
	if err := c.ShouldBindJSON(p); err != nil { // 只能检查字段是否正确 类型是否正确
		zap.L().Error("SignUp with invalid param", zap.Error(err))

		// 获取validator.ValidationErrors类型的errors
		errors, ok := err.(validator.ValidationErrors) // 类型断言
		if !ok {                                       // 如果不是 就是json序列化出错了
			ResponseError(c, CodeInvalidParams)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParams, ClearTopErr(errors.Translate(trans)))

		return
	}

	// 2. 业务处理  调用逻辑层的注册函数
	if err := logic.SignUp(p); err != nil && errors.Is(err, mysql.ErrorUserExist) { // 如果出错且错误是用户名已存在
		ResponseError(c, CodeUserExist)
		return
	} else if err != nil { // 数据库查询错误
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServeBusy, err.Error())
		return
	}

	// 3. 返回注册成功的响应
	ResponseSuccess(c, nil)
}

// LoginHandler 用户登录处理的接口
// @Summary 用户登陆的接口
// @Description 输入用户名和密码完成登陆
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param object json models.ParamLogin true "登陆参数"
// @Security ApiKeyAuth
// @Success 200 {object}
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	var p = new(models.ParamLogin)
	// 1. 获取请求参数
	if err := c.ShouldBindJSON(p); err != nil {

		errors, ok := err.(validator.ValidationErrors) // 是否是参数校验的错误
		if !ok {                                       // 表明是序列化出错
			ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParams, ClearTopErr(errors.Translate(trans)))
		return
	}

	// 2. 请求登录
	loginUser, err := logic.Login(p)

	// 登录失败
	if err != nil {
		if errors.Is(err, mysql.ErrorInvalidPassword) { // 用户名或密码错误
			ResponseError(c, CodeInvalidPassword)
		} else if errors.Is(err, mysql.ErrorUserNotExist) { // 用户不存在
			ResponseError(c, CodeUserNotExist)
		} else if errors.Is(err, jwt.ErrorInvalidToken) { // token生成错误
			zap.L().Error("gen jwt failed, err:", zap.Error(err))
			ResponseError(c, CodeServeBusy)
		} else { // 数据库其他错误
			zap.L().Error("login with invalid failed, err:", zap.Error(err))
			ResponseErrorWithMsg(c, CodeServeBusy, err.Error())
		}
		return
	}

	// 登录成功 返回数据
	ResponseSuccess(c, loginUser)
}
