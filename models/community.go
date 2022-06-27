package models

import "time"

// CommunityBase　社区的基本信息 ID　和　名字
type CommunityBase struct {
	ID            int64  `json:"id" db:"community_id"`               // 社区ID
	CommunityName string `json:"community_name" db:"community_name"` // 社区名称
}

// CommunityDetail 社区详情
type CommunityDetail struct {
	*CommunityBase           // 社区ID 和 名字
	Introduction   string    `json:"introduction" db:"introduction"` // 描述信息
	CreateTime     time.Time `json:"create_time" db:"create_time"`   // 创建时间
}
