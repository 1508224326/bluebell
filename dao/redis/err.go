package redis

import "errors"

var (
	ErrorPostExpired = errors.New("帖子发布已过期")
	ErrorRepeatVote  = errors.New("不可以重复投票")
	ErrorEmptyIds    = errors.New("未查询到相关")
)
