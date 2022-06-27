package models

import "time"

type PostBase struct {
	Status      int32     `json:"status" db:"status"`
	ID          int64     `json:"id,string" db:"post_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	AuthorID    int64     `json:"author_id,string" db:"author_id"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreatTime   time.Time `json:"creat_time" db:"create_time"`
	UpdateTime  time.Time `json:"update_time,omitempty" db:"update_time"`
}

// ApiPostDetail 帖子的详情页
type ApiPostDetail struct {
	*PostBase                           // 帖子的具体信息
	*ApiUser         `json:"author"`    // 作者的具体信息
	*CommunityDetail `json:"community"` // 社区的具体信息
	UpNumber         int64              `json:"up_number"`
	DownNumber       int64              `json:"down_number"`
}

func NewApiPostDetail(base *PostBase, user *ApiUser, detail *CommunityDetail, up, down int64) *ApiPostDetail {
	return &ApiPostDetail{
		base,
		user,
		detail,
		up,
		down,
	}
}
