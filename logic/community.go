package logic

import (
	"github.com/aiyouyo/bluebell/dao/mysql"
	"github.com/aiyouyo/bluebell/dao/redis"
	"github.com/aiyouyo/bluebell/models"
)

func GetCommunityList() ([]*models.CommunityBase, error) {
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	return mysql.QueryCommunityByID(id)
}

// TopCommunity 查询当前最热社区
func TopCommunity(topK *models.TopKCommunity) ([]*models.CommunityBase, error) {

	topCommunity, err := redis.TopCommunity(topK.TopK)
	return topCommunity, err

}
