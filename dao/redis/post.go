package redis

import (
	"fmt"
	"github.com/aiyouyo/bluebell/models"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

func CreatPost(postID, communityID int64) (err error) {
	pipeline := rdb.TxPipeline()
	// 增加一个帖子和时间的集合
	pipeline.ZAdd(NewRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	pipeline.ZAdd(NewRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 社区增加帖子
	cid := strconv.FormatInt(communityID, 10)
	pipeline.SAdd(NewRedisKey(KeyCommunityPostSet, cid), postID)
	_, err = pipeline.Exec()
	return
}

// GetPostsIDsInOrder 查询帖子的ID 指定数量 和排序规则 按时间 按分数
func GetPostsIDsInOrder(page, size int64, orderBy string) (ids []string, err error) {
	// 默认以时间查询  获取对应的Key
	key := getOrderKey(orderBy)
	// 倒序查询指定数量的元素
	ids, err = getRangeIDsByKey(page, size, key)
	fmt.Println(ids)
	return
}

func GetVotedIdsByUserID(uid string, page, size int64) (ids []string, err error) {
	key := NewRedisKey(KeyAuthorPostZSet, uid) // 用户点赞的存储帖子
	fmt.Println("-----------", key, "-------------")
	// 查询点赞过的帖子id
	ids, err = getRangeIDsByKey(page, size, key)
	return
}

// GetCommunityPostIDsInOrder 查询指定社区帖子的ID 指定数量 和排序规则 按时间 按分数
func GetCommunityPostIDsInOrder(communityID, page, size int64, orderBy string) (ids []string, err error) {
	// 社区的ID-->str:
	communityIDStr := strconv.FormatInt(communityID, 10)
	// 社区的key:
	cKey := NewRedisKey(KeyCommunityPostSet, communityIDStr)
	orderKey := getOrderKey(orderBy) // 获取按时间还是按分数排序的Key
	// 利用 ZInterStore 把社区帖子和分数帖子和时间帖子生生新的Zset
	// 同时利用缓存减少inter操作的耗时
	key := NewRedisKey("community", communityIDStr, orderBy) // 生成新的Key
	if rdb.Exists(key).Val() < 1 {                           // 缓存的key不存在, 利用 ZInterStore 把社区帖子和分数帖子和时间帖子生生新的Zset
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey)
		// 设置缓存过期时间
		pipeline.Expire(key, 180*time.Second) // 可以写到配置文件中
		_, err = pipeline.Exec()
		if err != nil {
			return
		}
	}

	ids, err = getRangeIDsByKey(page, size, key)
	return
}

// 根据key查询指定数量的数据  倒序 从大到小
func getRangeIDsByKey(page, size int64, key string) ([]string, error) {
	// 计算元素的范围
	start := (page - 1) * size
	end := start + size - 1
	// 获取元素
	ids, err := rdb.ZRevRange(key, start, end).Result()
	return ids, err
}

func getOrderKey(orderBy string) (key string) {
	key = NewRedisKey(KeyPostTimeZSet)
	if orderBy == models.OderByScore { // 按分数查询
		key = NewRedisKey(KeyPostScoreZSet)
	}
	return
}
