package request

type New struct {
	Title   string `json:"title" binding:"required,max=20"`
	Content string `json:"content"`
}
