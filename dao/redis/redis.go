package redis

import (
	"fmt"
	"github.com/aiyouyo/bluebell/models"
	"github.com/aiyouyo/bluebell/settingts"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var rdb *redis.Client

func Init() (err error) {

	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString("redis.host"),
			viper.GetInt("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool_size"),
	})

	_, err = rdb.Ping().Result()

	return
}

func InitV2(cfg *settingts.RedisConfig) (err error) {

	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	_, err = rdb.Ping().Result()

	return
}

func CLose() {
	_ = rdb.Close()
}

// SetAllPosts 同步数据
func SetAllPosts(posts []*models.PostBase) (err error) {
	pipeline := rdb.TxPipeline()
	var cKey string
	for _, post := range posts {
		pid := strconv.FormatInt(post.ID, 10)
		pipeline.ZCount(NewRedisKey(KeyPostVoteZSet, pid), "-1", "-1")
		pipeline.ZCount(NewRedisKey(KeyPostVoteZSet, pid), "1", "1")
		fmt.Println(post.ID, post.CreatTime)
	}
	cmder, err := pipeline.Exec()
	if err != nil {
		return err
	}
	m := make(map[int64]int64, len(posts))
	for i, l := 0, len(cmder); i < l; i += 2 {
		v1 := cmder[i].(*redis.IntCmd).Val()
		v2 := cmder[i+1].(*redis.IntCmd).Val()
		fmt.Println(v1, v2)
		score := -v1*432 + v2*432
		m[posts[i/2].ID] = score
	}

	fmt.Println(m)
	pipeline = rdb.Pipeline()
	for _, post := range posts {
		pipeline.ZAdd(NewRedisKey(KeyPostTimeZSet), redis.Z{
			Score:  float64(post.CreatTime.Unix()),
			Member: post.ID,
		})
		cKey = strconv.FormatInt(post.CommunityID, 10)
		pipeline.SAdd(NewRedisKey(KeyCommunityPostSet, cKey), post.ID)
		pipeline.ZAdd(NewRedisKey(KeyPostScoreZSet), redis.Z{
			Score:  float64(post.CreatTime.Unix()) + float64(m[post.ID]),
			Member: post.ID,
		})
	}
	_, err = pipeline.Exec()

	return
}
