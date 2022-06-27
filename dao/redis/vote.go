package redis

import (
	"errors"
	"fmt"
	"github.com/aiyouyo/bluebell/settingts"
	"github.com/go-redis/redis"
	"math"
	"time"
)

const (
	VotePreValue = 432 // 指定一个赞成票是432分  432*200 = 86400秒 = 一天 也就是说获得200张票可以将帖子续一天
)

/*
分为三种情况

value:   表示当前票数
oValue : 过去投票
|value - oValue|

value = 1: 投赞成票
	1. 之前投反对票  现在投赞成票          --->   更新票数    |value - oValue| = 2  -->   + 432*2   (value - oValue)>0
	2. 之前没投票   现在投赞成票           --->   更新票数    |value - oValue| = 1  -->   + 432*1   (value - oValue)>0

value = 0: 取消投票
	1. 之前投返对票  现在取消投票          --->   更新票数    |value - oValue| = 1  -->   + 432*1   (value - oValue)>0
	2. 之前投赞成票  现在取消投票          --->   更新票数    |value - oValue| = 1  -->   - 432*1   (value - oValue)<0

value = -1: 投反对票
	1. 之前没投票   现在投反对票          --->  更新票数      |value - oValue| = 1  -->   - 432*1   (value - oValue)<0
	2. 之前投赞成票 现在投反对票           --->   更新票数     |value - oValue| = 2  -->   - 432*2   (value - oValue)<0
*/

// VoteForPost 给帖子投票
func VoteForPost(authorID, postID string, value float64) error {

	fmt.Println(rdb.Ping().Result())

	// 判断帖子是否过期
	publishTime, err := rdb.ZScore(NewRedisKey(KeyPostTimeZSet), postID).Result() // 获得当前帖子的发布时间
	if err != nil {

		if errors.Is(err, redis.Nil) {
			return ErrorPostExpired
		}
		return err
	}
	// 过期了
	duration := float64(settingts.Conf.ExpireConfig.PostExpire * 24 * 3600) // 过期天数
	if float64(time.Now().Unix())-publishTime > duration {
		return ErrorPostExpired
	}

	// 获取用户是否给当前帖子投过票
	oValue, err := rdb.ZScore(NewRedisKey(KeyPostVoteZSet, postID), authorID).Result()
	fmt.Println("----------------------------")
	fmt.Println(oValue, value)
	fmt.Println("----------------------------")
	// value - oValue
	var op float64
	if value > oValue { //
		op = 1.0
	} else if value < oValue {
		op = -1.0
	} else {
		//oValue == value  // 投相同的票是不可以的
		return ErrorRepeatVote //
	}
	// 计算两次投票的差值
	diff := math.Abs(oValue - value)

	pipeline := rdb.TxPipeline()
	// 为该帖子投票
	pipeline.ZIncrBy(NewRedisKey(KeyPostScoreZSet), op*diff*VotePreValue, postID)
	// 计算该用户为该帖子的投票数据
	if value == 0 { // 取消点赞
		pipeline.ZRem(NewRedisKey(KeyAuthorPostZSet, authorID), postID)
	} else { // 记录点赞的帖子信息
		pipeline.ZAdd(NewRedisKey(KeyAuthorPostZSet, authorID), redis.Z{
			Score:  float64(time.Now().Unix() + settingts.Conf.ExpireConfig.VotedExpire), // 记录点赞的时间加上过期时间(1周)
			Member: postID,
		})
		//postID: vote: user

		pipeline.ZAdd(NewRedisKey(KeyPostVoteZSet, postID), redis.Z{
			Score:  value,
			Member: authorID,
		})
	}
	_, err = pipeline.Exec()

	return err
}

// GetVotesData 获得帖子的投票数量
func GetVotesData(postIDs []string) (upVotes, downVotes []int64, err error) {

	// 使用一个pipeline将命令发送过去
	pipeline := rdb.Pipeline()

	for _, pid := range postIDs {
		key := NewRedisKey(KeyPostVoteZSet, pid)
		pipeline.ZCount(key, "1", "1")
		pipeline.ZCount(key, "-1", "-1")
	}
	cmder, err := pipeline.Exec()

	if err != nil {
		return
	}
	size := len(postIDs)
	upVotes = make([]int64, 0, size)
	downVotes = make([]int64, 0, size)

	for i, l := 0, len(cmder); i < l; i += 2 {
		up := cmder[i].(*redis.IntCmd).Val() // 取值
		down := cmder[i+1].(*redis.IntCmd).Val()
		upVotes = append(upVotes, up)
		downVotes = append(downVotes, down)
	}

	return
}
