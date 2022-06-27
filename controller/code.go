package controller

// 定义一些常用的响应码

type ResponseCode int

const (
	CodeSuccess       ResponseCode = 1000 + iota // 成功的响应码
	CodeInvalidParams                            // 请求参数错误的响应码
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServeBusy

	CodeInvalidToken
	CodeWithoutToken
	CodeNeedLogin

	CodeVotedExpired
	CodeVoteRepeat
)

var codeMessage = map[ResponseCode]string{
	CodeSuccess:         "success",
	CodeInvalidParams:   "请求参数错误",
	CodeUserExist:       "用户名已存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServeBusy:       "服务器繁忙",
	CodeInvalidToken:    "无效的Token",
	CodeWithoutToken:    "未携带Token",
	CodeNeedLogin:       "需要登录",
	CodeVotedExpired:    "投票时间已过",
	CodeVoteRepeat:      "不可以重复投票",
}

// Msg 返回对应状态码的信息
func (rc ResponseCode) Msg() string {
	msg, ok := codeMessage[rc]
	if !ok {
		msg = codeMessage[CodeServeBusy]
	}
	return msg
}
