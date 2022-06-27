package redis

import (
	"github.com/aiyouyo/bluebell/dao/mysql"
	"github.com/aiyouyo/bluebell/models"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

// TopCommunity 查询最热社区
func TopCommunity(topK int32) ([]*models.CommunityBase, error) {
	communityBase, err := mysql.GetCommunityList() // 查询所有的社区信息
	if err != nil {
		return nil, err
	}
	topKey := NewRedisKey("TopCommunity")
	if rdb.Exists(topKey).Val() < 1 { // 如果不存在缓存的topCommunity 就开始计算

		pipeline := rdb.Pipeline()
		for _, c := range communityBase { // 计算每条社区下的帖子数量
			cid := strconv.FormatInt(c.ID, 10)
			pipeline.SCard(NewRedisKey(KeyCommunityPostSet, cid))
		}
		exec, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}

		pipeline = rdb.Pipeline()
		for idx, cmder := range exec { // 添加帖子数量和社区的id到redis中 然后返回前3
			size := cmder.(*redis.IntCmd).Val()
			pipeline.ZAdd(topKey, redis.Z{
				Score:  float64(size),
				Member: idx,
			})
		}
		pipeline.Expire(topKey, time.Second*180) // 三分钟过期
		_, err = pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	index, err := getRangeIDsByKey(1, int64(topK), topKey)
	if err != nil {
		return nil, err
	}

	// 构造一个top3社区数组
	topCommunity := make([]*models.CommunityBase, len(index), topK)
	for i, strIdx := range index {
		idx, _ := strconv.ParseInt(strIdx, 10, 32)
		topCommunity[i] = communityBase[idx]
	}

	return topCommunity, nil
}
