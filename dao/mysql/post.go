package mysql

import (
	"github.com/aiyouyo/bluebell/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

const (
	MAXSIZE = 10 // 最大查询数
)

// 关于帖子的一些操作

func CreatePost(post *models.PostBase) error {
	sqlStr := `insert into post(post_id, title, content, author_id, community_id) values(?, ?, ?, ?, ?)`
	_, err := db.Exec(sqlStr, post.ID, post.Title, post.Content, post.AuthorID, post.CommunityID)
	return err
}

func GetPostDetailByID(pid int64) (data *models.PostBase, err error) {
	sqlStr := "select post_id, title, content, author_id, community_id, status, create_time, update_time from post where post_id = ?"
	data = new(models.PostBase)
	err = db.Get(data, sqlStr, pid)
	return
}

// GetPostList 获取所有帖子
func GetPostList(page, size int64) (post []*models.PostBase, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, status, create_time, update_time from post 
			   ORDER BY create_time 
			   DESC
			   limit ?, ?` // 保证查询的最新的帖子
	if size > MAXSIZE {
		size = MAXSIZE
	}
	post = make([]*models.PostBase, 0, size)
	err = db.Select(&post, sqlStr, (page-1)*size, size)
	return
}

func GetPostListByIDs(ids []string) (post []*models.PostBase, err error) {

	//strIDs := make([]string, 0, len(ids))
	// 利用sqlx的 sqlx.In方法查询一批
	sqlStr := `select post_id, title, content, author_id, community_id, status, create_time, update_time from post 
               where post_id in (?)
               order by FIND_IN_SET(post_id, ?)`

	query, args, err1 := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err1 != nil {
		return nil, err1
	}

	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = db.Rebind(query)
	err = db.Select(&post, query, args...)
	return

}
