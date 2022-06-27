package logic

import (
	"github.com/aiyouyo/bluebell/dao/mysql"
	"github.com/aiyouyo/bluebell/models"
	"github.com/aiyouyo/bluebell/pkg/jwt"
	"github.com/aiyouyo/bluebell/pkg/snowflake"
)

func SignUp(param *models.ParamSignUp) (err error) {
	// 判断用户是否存在

	if err = mysql.CheckUserExist(param.Username); err != nil { // 数据库查询出错

		return err
	}

	// 生成uid
	uid := snowflake.GenID()
	// 构造用户实例
	u := &models.UserBase{
		UserID:   uid,
		UserName: param.Username,
		Password: param.Password,
	}

	// 保存到数据库中
	err = mysql.InsertUser(u)
	return err
}

func Login(param *models.ParamLogin) (loginUser *models.LoginUser, err error) {
	user := &models.UserBase{
		UserName: param.Username,
		Password: param.Password,
	}

	// 登录
	if err = mysql.Login(user); err != nil { // 登录不成功
		return
	}

	var token string
	// 登录成功 生成token
	token, err = jwt.GenToken(user.UserName, user.UserID)

	loginUser = new(models.LoginUser)
	loginUser.UserID = user.UserID
	loginUser.UserName = user.UserName
	loginUser.CreateTime = user.CreateTime
	loginUser.Token = token

	return
}
