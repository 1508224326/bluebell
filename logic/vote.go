package logic

import (
	"github.com/aiyouyo/bluebell/dao/redis"
	"github.com/aiyouyo/bluebell/models"
	"strconv"
)

// VoteForPost 给帖子投票
func VoteForPost(vote *models.ParamVotePost) (err error) {

	strAuthorID := strconv.FormatInt(vote.AuthorID, 10)
	strPostID := strconv.FormatInt(vote.PostID, 10)
	float64VoteValue := float64(vote.VoteValue)
	err = redis.VoteForPost(strAuthorID, strPostID, float64VoteValue)
	return
}
