package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/aiyouyo/bluebell/models"
)

// 加密的字符串
const encrypt = "bluebell.com"

// CheckUserExist 检查新用户名用户是否存在了
func CheckUserExist(username string) (err error) {
	sqlStr := "select count(user_id) from user where username = ?"
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 { // 用户名已存在
		return ErrorUserExist
	}
	return nil
}

// InsertUser 插入用户到数据库中
func InsertUser(user *models.UserBase) (err error) {
	//	执行sql 语句入库
	sqlStr := "insert into user(user_id, username, password) values(?, ?, ?)"
	// 密码加密
	user.Password = encryptPwd(user.Password)
	_, err = db.Exec(sqlStr, user.UserID, user.UserName, user.Password)
	return
}

// Login 用户登录校验
func Login(user *models.UserBase) (err error) {
	oPwd := user.Password // 记录一下密码
	sqlStr := "select user_id, username, password, create_time from user where username = ?"

	// 序列化进user里
	err = db.Get(user, sqlStr, user.UserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // 用户不存在
			return ErrorUserNotExist
		} else {
			return err
		}
	}

	password := encryptPwd(oPwd) // 加密一下密码
	// 查询出来的密码是否相等
	if password != user.Password {
		return ErrorInvalidPassword
	}

	return
}

func GetUserByTD(uid int64) (user *models.ApiUser, err error) {
	sqlStr := "select user_id, username, create_time from user where user_id = ?"
	user = new(models.ApiUser)
	err = db.Get(user, sqlStr, uid)
	return
}

// 密码加密
func encryptPwd(oPwd string) string {
	hash := md5.New()
	hash.Write([]byte(encrypt))
	return hex.EncodeToString(hash.Sum([]byte(oPwd)))
}
