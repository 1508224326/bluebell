package models

// ParamSignUp 定义注册请求参数的结构体
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`                                  // 用户名 不为空
	Password   string `json:"password" binding:"required,min=6,max=20"`                     // 密码 不为空
	RePassword string `json:"re_password" binding:"required,min=6,max=20,eqfield=Password"` // 重复密码 不为空
}

// ParamLogin 定义登录时的参数结构体
type ParamLogin struct {
	Username string `json:"username" binding:"required"` // 用户名 不为空
	Password string `json:"password" binding:"required"` // 密码 不为空
}

// ParamVotePost 帖子投票的参数
type ParamVotePost struct {
	AuthorID  int64 `json:"author_id,string"`                         // 用户ID 可以为空
	PostID    int64 `json:"post_id,string" binding:"required"`        // 帖子ID 可以为空
	VoteValue int8  `json:"vote_value,string" binding:"oneof=-1 0 1"` // 投票值 [-1, 0, 1]
}

// ParamPosts 查询帖子的请求参数
type ParamPosts struct {
	Page        int64  `json:"page" form:"page" binding:"min=1"` // 页码
	Size        int64  `json:"size" form:"size"`                 // 每页数据量
	OrderBy     string `json:"order_by" form:"order_by"`         // 排序规则
	CommunityID int64  `json:"community_id" form:"community_id"` // 社区ID 可按照社区划分帖子
}

type TopKCommunity struct {
	TopK int32 `json:"top_k" form:"top_k"`
}

//type

const (
	OderByTime  = "time"  // 按照时间排序
	OderByScore = "score" // 按照分数排序
)
