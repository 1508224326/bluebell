package mysql

import (
	"database/sql"
	"github.com/aiyouyo/bluebell/models"
	"go.uber.org/zap"
)

// GetCommunityList 获取所有社区名字和id
func GetCommunityList() (communityList []*models.CommunityBase, err error) {

	sqlStr := "select community_id, community_name from community"

	err = db.Select(&communityList, sqlStr)

	if err != nil && err == sql.ErrNoRows { // 未查询到 空内容
		zap.L().Warn("there is no community in db")
		return communityList, nil
	}
	return communityList, err
}

// QueryCommunityByID 根据社区的id查询详细信息
func QueryCommunityByID(id int64) (community *models.CommunityDetail, err error) {
	sqlStr := "select community_id, community_name, introduction, create_time from community where id = ?"
	community = new(models.CommunityDetail)

	err = db.Get(community, sqlStr, id)

	if err != nil && err == sql.ErrNoRows {
		zap.L().Warn("query community by ID failed")
		err = ErrorInvalidID
	}
	return
}
