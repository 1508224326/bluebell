package models

import "time"

type UserBase struct {
	UserID   int64  `db:"user_id"`
	UserName string `db:"username"`
	Password string `db:"password"`
	//Email      string    `json:"email" db:"email"`
	CreateTime time.Time `db:"create_time"`
}

// LoginUser 用户登录成功后返回的数据
type LoginUser struct {
	UserID     int64     `json:"user_id,string" db:"user_id"`
	UserName   string    `json:"user_name" db:"username"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
	Token      string    `json:"token"`
}

// ApiUser 获取详细的作者信息
type ApiUser struct {
	UserID   int64  `json:"user_id,string" db:"user_id"`
	UserName string `json:"user_name" db:"username"`
	//Email      string    `json:"email" db:"email"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
}
