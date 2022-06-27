package logic

import (
	"fmt"
	"github.com/aiyouyo/bluebell/dao/mysql"
	"github.com/aiyouyo/bluebell/dao/redis"
	"github.com/aiyouyo/bluebell/models"
	"github.com/aiyouyo/bluebell/pkg/snowflake"
	"go.uber.org/zap"
	"strconv"
)

func CreatePost(post *models.PostBase) (err error) {

	// 1. 生成帖子ID
	post.ID = snowflake.GenID()

	// 2. 保存到数据库
	err = mysql.CreatePost(post)
	if err != nil {
		return
	}
	err = redis.CreatPost(post.ID, post.CommunityID)
	return
}

func GetPostDetailByID(pid string) (data *models.ApiPostDetail, err error) {

	id, err := strconv.ParseInt(pid, 10, 64)
	if err != nil {
		return
	}

	// 1. 查询帖子信息
	post, err := mysql.GetPostDetailByID(id) // 这是帖子的详细信息
	if err != nil {
		zap.L().Error("mysql.GetPostDetailByID(pid) failed",
			zap.Int64("pid", id), zap.Error(err))
		return
	}

	// 2. 根据作者id查询作者信息
	user, err := mysql.GetUserByTD(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByTD(post.AuthorID) failed",
			zap.Int64("uid", post.AuthorID), zap.Error(err))
		return
	}

	// 3. 查询帖子的社区详细信息
	community, err := mysql.QueryCommunityByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.QueryCommunityByID(post.CommunityID) failed",
			zap.Int64("cid", post.CommunityID), zap.Error(err))
		return
	}

	// 5. 查询投票信息
	fmt.Println(pid)
	upVotes, downVotes, err := redis.GetVotesData([]string{pid})
	if err != nil {
		zap.L().Error("redis.GetVotesData([]string{strconv.FormatInt(pid, 64)})",
			zap.Int64("pid", post.ID), zap.Error(err))
		return
	}

	fmt.Println("[", upVotes, downVotes, "]")

	// 4. 组合数据
	data = &models.ApiPostDetail{
		PostBase:        post,
		ApiUser:         user,
		CommunityDetail: community,
		UpNumber:        upVotes[0],
		DownNumber:      downVotes[0],
	}
	return
}

// fillPostDetail 填充帖子的详细信息
func fillPostDetail(basePosts []*models.PostBase, ups, downs []int64) (data []*models.ApiPostDetail, err error) {

	data = make([]*models.ApiPostDetail, 0, len(basePosts)) // 创建空间
	var (
		user      *models.ApiUser
		community *models.CommunityDetail
	)

	flag := ups == nil && downs == nil
	var postDetail *models.ApiPostDetail
	for idx, post := range basePosts {
		// 2. 根据作者id查询作者信息
		user, err = mysql.GetUserByTD(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByTD(post.AuthorID) failed",
				zap.Int64("uid", post.AuthorID), zap.Error(err))
			continue
		}

		// 3. 查询帖子的社区详细信息
		community, err = mysql.QueryCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.QueryCommunityByID(post.CommunityID) failed",
				zap.Int64("cid", post.CommunityID), zap.Error(err))
			continue
		}
		if flag {
			postDetail = &models.ApiPostDetail{ // 这样才可以
				PostBase:        post,
				ApiUser:         user,
				CommunityDetail: community,
				UpNumber:        0,
				DownNumber:      0,
			}
		} else {
			postDetail = &models.ApiPostDetail{ // 这样才可以
				PostBase:        post,
				ApiUser:         user,
				CommunityDetail: community,
				UpNumber:        ups[idx],
				DownNumber:      downs[idx],
			}
		}

		data = append(data, postDetail) // 追加进去
	}
	return
}

// GetPostList 获取所有帖子的详细信息
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {

	var posts []*models.PostBase
	// 获取所有帖子
	posts, err = mysql.GetPostList(page, size)

	if err != nil {
		return nil, err
	}

	// 填充帖子的详细数据
	data, err = fillPostDetail(posts, nil, nil)
	if err != nil {
		zap.L().Error("fillPostDetail failed, err:", zap.Error(err))
		fmt.Println(err)
		return
	}

	return
}

// GetPostListV2 获取帖子 按照时间或分数排序返回
func GetPostListV2(queryStr *models.ParamPosts) (data []*models.ApiPostDetail, err error) {
	// 1. 从Redis中查询最新的size条帖子id
	var ids []string
	ids, err = getIDsFromRedis(queryStr)
	// 如果是空的
	if len(ids) <= 0 {
		zap.L().Error("ids is empty, err: ", zap.Error(err))
		return data, nil
	}
	// 查询数据
	data, err = expandPostDataByIds(ids)

	// 4. 返回数据
	return
}

// GetVotedPostList 获得点赞的帖子
func GetVotedPostList(uid int64, queryStr *models.ParamPosts) (data []*models.ApiPostDetail, err error) {

	ids, err := getVotedIDsFromRedis(uid, queryStr)
	if err != nil {
		return data, err
	}
	// 如果是空的
	if len(ids) <= 0 {
		zap.L().Error("ids is empty, err: ", zap.Error(err))
		return data, nil
	}

	// 填充ids的数据
	data, err = expandPostDataByIds(ids)
	return
}

// 获取帖子的ID
func getIDsFromRedis(queryStr *models.ParamPosts) (ids []string, err error) {
	if queryStr.CommunityID < 0 { // 表明 普通查询 不按照社区划分
		ids, err = redis.GetPostsIDsInOrder(queryStr.Page, queryStr.Size, queryStr.OrderBy)
	} else {
		ids, err = redis.GetCommunityPostIDsInOrder(queryStr.CommunityID, queryStr.Page, queryStr.Size, queryStr.OrderBy)
	}

	return
}

// 查询用户点赞的帖子id
func getVotedIDsFromRedis(uid int64, queryStr *models.ParamPosts) (ids []string, err error) {
	// 将用户id转为字符串
	uidStr := strconv.FormatInt(uid, 10)

	ids, err = redis.GetVotedIdsByUserID(uidStr, queryStr.Page, queryStr.Size)

	return

}

// 根据帖子id填充帖子的数据
func expandPostDataByIds(ids []string) (data []*models.ApiPostDetail, err error) {

	// 2. 查询帖子的投票数据
	ups, downs, err := redis.GetVotesData(ids)
	if err != nil {
		zap.L().Error("redis.GetVotesData failed, err: ",
			zap.Strings("ids", ids), zap.Error(err))
		return
	}

	var postBase []*models.PostBase
	// 2. 拿着这些id去MySQL查数据
	postBase, err = mysql.GetPostListByIDs(ids)
	if err != nil {
		zap.L().Error("mysql.GetPostListByIDs(ids) failed, err: ",
			zap.Strings("ids", ids), zap.Error(err))
		return
	}

	// 3. 填充数据
	data, err = fillPostDetail(postBase, ups, downs)
	if err != nil {
		zap.L().Error("fillPostDetail failed, err:", zap.Error(err))
		fmt.Println(err)
		return
	}

	// 4. 返回数据
	return
}
