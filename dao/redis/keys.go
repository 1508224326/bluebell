package redis

import "strings"

const (
	KeySep              = ":" // 分隔符
	KeyPrefix           = "bluebell"
	KeyPostTimeZSet     = "post:time"
	KeyPostScoreZSet    = "post:score"
	KeyAuthorPostZSet   = "author:post"
	KeyCommunityPostSet = "community:post"
	KeyPostVoteZSet     = "post:vote"
)

// NewRedisKey 构造一个Key
func NewRedisKey(keys ...string) (StrKey string) {
	subKeys := make([]string, 0, len(keys)+1)
	subKeys = append(subKeys, KeyPrefix)
	subKeys = append(subKeys, keys...)
	StrKey = strings.Join(subKeys, KeySep)
	return
}
