package main

type PostId string

type MediaId string

type UserId string

// 帖子，除了文本内容，可以附加相片，视频
type Post struct {
	Id      PostId `json:"id,omitempty" genji:"id"`
	UserId  UserId `json:"userId,omitempty" genji:"userId"`
	Content string `json:"content,omitempty" genji:"content"`
	// 附加的媒体文件ID
	Medias []*MediaId `json:"medias,omitempty" genji:"medias"`
	// 将针对Post的评论直接嵌入Media
	Comments []*Comment `json:"comments,omitempty" genji:"comments"`
	// 发布时间
	Time int64 `json:"time,omitempty" genji:"time"`
}

// 媒体，包括相片，视频
// 一个Media必须属于某个Post。如果某个Media不属于任何Post，
// 可能是Post发布未完成导致的，可以将该Media删除
type Media struct {
	Id MediaId `json:"id,omitempty" genji:"id"`
	// 这里的UserId是为了方便浏览用户的媒体文件
	UserId UserId    `json:"userId,omitempty" genji:"userId"`
	Type   MediaType `json:"type,omitempty" genji:"type"`
	URL    string    `json:"url,omitempty" genji:"url"`
	// 将针对Media的评论直接嵌入Media
	Comments []*Comment `json:"comments,omitempty" genji:"comments"`
	Posted   *bool      `json:"posted,omitempty" genji:"posted"`
	// 上传时间
	Time int64 `json:"time,omitempty" genji:"time"`
}

type Comment struct {
	UserId  UserId `json:"userId,omitempty" genji:"userId"`
	Content string `json:"content,omitempty" genji:"content"`
	Time    int64  `json:"time,omitempty" genji:"time"` // 发布时间
}

// 用户
type User struct {
	Id       UserId
	Username string
	Password string
	Avatar   string
	Time     int64 // 注册时间
}

type MediaType int

const (
	MediaType_Unknown MediaType = iota + 1
	MediaType_Picture
	MediaType_Video
)
