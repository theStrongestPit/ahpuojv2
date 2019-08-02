package request

type Issue struct {
	Title     string `json:"title" binding:"required,max=20"`
	ProblemId int    `json:"problem_id"  binding:"gte=0"`
}

type Reply struct {
	Content     string `json:"content" binding:"required"`
	ReplyId     int    `json:"reply_id"  binding:"gte=0"`
	ReplyUserId int    `json:"reply_user_id"  binding:"gte=0"`
}
